package main

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
)

type codegenTypeRefMode int

const (
	typeRefLocal codegenTypeRefMode = iota
	typeRefQualified
)

type typeRefConfig struct {
	mode            codegenTypeRefMode
	typesImportPath string
}

func generate(
	s *TLOSchema,
	typesPackageName string,
	functionsPackageName string,
	typesImportPath string,
) (map[string]string, map[string]string, error) {
	files := map[string]*jen.File{}
	fileFor := func(key string) *jen.File {
		if key == "" {
			key = "common"
		}
		f, ok := files[key]
		if ok {
			return f
		}
		f = jen.NewFile(typesPackageName)
		files[key] = f
		return f
	}
	typeRefLocalCfg := typeRefConfig{mode: typeRefLocal}

	// Interfaces + structs for each abstract type
	for _, t := range s.Types {
		if t.IsBuiltinType() {
			continue
		}
		if len(t.Constructors) == 0 {
			continue
		}
		f := fileFor(fileKeyFromQualifiedName(t.Name))
		if len(t.Constructors) > 1 {
			f.Add(genInterface(t))
			f.Add(genInterfaceDecoderHelper(t))
			f.Line()
		}
		for _, c := range t.Constructors {
			f.Add(genStruct(c, s, len(t.Constructors) > 1, typeRefLocalCfg))
			f.Line()
		}
	}

	functionsFiles := map[string]*jen.File{}
	functionsByFileKey := map[string][]*TLOCombinator{}
	prefixByFileKey := map[string]string{}
	functionFileFor := func(key string) *jen.File {
		if key == "" {
			key = "common"
		}
		f, ok := functionsFiles[key]
		if ok {
			return f
		}
		f = jen.NewFile(functionsPackageName)
		functionsFiles[key] = f
		return f
	}
	typeRefFunctionsCfg := typeRefConfig{mode: typeRefQualified, typesImportPath: typesImportPath}

	// Functions
	if len(s.Functions) > 0 {
		commonFile := functionFileFor("common")
		commonFile.Add(genRequesterInterface())
		commonFile.Add(genRequestFunc())
		for _, fn := range s.Functions {
			prefix, key := functionPrefixAndFileKey(fn.Name)
			f := functionFileFor(key)
			f.Add(genFunctionType(fn, s, typeRefFunctionsCfg))
			f.Line()
			functionsByFileKey[key] = append(functionsByFileKey[key], fn)
			if _, ok := prefixByFileKey[key]; !ok {
				prefixByFileKey[key] = prefix
			}
		}
		for key, fns := range functionsByFileKey {
			functionFileFor(key).Add(genFunctionWrapper(prefixByFileKey[key], fns, s, typeRefFunctionsCfg))
		}
	}

	typesOut := make(map[string]string, len(files))
	for name, f := range files {
		typesOut[name+".go"] = fmt.Sprintf("%#v", f)
	}
	functionsOut := make(map[string]string, len(functionsFiles))
	for name, f := range functionsFiles {
		functionsOut[name+".go"] = fmt.Sprintf("%#v", f)
	}
	return typesOut, functionsOut, nil
}

// --- naming helpers ----------------------------------------------------------

var capitalizePatterns = []string{"api", "id", "json", "p2p", "sha", "srp", "ttl", "url"}

func goName(name string) string {
	if name == "" {
		return ""
	}
	// "client.runQuery" → "ClientRunQuery"
	parts := strings.SplitN(name, ".", 2)
	res := ""
	for _, p := range parts {
		delim := strcase.ToDelimited(p, '|')
		for _, item := range strings.Split(delim, "|") {
			item = strings.ToLower(item)
			upper := false
			for _, pat := range capitalizePatterns {
				if item == pat {
					item = strings.ToUpper(item)
					upper = true
					break
				}
			}
			if !upper {
				runes := []rune(item)
				if len(runes) == 0 {
					continue
				}
				runes[0] = unicode.ToUpper(runes[0])
				item = string(runes)
			}
			res += item
		}
	}
	return res
}

func ifaceName(t *TLOType) string { return "I" + goName(t.Name) }

func fileKeyFromQualifiedName(name string) string {
	_, key := functionPrefixAndFileKey(name)
	return key
}

func functionPrefixAndFileKey(name string) (string, string) {
	if name == "" {
		return "", "common"
	}
	prefix, _, ok := strings.Cut(name, ".")
	if !ok || prefix == "" {
		return "", "common"
	}
	return prefix, strings.ToLower(prefix)
}

// --- type expression → jen.Code ---------------------------------------------

