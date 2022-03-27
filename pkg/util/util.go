package util

import (
	"encoding/json"
	"fmt"
	setting "general_ledger_golang/pkg/config"
	"io"
)

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

func ParseReqBodyToMap(reqBody io.ReadCloser) map[string]interface{} {
	x := map[string]interface{}{}

	bodyBytes, err := io.ReadAll(reqBody)

	if err != nil {
		fmt.Printf("Req.Body reading failed from ParseReqBodyToMap, error: %+v", err)
		panic(err)
	}

	err = json.Unmarshal(bodyBytes, &x)

	if err != nil {
		fmt.Printf("Req.Body parsing failed inside UseRequestBody, error: %+v", err)
		panic(err)
	}

	return x
}
