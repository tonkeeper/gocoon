// int ? = Int;
// long ? = Long;
// double ? = Double;
// string ? = String;
// object ? = Object;
// function ? = Function;
// bytes data:string = Bytes;
// true = True;
// boolTrue = Bool;
// boolFalse = Bool;

// vector {t:Type} # [ t ] = Vector t;

// int128 4*[ int ] = Int128;
// int256 8*[ int ] = Int256;

tokensUsed prompt_tokens_used:long cached_tokens_used:long completion_tokens_used:long reasoning_tokens_used:long
           total_tokens_used:long = TokensUsed;

proxy.params#d5c5609f flags:#
  proxy_public_key:int256
  proxy_owner_address:string
  proxy_sc_address:string
  is_test:flags.0?Bool
  proto_version:flags.1?int
  = proxy.Params;

client.params#40fdca64 flags:#
  client_owner_address:string
  is_test:flags.0?Bool
  min_proto_version:flags.1?int
  max_proto_version:flags.1?int
  = client.Params;

proxy.signedPayment data:bytes = proxy.SignedPayment;
proxy.signedPaymentEmpty = proxy.SignedPayment;

client.proxyConnectionAuthShort secret_hash:int256 nonce:long = client.ProxyConnectionAuth;
client.proxyConnectionAuthLong nonce:long = client.ProxyConnectionAuth;

client.authorizationWithProxySuccess signed_payment:proxy.SignedPayment tokens_committed_to_db:long max_tokens:long = client.AuthorizationWithProxy;
client.authorizationWithProxyFailed error_code:int error:string = client.AuthorizationWithProxy;

client.connectedToProxy params:proxy.params client_sc_address:string auth:client.ProxyConnectionAuth signed_payment:proxy.SignedPayment = client.ConnectedToProxy;

client.queryAnswer answer:bytes is_completed:Bool request_id:int256 request_tokens_used:tokensUsed = client.QueryAnswer;
client.queryAnswerError error_code:int error:string request_id:int256 request_tokens_used:tokensUsed = client.QueryAnswer;

client.queryFinalInfo flags:# tokens_used:tokensUsed worker_debug:flags.0?string proxy_debug:flags.0?string proxy_start_time:flags.1?double proxy_end_time:flags.1?double worker_start_time:flags.1?double worker_end_time:flags.1?double = client.QueryFinalInfo;
client.queryAnswerEx request_id:int256 answer:bytes flags:# final_info:flags.0?client.queryFinalInfo = client.QueryAnswerEx;
client.queryAnswerErrorEx request_id:int256 error_code:int error:string flags:# final_info:flags.0?client.queryFinalInfo = client.QueryAnswerEx;
client.queryAnswerPartEx request_id:int256 answer:bytes flags:# final_info:flags.0?client.queryFinalInfo = client.QueryAnswerEx;

client.refund data:bytes = client.Refund;
client.refundRejected active_queries:long = client.Refund;

client.workerType name:string active_workers:int coefficient_min:int coefficient_bucket50:int coefficient_max:int = client.WorkerType;
client.workerTypes types:(vector client.workerType) = client.WorkerTypes;

client.workerInstanceV2#3ea93d00 flags:# coefficient:int active_requests:int max_active_requests:int = client.WorkerInstanceV2;
client.workerTypeV2 name:string workers:(vector client.workerInstanceV2) = client.WorkerTypeV2;
client.workerTypesV2 types:(vector client.workerTypeV2) = client.WorkerTypesV2;

client.paymentStatus signed_payment:proxy.SignedPayment db_tokens:long max_tokens:long = client.PaymentStatus;

---functions---

client.connectToProxy params:client.params min_config_version:int = client.ConnectedToProxy;
client.authorizeWithProxyShort data:bytes = client.AuthorizationWithProxy;
client.authorizeWithProxyLong = client.AuthorizationWithProxy;
client.runQuery model_name:string query:bytes max_coefficient:int max_tokens:int timeout:double request_id:int256 min_config_version:int = client.QueryAnswer;
client.runQueryEx model_name:string query:bytes max_coefficient:int max_tokens:int timeout:double request_id:int256 min_config_version:int flags:# enable_debug:flags.0?Bool = client.QueryAnswerEx;
client.requestRefund = client.Refund;
client.updatePaymentStatus = client.PaymentStatus;
client.getWorkerTypes = client.WorkerTypes;
client.getWorkerTypesV2 = client.WorkerTypesV2;