func typeExprCode(node *TypeExprNode, s *TLOSchema, cfg typeRefConfig) jen.Code {
	if node == nil {
		return jen.Interface()
	}
	switch node.Kind {
	case "nat_const", "nat_var", "type_var":
		// flags / nat fields rendered as uint32
		return jen.Uint32()

	case "type_expr":
		return builtinOrUserType(node, s, cfg)

	case "array":
		if len(node.ArrayArgs) == 1 {
			return jen.Index().Add(typeExprCode(node.ArrayArgs[0].TypeExpr, s, cfg))
		}
		return jen.Index().Interface()

	default:
		return jen.Interface()
	}
}

func builtinOrUserType(node *TypeExprNode, s *TLOSchema, cfg typeRefConfig) jen.Code {
	makeTypeRef := func(name string) jen.Code {
		if cfg.mode == typeRefQualified {
			return jen.Qual(cfg.typesImportPath, name)
		}
		return jen.Id(name)
	}
	switch node.TypeID {
	case builtinNat:
		return jen.Uint32()
	case builtinInt:
		return jen.Int32()
	case builtinLong:
		return jen.Int64()
	case builtinDouble:
		return jen.Float64()
	case builtinString:
		return jen.String()
	case builtinBoolFalse, builtinBoolTrue:
		return jen.Bool()
	case builtinVector:
		// Vector<T>: one child type param
		if len(node.Children) == 1 {
			return jen.Index().Add(typeExprCode(node.Children[0], s, cfg))
		}
		return jen.Index().Interface()
	}

	// Look up user type
	t, ok := s.TypesByID[node.TypeID]
	if !ok {
		return jen.Interface() // unknown
	}
	switch t.Name {
	case "Int128":
		return jen.Index(jen.Lit(16)).Byte()
	case "Int256":
		return jen.Index(jen.Lit(32)).Byte()
	case "Bytes":
		return jen.Index().Byte()
	case "Bool", "True":
		return jen.Bool()
	}

	if len(t.Constructors) == 0 {
		return jen.Interface()
	}
	if len(t.Constructors) > 1 {
		return makeTypeRef(ifaceName(t))
	}
	return makeTypeRef(goName(t.Constructors[0].Name))
}

// --- code generators ---------------------------------------------------------

func genInterface(t *TLOType) jen.Code {
	iface := jen.Type().Id(ifaceName(t)).Interface(
		jen.Id("CRC").Params().Uint32(),
		jen.Id("MarshalTL").Params().Params(jen.Index().Byte(), jen.Error()),
		jen.Id("UnmarshalTL").Params(jen.Qual("io", "Reader")).Error(),
		jen.Id("_"+ifaceName(t)).Params(),
	)
	checks := make([]jen.Code, len(t.Constructors))
	for i, c := range t.Constructors {
		checks[i] = jen.Id("_").Id(ifaceName(t)).Op("=").
			Call(jen.Op("*").Id(goName(c.Name))).Call(jen.Nil())
	}
	return jen.Add(iface, jen.Line(), jen.Var().Defs(checks...).Line())
}

func genInterfaceDecoderHelper(t *TLOType) jen.Code {
	if len(t.Constructors) <= 1 {
		return jen.Null()
	}
	ifName := ifaceName(t)
	localName := "decode" + ifName
	pubName := "Decode" + ifName

	local := jen.Func().Id(localName).
		Params(jen.Id("r").Qual("io", "Reader")).
		Params(jen.Id(ifName), jen.Error()).
		Block(
			jen.Var().Id("tag").Uint32(),
			jen.Err().Op(":=").Qual("github.com/tonkeeper/tongo/tl", "Unmarshal").Call(jen.Id("r"), jen.Op("&").Id("tag")),
			jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Nil(), jen.Err())),
			jen.Var().Id("res").Id(ifName),
			jen.Switch(jen.Id("tag")).BlockFunc(func(g *jen.Group) {
				for _, c := range t.Constructors {
					g.Case(jen.Lit(c.ID)).Block(
						jen.Id("res").Op("=").Op("&").Id(goName(c.Name)).Values(),
					)
				}
				g.Default().Block(
					jen.Return(
						jen.Nil(),
						jen.Qual("fmt", "Errorf").Call(jen.Lit("invalid crc code: got 0x%08x"), jen.Id("tag")),
					),
				)
			}),
			jen.Err().Op("=").Id("res").Dot("UnmarshalTL").Call(jen.Id("r")),
			jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Nil(), jen.Err())),
			jen.Return(jen.Id("res"), jen.Nil()),
		)

	pub := jen.Func().Id(pubName).
		Params(jen.Id("r").Qual("io", "Reader")).
		Params(jen.Id(ifName), jen.Error()).
		Block(
			jen.Return(jen.Id(localName).Call(jen.Id("r"))),
		)
	return jen.Add(local, jen.Line(), pub)
}

