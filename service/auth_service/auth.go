package auth_service

import "general_ledger_golang/pkg/config"

type Auth struct {
	ServiceName  string
	ServiceToken string
}

type CheckType string

func (c *CheckType) String() string {
	if c == nil {
		return ""
	}
	return string(*c)
}

var (
	READ  CheckType = "READ"
	WRITE CheckType = "WRITE"
)

func (a *Auth) Check(token string, checkType CheckType, serviceName string) bool {
	conf := config.GetConfig()
	allowedTokens := conf.ServerSetting.ServiceTokenWhitelist

	for service := range allowedTokens {
		RWToken := allowedTokens[service]
		if serviceName == service {
			allowedToken := RWToken[checkType.String()]
			// If you provide write token and ask to read, allowed
			if checkType == READ && RWToken[WRITE.String()] == token {
				return true
			}
			// else, write token writes, read token reads.
			if token == allowedToken {
				return true
			}
		}
	}
	return false
}
