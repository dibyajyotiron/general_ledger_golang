package util

import (
	"encoding/json"
	"fmt"
	"io"

	setting "general_ledger_golang/pkg/config"
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
		return nil
	}

	err = json.Unmarshal(bodyBytes, &x)

	if err != nil {
		fmt.Printf("Req.Body parsing failed inside UseRequestBody, error: %+v", err)
		return nil
	}

	return x
}

func StructToJSON(m interface{}) map[string]interface{} {
	resultStr, _ := json.Marshal(m)
	var result map[string]interface{}
	err := json.Unmarshal(resultStr, &result)
	if err != nil {
		fmt.Printf("Unable to convert model to Struct %+v", err)
		return nil
	}
	return result
}

type CopyableMap map[string]interface{}
type CopyableSlice []interface{}

// DeepCopyMap will create a deep copy of the given map. The depth of this
// copy is all-inclusive. Both maps and slices will be considered when
// making the copy.
func DeepCopyMap(m map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}

	for k, v := range m {
		// Handle maps
		mapvalue, isMap := v.(map[string]interface{})
		if isMap {
			result[k] = DeepCopyMap(mapvalue)
			continue
		}

		// Handle slices
		sliceValue, isSlice := v.([]interface{})
		if isSlice {
			result[k] = DeepCopySlice(sliceValue)
			continue
		}

		result[k] = v
	}

	return result
}

// DeepCopySlice will create a deep copy of this slice. The depth of this
// copy is all-inclusive. Both maps and slices will be considered when
// making the copy.
func DeepCopySlice(s []interface{}) []interface{} {
	var result []interface{}

	for _, v := range s {
		// Handle maps
		mapvalue, isMap := v.(map[string]interface{})
		if isMap {
			result = append(result, DeepCopyMap(CopyableMap(mapvalue)))
			continue
		}

		// Handle slices
		sliceValue, isSlice := v.([]interface{})
		if isSlice {
			result = append(result, DeepCopySlice(CopyableSlice(sliceValue)))
			continue
		}

		result = append(result, v)
	}

	return result
}