func genStruct(c *TLOCombinator, s *TLOSchema, hasIface bool, cfg typeRefConfig) jen.Code {
	typName := goName(c.Name)
	fields := make([]jen.Code, 0, len(c.Args))
	flagFieldByVarNum := map[int32]string{}
	for _, arg := range c.Args {
		if isFlagsArg(arg, s) {
			flagName := goName(arg.Name)
			if flagName == "" {
				flagName = fmt.Sprintf("Flags%d", arg.VarNum)
			}
			fields = append(fields, jen.Id(flagName).Uint32().Tag(map[string]string{"tl": bitflagTag(arg)}))
			flagFieldByVarNum[arg.VarNum] = flagName
			continue
		}
		fieldName := goName(arg.Name)
		fieldType := typeExprCode(arg.TypeExpr, s, cfg)
		if arg.IsOptional && shouldPointerWrapOptional(arg.TypeExpr, s) {
			fieldType = jen.Op("*").Add(fieldType)
		}
		if arg.IsOptional {
			target := flagFieldByVarNum[arg.ExistVarNum]
			if target == "" {
				target = fmt.Sprintf("%d", arg.ExistVarNum)
			}
			tag := fmt.Sprintf(",omitempty:%s:%d", target, arg.ExistVarBit)
			fields = append(fields, jen.Id(fieldName).Add(fieldType).Tag(map[string]string{"tl": tag}))
		} else {
			fields = append(fields, jen.Id(fieldName).Add(fieldType))
		}
	}

	typ := jen.Type().Id(typName).Struct(fields...)
	crc := jen.Func().Params(jen.Op("*").Id(typName)).Id("CRC").Params().Uint32().
		Block(jen.Return(jen.Lit(c.ID)))

	marshalTL := genMarshalTLMethod(typName, c.Args, s, flagFieldByVarNum)
	unmarshalTL := genUnmarshalTLMethod(typName, c.Args, s, flagFieldByVarNum)
	res := jen.Add(typ, jen.Line(), crc, jen.Line(), marshalTL, jen.Line(), unmarshalTL)
	if hasIface && c.ResultType != nil {
		marker := jen.Func().Params(jen.Op("*").Id(typName)).
			Id("_" + ifaceName(c.ResultType)).Params().Block()
		res = res.Add(jen.Line(), marker)
	}
	return res
}

func genRequesterInterface() jen.Code {
	return jen.Type().Id("Requester").Interface(
		jen.Id("MakeRequest").
			Params(jen.Id("ctx").Qual("context", "Context"), jen.Id("msg").Index().Byte()).
			Params(jen.Index().Byte(), jen.Error()),
	).Line()
}

func genRequestFunc() jen.Code {
	raw := jen.Func().Id("requestRaw").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("m").Id("Requester"),
			jen.Id("in").Interface(
				jen.Id("CRC").Params().Uint32(),
				jen.Id("MarshalTL").Params().Params(jen.Index().Byte(), jen.Error()),
			),
		).
		Params(jen.Index().Byte(), jen.Error()).
		Block(
			jen.List(jen.Id("body"), jen.Err()).Op(":=").Id("in").Dot("MarshalTL").Call(),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("marshaling: %w"), jen.Err())),
			),
			jen.Id("msg").Op(":=").Make(jen.Index().Byte(), jen.Lit(4).Op("+").Len(jen.Id("body"))),
			jen.Qual("encoding/binary", "LittleEndian").Dot("PutUint32").Call(jen.Id("msg"), jen.Id("in").Dot("CRC").Call()),
			jen.Id("copy").Call(jen.Id("msg").Index(jen.Lit(4), jen.Empty()), jen.Id("body")),
			jen.List(jen.Id("respRaw"), jen.Err()).Op(":=").Id("m").Dot("MakeRequest").Call(jen.Id("ctx"), jen.Id("msg")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("sending: %w"), jen.Err())),
			),
			jen.Return(jen.Id("respRaw"), jen.Nil()),
		).Line()

	normal := jen.Func().Id("request").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("m").Id("Requester"),
			jen.Id("in").Interface(
				jen.Id("CRC").Params().Uint32(),
				jen.Id("MarshalTL").Params().Params(jen.Index().Byte(), jen.Error()),
			),
			jen.Id("out").Interface(
				jen.Id("CRC").Params().Uint32(),
				jen.Id("UnmarshalTL").Params(jen.Qual("io", "Reader")).Error(),
			),
		).
		Params(jen.Error()).
		Block(
			jen.List(jen.Id("respRaw"), jen.Err()).Op(":=").Id("requestRaw").Call(jen.Id("ctx"), jen.Id("m"), jen.Id("in")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Err()),
			),
			jen.If(jen.Len(jen.Id("respRaw")).Op("<").Lit(4)).Block(
				jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("response: too short: %d"), jen.Len(jen.Id("respRaw")))),
			),
			jen.Id("got").Op(":=").Qual("encoding/binary", "LittleEndian").Dot("Uint32").Call(jen.Id("respRaw")),
			jen.Id("want").Op(":=").Id("out").Dot("CRC").Call(),
			jen.If(jen.Id("got").Op("!=").Id("want")).Block(
				jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("response: invalid crc code: got 0x%08x; want 0x%08x"), jen.Id("got"), jen.Id("want"))),
			),
			jen.Err().Op("=").Id("out").Dot("UnmarshalTL").Call(jen.Qual("bytes", "NewReader").Call(jen.Id("respRaw").Index(jen.Lit(4), jen.Empty()))),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("response: %w"), jen.Err())),
			),
			jen.Return(jen.Nil()),
		).Line()
	return jen.Add(raw, jen.Line(), normal)
}

