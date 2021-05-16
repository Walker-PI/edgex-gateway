package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// GetMarshalStr 格式化对象输出（json序列化）
func GetMarshalStr(obj interface{}) string {
	if obj == nil {
		return ""
	}
	bytes, err := json.Marshal(obj)
	if err != nil {
		return fmt.Sprintf("Marshal Error. obj=%v err=%v", obj, err)
	}
	return string(bytes)
}

// GetBodyStr 获取请求Body
func GetBodyStr(r *http.Request) string {
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return ""
	}
	buf := bytes.NewBuffer(body)
	r.Body = ioutil.NopCloser(buf)
	return string(body)
}
