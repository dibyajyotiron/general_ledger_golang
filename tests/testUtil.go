package tests

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type SM []map[string]interface{}

type TestRequestBuilder struct {
	URL    string
	METHOD string
	BODY   string
}

func (t *TestRequestBuilder) DepositOrWithdraw() error {
	_, err := t.MakeApiCall()

	if err != nil {
		return err
	}
	return nil
}

func (t *TestRequestBuilder) MakeApiCall() (*string, error) {
	method, url, reqBodyAsStr := t.METHOD, t.URL, t.BODY
	payload := strings.NewReader(reqBodyAsStr)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	strResp := string(body)
	return &strResp, nil
}

func (t *TestRequestBuilder) MapToJSONString(object map[string]interface{}) string {
	bodyBytes, _ := json.Marshal(object)
	return string(bodyBytes)
}

// JSONStringToMap returns a map version of a stringified json, in case of error, returns nil
func (t *TestRequestBuilder) JSONStringToMap(object string) map[string]interface{} {
	m := map[string]interface{}{}
	err := json.Unmarshal([]byte(object), &m)
	if err != nil {
		return nil
	}
	return m
}