func genFunctionType(fn *TLOCombinator, s *TLOSchema, cfg typeRefConfig) jen.Code {
	makeTypeRef := func(name string) jen.Code {
		if cfg.mode == typeRefQualified {
			return jen.Qual(cfg.typesImportPath, name)
		}
		return jen.Id(name)
	}
	reqName := goName(fn.Name) + "Request"
	fields := make([]jen.Code, 0, len(fn.Args))
	flagFieldByVarNum := map[int32]string{}
	for _, arg := range fn.Args {
		if isFlagsArg(arg, s) {
			flagName := goName(arg.Name)
			if flagName == "" {
				flagName = fmt.Sprintf("Flags%d", arg.VarNum)
			}
			fields = append(fields, jen.Id(flagName).Uint32().Tag(map[string]string{"tl": bitflagTag(arg)}))
			flagFieldByVarNum[arg.VarNum] = flagName
			continue
		}
		fieldType := typeExprCode(arg.TypeExpr, s, cfg)
		if arg.IsOptional && shouldPointerWrapOptional(arg.TypeExpr, s) {
			fieldType = jen.Op("*").Add(fieldType)
		}
		field := jen.Id(goName(arg.Name)).Add(fieldType)
		if arg.IsOptional {
			target := flagFieldByVarNum[arg.ExistVarNum]
			if target == "" {
				target = fmt.Sprintf("%d", arg.ExistVarNum)
			}
			tag := fmt.Sprintf(",omitempty:%s:%d", target, arg.ExistVarBit)
			field = field.Tag(map[string]string{"tl": tag})
		}
		fields = append(fields, field)
	}

	reqStruct := jen.Type().Id(reqName).Struct(fields...)
	crc := jen.Func().Params(jen.Op("*").Id(reqName)).Id("CRC").Params().Uint32().
		Block(jen.Return(jen.Lit(fn.ID)))
	marshalTL := genMarshalTLMethod(reqName, fn.Args, s, flagFieldByVarNum)

	var retType jen.Code
	if fn.ResultType != nil && len(fn.ResultType.Constructors) > 1 {
		retType = makeTypeRef(ifaceName(fn.ResultType))
	} else if fn.ResultType != nil && len(fn.ResultType.Constructors) == 1 {
		retType = makeTypeRef(goName(fn.ResultType.Constructors[0].Name))
	} else {
		retType = jen.Interface()
	}

	fnCode := jen.Func().Id(goName(fn.Name)).
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("m").Id("Requester"),
			jen.Id("i").Id(reqName),
		).
		Params(retType, jen.Error()).
		BlockFunc(func(g *jen.Group) {
			if fn.ResultType != nil && len(fn.ResultType.Constructors) > 1 {
				g.List(jen.Id("respRaw"), jen.Err()).Op(":=").Id("requestRaw").Call(jen.Id("ctx"), jen.Id("m"), jen.Op("&").Id("i"))
				g.If(jen.Err().Op("!=").Nil()).Block(
					jen.Var().Id("zero").Add(retType),
					jen.Return(jen.Id("zero"), jen.Err()),
				)
				g.List(jen.Id("res"), jen.Err()).Op(":=").Qual(cfg.typesImportPath, "Decode"+ifaceName(fn.ResultType)).Call(
					jen.Qual("bytes", "NewReader").Call(jen.Id("respRaw")),
				)
				g.If(jen.Err().Op("!=").Nil()).Block(
					jen.Var().Id("zero").Add(retType),
					jen.Return(jen.Id("zero"), jen.Qual("fmt", "Errorf").Call(jen.Lit("response: %w"), jen.Err())),
				)
				g.Return(jen.Id("res"), jen.Nil())
				return
			}
			g.Var().Id("res").Add(retType)
			g.Return(jen.Id("res"), jen.Id("request").Call(
				jen.Id("ctx"), jen.Id("m"), jen.Op("&").Id("i"), jen.Op("&").Id("res"),
			))
		})

	return jen.Add(reqStruct, jen.Line(), crc, jen.Line(), marshalTL, jen.Line(), fnCode)
}

