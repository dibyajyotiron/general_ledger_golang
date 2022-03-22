package e

var MsgFlags = map[int]string{
	SUCCESS:             "SUCCESS",
	DEBIT:               "DEBIT",
	CREDIT:              "CREDIT",
	ERROR:               "ERROR",
	MISSING_AUTH_HEADER: "MISSING_AUTH_HEADER",
	INVALID_PARAMS:      "INVALID_PARAMS",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
