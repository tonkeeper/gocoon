package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	magicV2 uint32 = 0x3a2f9be2
	magicV3 uint32 = 0xe4a8604b
	magicV4 uint32 = 0x90ac88d7

	tlsTypeMagic             uint32 = 0x12eb4386
	tlsCombinatorMagic       uint32 = 0x5c0a1ed5
	tlsCombinatorLeft        uint32 = 0x4c12c6d9
	tlsCombinatorLeftBuiltin uint32 = 0xcd211f63
	tlsCombinatorRightV2     uint32 = 0x2c064372
	tlsArgV2Magic            uint32 = 0x29dfe61b

	tlsExprNat  uint32 = 0xdcb49bd8
	tlsExprType uint32 = 0xecc9da78

	tlsTypeVar  uint32 = 0x0142ceae
	tlsTypeExpr uint32 = 0xc1863d08
	tlsArray    uint32 = 0xd9fb20de

	tlsNatConst    uint32 = 0x8ce940b1
	tlsNatConstOld uint32 = 0xdcb49bd8
	tlsNatVar      uint32 = 0x4e8a14f0

	flagNoVar = 1 << 21
)

// Well-known builtin type IDs.
const (
	builtinNat       uint32 = 0x2cecf817 // "#"
	builtinInt       uint32 = 0xa8509bda
	builtinLong      uint32 = 0x22076cba
	builtinDouble    uint32 = 0x2210c154
	builtinString    uint32 = 0xb5286e24
	builtinVector    uint32 = 0x1cb5c415
	builtinBoolFalse uint32 = 0xbc799737
	builtinBoolTrue  uint32 = 0x997275b5
)

type TLOSchema struct {
	Version      int
	Types        []*TLOType
	Constructors []*TLOCombinator
	Functions    []*TLOCombinator
	TypesByID    map[uint32]*TLOType
}

type TLOType struct {
	ID              uint32
	Name            string
	ConstructorsNum int32
	Flags           int32
	Arity           int32
	Constructors    []*TLOCombinator
}

type TLOCombinator struct {
	ID         uint32
	Name       string
	TypeID     uint32
	ResultType *TLOType
	Args       []*TLOArg
	IsBuiltin  bool
	ResultExpr *TypeExprNode
}

type TLOArg struct {
	Flags       int32
	VarNum      int32
	Name        string
	IsOptional  bool
	ExistVarNum int32
	ExistVarBit int32
	TypeExpr    *TypeExprNode
}

type TypeExprNode struct {
	Kind         string // "type_expr" | "type_var" | "array" | "nat_const" | "nat_var"
	TypeID       uint32
	Flags        uint32
	Children     []*TypeExprNode // type parameters (e.g. element type of Vector)
	VarNum       int32
	Value        int32         // nat_const value
	Diff         int32         // nat_var diff
	ArrayArgs    []*TLOArg     // array repeated fields
	Multiplicity *TypeExprNode // array multiplicity
}

// --- reader ------------------------------------------------------------------

type tloReader struct {
	r             *bytes.Reader
	schemaVersion int
	schemaOptFlag int32
	schemaHasVars int32
	typesByID     map[uint32]*TLOType
}