func genFunctionWrapper(prefix string, fns []*TLOCombinator, s *TLOSchema, cfg typeRefConfig) jen.Code {
	if prefix == "" || len(fns) == 0 {
		return jen.Null()
	}
	wrapperName := goName(prefix)
	iface := jen.Type().Id(wrapperName + "API").InterfaceFunc(func(g *jen.Group) {
		for _, fn := range fns {
			reqName, retType, methodName := fnRequestAndReturnAndMethod(fn, s, cfg, wrapperName)
			g.Id(methodName).
				Params(
					jen.Id("ctx").Qual("context", "Context"),
					jen.Id("i").Id(reqName),
				).
				Params(retType, jen.Error())
		}
	})

	wrapperType := jen.Type().Id(wrapperName).Struct(
		jen.Id("requester").Id("Requester"),
	)
	ctor := jen.Func().Id("New" + wrapperName).
		Params(jen.Id("requester").Id("Requester")).
		Params(jen.Op("*").Id(wrapperName)).
		Block(
			jen.Return(jen.Op("&").Id(wrapperName).Values(jen.Dict{
				jen.Id("requester"): jen.Id("requester"),
			})),
		)

	var methods []jen.Code
	for _, fn := range fns {
		reqName, retType, methodName := fnRequestAndReturnAndMethod(fn, s, cfg, wrapperName)
		method := jen.Func().
			Params(jen.Id("c").Op("*").Id(wrapperName)).
			Id(methodName).
			Params(
				jen.Id("ctx").Qual("context", "Context"),
				jen.Id("i").Id(reqName),
			).
			Params(retType, jen.Error()).
			Block(
				jen.Return(
					jen.Id(goName(fn.Name)).Call(
						jen.Id("ctx"),
						jen.Id("c").Dot("requester"),
						jen.Id("i"),
					),
				),
			)
		methods = append(methods, method, jen.Line())
	}

	implCheck := jen.Var().Id("_").Id(wrapperName + "API").Op("=").Call(jen.Op("*").Id(wrapperName)).Call(jen.Nil())
	return jen.Add(iface, jen.Line(), wrapperType, jen.Line(), ctor, jen.Line(), jen.Add(methods...), implCheck, jen.Line())
}

func fnRequestAndReturnAndMethod(fn *TLOCombinator, s *TLOSchema, cfg typeRefConfig, wrapperName string) (string, jen.Code, string) {
	makeTypeRef := func(name string) jen.Code {
		if cfg.mode == typeRefQualified {
			return jen.Qual(cfg.typesImportPath, name)
		}
		return jen.Id(name)
	}
	reqName := goName(fn.Name) + "Request"

	var retType jen.Code
	if fn.ResultType != nil && len(fn.ResultType.Constructors) > 1 {
		retType = makeTypeRef(ifaceName(fn.ResultType))
	} else if fn.ResultType != nil && len(fn.ResultType.Constructors) == 1 {
		retType = makeTypeRef(goName(fn.ResultType.Constructors[0].Name))
	} else {
		retType = jen.Interface()
	}

	fnName := goName(fn.Name)
	methodName := strings.TrimPrefix(fnName, wrapperName)
	if methodName == "" {
		methodName = fnName
	}
	return reqName, retType, methodName
}

// IsBuiltinType returns true for types with no user-facing constructors.
func (t *TLOType) IsBuiltinType() bool {
	switch t.ID {
	case builtinNat, builtinInt, builtinLong, builtinDouble, builtinString, builtinVector,
		builtinBoolFalse, builtinBoolTrue:
		return true
	}
	switch t.Name {
	case "#", "Int128", "Int256", "Bytes", "Bool", "True":
		return true
	}
	return false
}

func isFlagsArg(arg *TLOArg, s *TLOSchema) bool {
	if arg == nil || arg.Name != "flags" || arg.TypeExpr == nil {
		return false
	}
	switch arg.TypeExpr.Kind {
	case "nat_const", "nat_var":
		return true
	case "type_expr":
		if arg.TypeExpr.TypeID == builtinNat {
			return true
		}
		t, ok := s.TypesByID[arg.TypeExpr.TypeID]
		return ok && t.Name == "#"
	default:
		return false
	}
}

