package util

import setting "general_ledger_golang/pkg/config"

// Setup Initialize the util
func Setup() {
	jwtSecret = []byte(setting.AppSetting.JwtSecret)
}

func Contains(e interface{}, s []string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
