package interval

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36 Edg/130.0.0.0"
)

var SessionId string
var (
	client = &http.Client{
		Timeout: 5 * time.Second,
	}
)

type ResponseBodyFormat = func([]byte) string

// GenerateCommonRequest 生成通用的Request请求体
func GenerateCommonRequest(method, url string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return request, err
	}
	request.Header.Set("User-Agent", userAgent)
	if SessionId != "" {
		request.Header.Set("Cookie", SessionId)
	}
	return request, nil
}

// ParseResponseByStruct 解析练习平台返回的响应数据 在响应数据为结构体的情况下使用
func ParseResponseByStruct[T any](resp *http.Response, format ResponseBodyFormat) (*T, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取请求响应数据失败, message: %s", err)
		return nil, err
	}
	originStringBody := format(body)
	value := new(T)
	if err = json.Unmarshal([]byte(originStringBody), value); err != nil {
		log.Printf("格式化请求响应数据失败, message: %s", err)
		return nil, err
	}
	return value, nil
}

// ParseResponseBySlice 解析练习平台的响应数据，在响应数据为数组的情况下使用
func ParseResponseBySlice[T any](resp *http.Response, format ResponseBodyFormat) ([]T, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取请求响应数据失败, message: %s", err)
		return nil, err
	}
	originStringBody := format(body)
	values := make([]T, 0)
	if err = json.Unmarshal([]byte(originStringBody), &values); err != nil {
		log.Printf("格式化请求响应数据失败, message: %s", err)
		return nil, err
	}
	return values, nil
}