func (r *tloReader) readUint32() (uint32, error) {
	var buf [4]byte
	if _, err := io.ReadFull(r.r, buf[:]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(buf[:]), nil
}

func (r *tloReader) readInt32() (int32, error) {
	v, err := r.readUint32()
	return int32(v), err
}

func (r *tloReader) readString() (string, error) {
	b := make([]byte, 1)
	if _, err := io.ReadFull(r.r, b); err != nil {
		return "", err
	}

	var length, headerSize int
	if b[0] < 254 {
		length = int(b[0])
		headerSize = 1
	} else if b[0] == 254 {
		extra := make([]byte, 3)
		if _, err := io.ReadFull(r.r, extra); err != nil {
			return "", err
		}
		length = int(extra[0]) | int(extra[1])<<8 | int(extra[2])<<16
		headerSize = 4
	} else {
		return "", fmt.Errorf("invalid string header byte: %d", b[0])
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(r.r, data); err != nil {
		return "", err
	}

	total := headerSize + length
	if pad := (4 - total%4) % 4; pad > 0 {
		if _, err := r.r.Seek(int64(pad), io.SeekCurrent); err != nil {
			return "", err
		}
	}

	return string(data), nil
}

func (r *tloReader) expect(expected uint32) error {
	got, err := r.readUint32()
	if err != nil {
		return err
	}
	if got != expected {
		return fmt.Errorf("expected 0x%08x, got 0x%08x", expected, got)
	}
	return nil
}

// --- top-level parse ---------------------------------------------------------

func ParseTLO(data []byte) (*TLOSchema, error) {
	if len(data)%4 != 0 {
		return nil, fmt.Errorf("invalid .tlo size %d: must be a multiple of 4 bytes", len(data))
	}
	r := &tloReader{r: bytes.NewReader(data)}

	magic, err := r.readUint32()
	if err != nil {
		return nil, fmt.Errorf("read magic: %w", err)
	}
	switch magic {
	case magicV2:
		r.schemaVersion = 2
		r.schemaOptFlag = 2
		r.schemaHasVars = 4
	case magicV3:
		r.schemaVersion = 3
		r.schemaOptFlag = 4
		r.schemaHasVars = 2
	case magicV4:
		r.schemaVersion = 4
		r.schemaOptFlag = 4
		r.schemaHasVars = 2
	default:
		return nil, fmt.Errorf("unknown TLO magic: 0x%08x", magic)
	}

	// date + version (ignored)
	if _, err := r.readInt32(); err != nil {
		return nil, fmt.Errorf("read date: %w", err)
	}
	if _, err := r.readInt32(); err != nil {
		return nil, fmt.Errorf("read version field: %w", err)
	}

	s := &TLOSchema{Version: r.schemaVersion, TypesByID: make(map[uint32]*TLOType)}
	r.typesByID = s.TypesByID

	typesCount, err := r.readInt32()
	if err != nil {
		return nil, fmt.Errorf("read types_count: %w", err)
	}
	for i := int32(0); i < typesCount; i++ {
		t, err := r.readType()
		if err != nil {
			return nil, fmt.Errorf("type[%d]: %w", i, err)
		}
		s.Types = append(s.Types, t)
		s.TypesByID[t.ID] = t
	}

	ctorCount, err := r.readInt32()
	if err != nil {
		return nil, fmt.Errorf("read constructors_count: %w", err)
	}
	var lastCtorName string
	for i := int32(0); i < ctorCount; i++ {
		posBefore := r.pos()
		c, err := r.readCombinator()
		if err != nil {
			return nil, fmt.Errorf("constructor[%d] at pos %d (after %q): %w", i, posBefore, lastCtorName, err)
		}
		lastCtorName = c.Name
		s.Constructors = append(s.Constructors, c)
		if t, ok := s.TypesByID[c.TypeID]; ok {
			t.Constructors = append(t.Constructors, c)
			c.ResultType = t
		}
	}

	fnCount, err := r.readInt32()
	if err != nil {
		return nil, fmt.Errorf("read functions_count: %w", err)
	}
	for i := int32(0); i < fnCount; i++ {
		f, err := r.readCombinator()
		if err != nil {
			return nil, fmt.Errorf("function[%d]: %w", i, err)
		}
		s.Functions = append(s.Functions, f)
		if t, ok := s.TypesByID[f.TypeID]; ok {
			f.ResultType = t
		}
	}

	if rest := r.r.Len(); rest != 0 {
		return nil, fmt.Errorf("unexpected trailing bytes: %d", rest)
	}

	return s, nil
}

func (r *tloReader) readType() (*TLOType, error) {
	if err := r.expect(tlsTypeMagic); err != nil {
		return nil, fmt.Errorf("TLS_TYPE: %w", err)
	}
	id, err := r.readUint32()
	if err != nil {
		return nil, err
	}
	name, err := r.readString()
	if err != nil {
		return nil, err
	}
	ctorNum, err := r.readInt32()
	if err != nil {
		return nil, err
	}
	flags, err := r.readInt32()
	if err != nil {
		return nil, err
	}
	flags &^= 1 | 8 | 16 | 1024 // clear bits 0,3,4,10
	arity, err := r.readInt32()
	if err != nil {
		return nil, err
	}
	// unused int64
	var unused [8]byte
	if _, err := io.ReadFull(r.r, unused[:]); err != nil {
		return nil, err
	}
	return &TLOType{ID: id, Name: name, ConstructorsNum: ctorNum, Flags: flags, Arity: arity}, nil
}

func (r *tloReader) readCombinator() (*TLOCombinator, error) {
	if err := r.expect(tlsCombinatorMagic); err != nil {
		return nil, fmt.Errorf("TLS_COMBINATOR: %w", err)
	}
	id, err := r.readUint32()
	if err != nil {
		return nil, err
	}
	name, err := r.readString()
	if err != nil {
		return nil, err
	}
	typeID, err := r.readUint32()
	if err != nil {
		return nil, err
	}

	leftMagic, err := r.readUint32()
	if err != nil {
		return nil, err
	}

	c := &TLOCombinator{ID: id, Name: name, TypeID: typeID}
	switch leftMagic {
	case tlsCombinatorLeft:
		c.Args, err = r.readArgsList()
		if err != nil {
			return nil, fmt.Errorf("args: %w", err)
		}
	case tlsCombinatorLeftBuiltin:
		c.IsBuiltin = true
	default:
		return nil, fmt.Errorf("unknown combinator-left magic: 0x%08x", leftMagic)
	}

	if err := r.expect(tlsCombinatorRightV2); err != nil {
		return nil, fmt.Errorf("TLS_COMBINATOR_RIGHT_V2: %w", err)
	}
	c.ResultExpr, err = r.readTypeExpr()
	if err != nil {
		return nil, fmt.Errorf("result expr: %w", err)
	}
	return c, nil
}

// readExpr is used only for type parameters (arity × read_expr inside TLS_TYPE_EXPR).
// It reads the outer TLS_EXPR_NAT / TLS_EXPR_TYPE wrapper first.

func (r *tloReader) readArgsList() ([]*TLOArg, error) {
	count, err := r.readInt32()
	if err != nil {
		return nil, err
	}
	if count < 0 {
		return nil, fmt.Errorf("negative args_count: %d", count)
	}
	args := make([]*TLOArg, count)
	for i := range args {
		args[i], err = r.readArg()
		if err != nil {
			return nil, fmt.Errorf("arg[%d]: %w", i, err)
		}
	}
	return args, nil
}

func (r *tloReader) readArg() (*TLOArg, error) {
	if err := r.expect(tlsArgV2Magic); err != nil {
		return nil, fmt.Errorf("TLS_ARG_V2: %w", err)
	}
	name, err := r.readString()
	if err != nil {
		return nil, err
	}
	flags, err := r.readInt32()
	if err != nil {
		return nil, err
	}

	rawFlags := flags
	isOptional := flags&r.schemaOptFlag != 0
	hasVars := flags&r.schemaHasVars != 0

	arg := &TLOArg{Name: name, Flags: rawFlags, VarNum: -1, IsOptional: isOptional, ExistVarNum: -1}

	if hasVars {
		arg.VarNum, err = r.readInt32()
		if err != nil {
			return nil, err
		}
	}
	if isOptional {
		arg.ExistVarNum, err = r.readInt32()
		if err != nil {
			return nil, err
		}
		arg.ExistVarBit, err = r.readInt32()
		if err != nil {
			return nil, err
		}
	}

	arg.TypeExpr, err = r.readTypeExpr()
	if err != nil {
		return nil, fmt.Errorf("type expr: %w", err)
	}
	return arg, nil
}

func (r *tloReader) readExpr() (*TypeExprNode, error) {
	tag, err := r.readUint32()
	if err != nil {
		return nil, err
	}
	switch tag {
	case tlsExprNat:
		return r.readNatExpr()
	case tlsExprType:
		return r.readTypeExpr()
	default:
		return nil, fmt.Errorf("unknown expr tag: 0x%08x", tag)
	}
}

func (r *tloReader) readNatExpr() (*TypeExprNode, error) {
	tag, err := r.readUint32()
	if err != nil {
		return nil, err
	}
	return r.readNatExprWithTag(tag)
}

func (r *tloReader) readNatExprWithTag(tag uint32) (*TypeExprNode, error) {
	switch tag {
	case tlsNatConst, tlsNatConstOld:
		val, err := r.readInt32()
		if err != nil {
			return nil, err
		}
		return &TypeExprNode{Kind: "nat_const", Value: val}, nil
	case tlsNatVar:
		diff, err := r.readInt32()
		if err != nil {
			return nil, err
		}
		varNum, err := r.readInt32()
		if err != nil {
			return nil, err
		}
		return &TypeExprNode{Kind: "nat_var", Diff: diff, VarNum: varNum}, nil
	default:
		return nil, fmt.Errorf("unknown nat expr tag: 0x%08x", tag)
	}
}

func (r *tloReader) pos() int64 {
	p, _ := r.r.Seek(0, io.SeekCurrent)
	return p
}

func (r *tloReader) peekBytes(n int) []byte {
	pos := r.pos()
	buf := make([]byte, n)
	nn, _ := r.r.Read(buf)
	r.r.Seek(pos, io.SeekStart)
	return buf[:nn]
}

func (r *tloReader) readTypeExpr() (*TypeExprNode, error) {
	tag, err := r.readUint32()
	if err != nil {
		return nil, err
	}
	switch tag {
	case tlsExprType: // some encoders wrap type exprs in the outer TLS_EXPR_TYPE tag
		return r.readTypeExpr()
	case tlsTypeVar:
		varNum, err := r.readInt32()
		if err != nil {
			return nil, err
		}
		flags, err := r.readUint32()
		if err != nil {
			return nil, err
		}
		return &TypeExprNode{Kind: "type_var", VarNum: varNum, Flags: flags}, nil

	case tlsTypeExpr:
		typeID, err := r.readUint32()
		if err != nil {
			return nil, err
		}
		flags, err := r.readUint32()
		if err != nil {
			return nil, err
		}
		arity, err := r.readInt32()
		if err != nil {
			return nil, err
		}
		if arity < 0 {
			return nil, fmt.Errorf("negative type arity: %d", arity)
		}
		if t, ok := r.typesByID[typeID]; ok && t.Arity != arity {
			return nil, fmt.Errorf("type arity mismatch for %q (0x%08x): got %d, expected %d", t.Name, typeID, arity, t.Arity)
		}
		node := &TypeExprNode{Kind: "type_expr", TypeID: typeID, Flags: flags | flagNoVar}
		for i := int32(0); i < arity; i++ {
			child, err := r.readExpr()
			if err != nil {
				return nil, fmt.Errorf("type param[%d]: %w", i, err)
			}
			node.Children = append(node.Children, child)
		}
		return node, nil

	case tlsArray:
		natTag, err := r.readUint32()
		if err != nil {
			return nil, err
		}
		mult, err := r.readNatExprWithTag(natTag)
		if err != nil {
			return nil, fmt.Errorf("array multiplicity: %w", err)
		}
		arrayArgs, err := r.readArgsList()
		if err != nil {
			return nil, fmt.Errorf("array args: %w", err)
		}
		return &TypeExprNode{Kind: "array", Multiplicity: mult, ArrayArgs: arrayArgs}, nil

	default:
		next := r.peekBytes(32)
		return nil, fmt.Errorf("unknown type expr tag: 0x%08x at pos %d, next bytes: %x", tag, r.pos(), next)
	}
}
