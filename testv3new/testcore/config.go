package testcore

// config here

// system
var (
	// TESTIP       = "59.36.173.89"
	TESTIP       = "203.56.175.233"
	APIADDR      = "http://" + TESTIP + ":10002"
	WSADDR       = "ws://" + TESTIP + ":10001"
	REGISTERADDR = APIADDR + "/auth/user_register"
	TOKENADDR    = APIADDR + "/auth/user_token"
)