func bitflagTag(arg *TLOArg) string {
	if arg != nil && arg.VarNum >= 0 {
		return fmt.Sprintf("%d,bitflag", arg.VarNum)
	}
	if arg != nil && arg.Name != "" {
		return arg.Name + ",bitflag"
	}
	return "flags,bitflag"
}

func genMarshalTLMethod(typeName string, args []*TLOArg, s *TLOSchema, flagFieldByVarNum map[int32]string) jen.Code {
	stmts := []jen.Code{
		jen.Var().Defs(
			jen.Err().Error(),
			jen.Id("b").Index().Byte(),
		),
		jen.Id("_").Op("=").Err(),
		jen.Id("_").Op("=").Id("b"),
		jen.Id("buf").Op(":=").Qual("bytes", "NewBuffer").Call(jen.Nil()),
	}

	// local mutable copies of bitflag fields, auto-populated from optional values.
	for varNum, fieldName := range flagFieldByVarNum {
		local := fmt.Sprintf("flagsVar%d", varNum)
		stmts = append(stmts, jen.Id(local).Op(":=").Id("t").Dot(fieldName))
	}

	for _, arg := range args {
		if isFlagsArg(arg, s) {
			continue
		}
		if !arg.IsOptional {
			continue
		}
		flagName := flagFieldByVarNum[arg.ExistVarNum]
		if flagName == "" {
			continue
		}
		local := fmt.Sprintf("flagsVar%d", arg.ExistVarNum)
		mask := uint32(1) << uint32(arg.ExistVarBit)
		field := jen.Id("t").Dot(goName(arg.Name))
		switch {
		case shouldPointerWrapOptional(arg.TypeExpr, s):
			stmts = append(stmts, jen.If(field.Op("!=").Nil()).Block(
				jen.Id(local).Op("|=").Lit(mask),
			))
		case arg.TypeExpr != nil && arg.TypeExpr.Kind == "array":
			stmts = append(stmts, jen.If(field.Op("!=").Nil()).Block(
				jen.Id(local).Op("|=").Lit(mask),
			))
		default:
			// Interfaces are nillable too.
			stmts = append(stmts, jen.If(field.Op("!=").Nil()).Block(
				jen.Id(local).Op("|=").Lit(mask),
			))
		}
	}

	for _, arg := range args {
		if isFlagsArg(arg, s) {
			local := fmt.Sprintf("flagsVar%d", arg.VarNum)
			stmts = append(stmts,
				jen.List(jen.Id("b"), jen.Err()).Op("=").Qual("github.com/tonkeeper/tongo/tl", "Marshal").Call(jen.Id(local)),
				jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Nil(), jen.Err())),
				jen.List(jen.Id("_"), jen.Err()).Op("=").Id("buf").Dot("Write").Call(jen.Id("b")),
				jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Nil(), jen.Err())),
			)
			continue
		}

		field := jen.Id("t").Dot(goName(arg.Name))
		if arg.IsOptional {
			local := fmt.Sprintf("flagsVar%d", arg.ExistVarNum)
			tmpBase := fmt.Sprintf("m%d%s", len(stmts), goName(arg.Name))
			stmts = append(stmts,
				jen.If(jen.Parens(jen.Id(local).Op(">>").Lit(arg.ExistVarBit)).Op("&").Lit(1).Op("==").Lit(1)).Block(
					marshalFieldStatements(arg.TypeExpr, field, s, shouldPointerWrapOptional(arg.TypeExpr, s), tmpBase)...,
				),
			)
			continue
		}

		tmpBase := fmt.Sprintf("m%d%s", len(stmts), goName(arg.Name))
		stmts = append(stmts, marshalFieldStatements(arg.TypeExpr, field, s, false, tmpBase)...)
	}

	stmts = append(stmts, jen.Return(jen.Id("buf").Dot("Bytes").Call(), jen.Nil()))

	return jen.Func().Params(jen.Id("t").Id(typeName)).Id("MarshalTL").Params().Params(jen.Index().Byte(), jen.Error()).Block(stmts...)
}

