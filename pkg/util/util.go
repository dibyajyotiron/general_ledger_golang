package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"

	setting "general_ledger_golang/pkg/config"
	"general_ledger_golang/pkg/logger"
)

// Setup Initialize the util
func Setup() {
	cfg := *setting.GetConfig()
	jwtSecret = []byte(cfg.AppSetting.JwtSecret)
}

func Includes(e interface{}, s []interface{}) bool {
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

func GetReqBodyFromCtx(c *gin.Context) map[string]interface{} {
	x := map[string]interface{}{}

	bodyBytes := c.MustGet("requestBodyBytes")
	err := json.Unmarshal(bodyBytes.([]byte), &x)

	if err != nil {
		fmt.Printf("Req.Body parsing failed inside GetReqBodyFromCtx, error: %+v\n", err)
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

// Copy copies all key/value pairs in src adding them to dst.
// When a key in src is already present in dst,
// The value in dst will be overwritten by the value associated
// with the key in src if `shouldReplace` is true.
func Copy(dest, src map[string]interface{}, shouldReplace bool) {
	for k, v := range src {
		if shouldReplace {
			dest[k] = v
		} else {
			var vS []interface{}
			if dest[k] != nil {
				vS = append(vS, v)
				dest[k] = vS
			} else {
				dest[k] = v
			}
		}
	}
}

// InterfaceToMapOfString Convert interface or map[string]interface to map[string]string
func InterfaceToMapOfString(src interface{}) (map[string]string, error) {
	byteSlice, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}
	res := map[string]string{}
	err = json.Unmarshal(byteSlice, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func ConvertToMapSlice(i interface{}) ([]interface{}, error) {
	v, ok := i.([]interface{})
	if !ok {
		logger.Logger.Errorf("There was a problem converting to a slice")
		return nil, errors.New("conversion to slice failed")
	}
	return v, nil
}
