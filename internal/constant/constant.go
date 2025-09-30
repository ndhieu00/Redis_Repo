package constant

var RespOk = []byte("+OK\r\n")
var RespNil = []byte("$-1\r\n")
var TtlKeyNotExist = []byte(":-2\r\n")
var TtlKeyExistNoExpire = []byte(":-1\r\n")