func genUnmarshalTLMethod(typeName string, args []*TLOArg, s *TLOSchema, flagFieldByVarNum map[int32]string) jen.Code {
	stmts := []jen.Code{
		jen.Var().Id("err").Error(),
		jen.Id("_").Op("=").Id("err"),
	}
	for _, arg := range args {
		field := jen.Id("t").Dot(goName(arg.Name))
		if isFlagsArg(arg, s) {
			stmts = append(stmts,
				jen.Err().Op("=").Qual("github.com/tonkeeper/tongo/tl", "Unmarshal").Call(jen.Id("r"), jen.Op("&").Add(field)),
				jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err())),
			)
			continue
		}

		tmpBase := fmt.Sprintf("tmp%d%s", len(stmts), goName(arg.Name))
		assignStmts := decodeFieldStatements(arg.TypeExpr, field, s, arg.IsOptional && shouldPointerWrapOptional(arg.TypeExpr, s), tmpBase)
		if arg.IsOptional {
			flagName := flagFieldByVarNum[arg.ExistVarNum]
			if flagName == "" {
				flagName = fmt.Sprintf("Flags%d", arg.ExistVarNum)
			}
			flagsField := jen.Id("t").Dot(flagName)
			stmts = append(stmts,
				jen.If(jen.Parens(flagsField.Op(">>").Lit(arg.ExistVarBit)).Op("&").Lit(1).Op("==").Lit(1)).Block(assignStmts...),
			)
			continue
		}
		stmts = append(stmts, assignStmts...)
	}
	stmts = append(stmts, jen.Return(jen.Nil()))
	return jen.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("UnmarshalTL").Params(jen.Id("r").Qual("io", "Reader")).Error().Block(stmts...)
}

func marshalFieldStatements(node *TypeExprNode, field *jen.Statement, s *TLOSchema, deref bool, tmpBase string) []jen.Code {
	if isDoubleType(node) {
		val := field
		if deref {
			val = jen.Op("*").Add(field)
		}
		bitsName := tmpBase + "Bits"
		rawName := tmpBase + "Raw"
		return []jen.Code{
			jen.Id(bitsName).Op(":=").Qual("math", "Float64bits").Call(val),
			jen.Var().Id(rawName).Index(jen.Lit(8)).Byte(),
			jen.Qual("encoding/binary", "LittleEndian").Dot("PutUint64").Call(jen.Id(rawName).Index(jen.Empty(), jen.Empty()), jen.Id(bitsName)),
			jen.List(jen.Id("_"), jen.Err()).Op("=").Id("buf").Dot("Write").Call(jen.Id(rawName).Index(jen.Empty(), jen.Empty())),
			jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Nil(), jen.Err())),
		}
	}
	if n, ok := fixedBytesLen(node, s); ok {
		data := field
		if deref {
			data = jen.Op("*").Add(field)
		}
		return []jen.Code{
			jen.List(jen.Id("_"), jen.Err()).Op("=").Id("buf").Dot("Write").Call(data.Index(jen.Empty(), jen.Empty())),
			jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Nil(), jen.Err())),
			jen.Id("_").Op("=").Lit(n),
		}
	}
	if it := interfaceTypeForExpr(node, s); it != nil {
		objName := tmpBase
		return []jen.Code{
			jen.Id(objName).Op(":=").Add(field),
			jen.If(jen.Id(objName).Op("==").Nil()).Block(
				jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("nil %s"), jen.Lit(goNameForExpr(node, s)))),
			),
			jen.List(jen.Id("b"), jen.Err()).Op("=").Qual("github.com/tonkeeper/tongo/tl", "Marshal").Call(jen.Id(objName).Dot("CRC").Call()),
			jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Nil(), jen.Err())),
			jen.List(jen.Id("_"), jen.Err()).Op("=").Id("buf").Dot("Write").Call(jen.Id("b")),
			jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Nil(), jen.Err())),
			jen.List(jen.Id("b"), jen.Err()).Op("=").Id(objName).Dot("MarshalTL").Call(),
			jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Nil(), jen.Err())),
			jen.List(jen.Id("_"), jen.Err()).Op("=").Id("buf").Dot("Write").Call(jen.Id("b")),
			jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Nil(), jen.Err())),
			jen.Id("_").Op("=").Lit(ifaceName(it)),
		}
	}
	marshalExpr := field
	if deref {
		marshalExpr = jen.Op("*").Add(field)
	}
	return []jen.Code{
		jen.List(jen.Id("b"), jen.Err()).Op("=").Qual("github.com/tonkeeper/tongo/tl", "Marshal").Call(marshalExpr),
		jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Nil(), jen.Err())),
		jen.List(jen.Id("_"), jen.Err()).Op("=").Id("buf").Dot("Write").Call(jen.Id("b")),
		jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Nil(), jen.Err())),
	}
}

