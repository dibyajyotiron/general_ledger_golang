package array

import "os"

func containsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// IsEnv takes a set of env names, and compares with APP_ENV
//
// IsEnv(constant.LOCAL) will check if it's local env
//
// IsEnv(constant.STAGE) will check if it's stage env
//
func IsEnv(localEnvNames []string) bool {
	currentEnv := os.Getenv("APP_ENV")
	if ContainsString(localEnvNames, currentEnv) {
		return true
	}
	return false
}