func decodeFieldStatements(node *TypeExprNode, field *jen.Statement, s *TLOSchema, pointerWrap bool, tmpBase string) []jen.Code {
	if isDoubleType(node) {
		bitsName := tmpBase + "Bits"
		if pointerWrap {
			valName := tmpBase + "Val"
			return []jen.Code{
				jen.Var().Id(bitsName).Uint64(),
				jen.Err().Op("=").Qual("github.com/tonkeeper/tongo/tl", "Unmarshal").Call(jen.Id("r"), jen.Op("&").Id(bitsName)),
				jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err())),
				jen.Id(valName).Op(":=").Qual("math", "Float64frombits").Call(jen.Id(bitsName)),
				jen.Add(field).Op("=").Op("&").Id(valName),
			}
		}
		return []jen.Code{
			jen.Var().Id(bitsName).Uint64(),
			jen.Err().Op("=").Qual("github.com/tonkeeper/tongo/tl", "Unmarshal").Call(jen.Id("r"), jen.Op("&").Id(bitsName)),
			jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err())),
			jen.Add(field).Op("=").Qual("math", "Float64frombits").Call(jen.Id(bitsName)),
		}
	}
	if _, ok := fixedBytesLen(node, s); ok {
		if pointerWrap {
			t := typeExprCode(node, s, typeRefConfig{mode: typeRefLocal})
			return []jen.Code{
				jen.Var().Id(tmpBase).Add(t),
				jen.List(jen.Id("_"), jen.Err()).Op("=").Qual("io", "ReadFull").Call(jen.Id("r"), jen.Id(tmpBase).Index(jen.Empty(), jen.Empty())),
				jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err())),
				jen.Add(field).Op("=").Op("&").Id(tmpBase),
			}
		}
		return []jen.Code{
			jen.List(jen.Id("_"), jen.Err()).Op("=").Qual("io", "ReadFull").Call(jen.Id("r"), field.Index(jen.Empty(), jen.Empty())),
			jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err())),
		}
	}
	if it := interfaceTypeForExpr(node, s); it != nil {
		fn := "decode" + ifaceName(it)
		return []jen.Code{
			jen.List(jen.Id(tmpBase), jen.Err()).Op(":=").Id(fn).Call(jen.Id("r")),
			jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err())),
			jen.Add(field).Op("=").Id(tmpBase),
		}
	}
	if pointerWrap {
		t := typeExprCode(node, s, typeRefConfig{mode: typeRefLocal})
		return []jen.Code{
			jen.Var().Id(tmpBase).Add(t),
			jen.Err().Op("=").Qual("github.com/tonkeeper/tongo/tl", "Unmarshal").Call(jen.Id("r"), jen.Op("&").Id(tmpBase)),
			jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err())),
			jen.Add(field).Op("=").Op("&").Id(tmpBase),
		}
	}
	return []jen.Code{
		jen.Err().Op("=").Qual("github.com/tonkeeper/tongo/tl", "Unmarshal").Call(jen.Id("r"), jen.Op("&").Add(field)),
		jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err())),
	}
}

func fixedBytesLen(node *TypeExprNode, s *TLOSchema) (int, bool) {
	if node == nil || node.Kind != "type_expr" {
		return 0, false
	}
	t, ok := s.TypesByID[node.TypeID]
	if !ok {
		return 0, false
	}
	switch t.Name {
	case "Int128":
		return 16, true
	case "Int256":
		return 32, true
	default:
		return 0, false
	}
}

func interfaceTypeForExpr(node *TypeExprNode, s *TLOSchema) *TLOType {
	if node == nil || node.Kind != "type_expr" {
		return nil
	}
	switch node.TypeID {
	case builtinNat, builtinInt, builtinLong, builtinDouble, builtinString, builtinVector, builtinBoolFalse, builtinBoolTrue:
		return nil
	}
	t, ok := s.TypesByID[node.TypeID]
	if !ok || len(t.Constructors) <= 1 {
		return nil
	}
	switch t.Name {
	case "Bool", "True", "Int128", "Int256", "Bytes":
		return nil
	}
	return t
}

func goNameForExpr(node *TypeExprNode, s *TLOSchema) string {
	if node == nil || node.Kind != "type_expr" {
		return "field"
	}
	t, ok := s.TypesByID[node.TypeID]
	if !ok {
		return "field"
	}
	return goName(t.Name)
}

func isDoubleType(node *TypeExprNode) bool {
	return node != nil && node.Kind == "type_expr" && node.TypeID == builtinDouble
}

func shouldPointerWrapOptional(node *TypeExprNode, s *TLOSchema) bool {
	if node == nil {
		return false
	}
	switch node.Kind {
	case "array":
		return false // slices are nillable
	case "type_expr":
		switch node.TypeID {
		case builtinVector:
			return false // slices are nillable
		case builtinNat, builtinInt, builtinLong, builtinDouble, builtinString, builtinBoolFalse, builtinBoolTrue:
			return true
		}
		t, ok := s.TypesByID[node.TypeID]
		if !ok {
			return true
		}
		switch t.Name {
		case "Bool", "True":
			return true
		case "Bytes":
			return false // []byte
		}
		if len(t.Constructors) > 1 {
			return false // interface is nillable
		}
		return true
	default:
		return true
	}
}
