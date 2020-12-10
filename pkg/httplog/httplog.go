package httplog

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	mylog "github.com/captainlee1024/fast-gin/log"
)

// HTTPGET 带有日志信息日志信息
func HTTPGET(trace *mylog.TraceContext, urlString string, urlParams url.Values, msTimeout int,
	header http.Header) (*http.Response, []byte, error) {
	startTime := time.Now().UnixNano()
	client := http.Client{
		Timeout: time.Duration(msTimeout) * time.Millisecond,
	}

	urlString = AddGetDataToURL(urlString, urlParams)
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		mylog.Log.Warn(req.URL.Path, trace, mylog.DLTagHTTPFailed, map[string]interface{}{
			"url":       urlString,
			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
			"method":    "GET",
			"args":      urlParams,
			"err":       err.Error(),
		})
		return nil, nil, err
	}

	if len(header) > 0 {
		req.Header = header
	}

	req = addTrace2Header(req, trace)
	resp, err := client.Do(req)
	if err != nil {
		mylog.Log.Warn(req.URL.Path, trace, mylog.DLTagHTTPFailed, map[string]interface{}{
			"url":       urlString,
			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
			"method":    "GET",
			"args":      urlParams,
			"err":       err.Error(),
		})
		return nil, nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		mylog.Log.Warn(req.URL.Path, trace, mylog.DLTagHTTPFailed, map[string]interface{}{
			"url":       urlString,
			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
			"method":    "GET",
			"result":    Substr(string(body), 0, 1024),
			"err":       err.Error(),
		})
		return nil, nil, err
	}
	mylog.Log.Info("", trace, mylog.DLTagHTTPFailed, map[string]interface{}{
		"url":       urlString,
		"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
		"method":    "GET",
		"args":      urlParams,
		"err":       err.Error(),
	})
	return resp, body, nil
}

// HTTPPOST 带有日志的 POST 方法
func HTTPPOST(trace *mylog.TraceContext, urlString string, urlParams url.Values, msTimeout int, header http.Header, contextType string) (*http.Response, []byte, error) {
	startTime := time.Now().UnixNano()
	client := http.Client{
		Timeout: time.Duration(msTimeout) * time.Millisecond,
	}
	if contextType == "" {
		contextType = "application/x-www-form-urlencoded"
	}
	urlParamEncode := urlParams.Encode()
	req, err := http.NewRequest("POST", urlString, strings.NewReader(urlParamEncode))
	if len(header) > 0 {
		req.Header = header
	}
	req = addTrace2Header(req, trace)
	req.Header.Set("Content-Type", contextType)
	resp, err := client.Do(req)
	if err != nil {
		mylog.Log.Warn(req.URL.Path, trace, mylog.DLTagHTTPFailed, map[string]interface{}{
			"url":       urlString,
			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
			"method":    "POST",
			"args":      Substr(urlParamEncode, 0, 1024),
			"err":       err.Error(),
		})
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		mylog.Log.Warn(req.URL.Path, trace, mylog.DLTagHTTPFailed, map[string]interface{}{
			"url":       urlString,
			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
			"method":    "POST",
			"args":      Substr(urlParamEncode, 0, 1024),
			"result":    Substr(string(body), 0, 1024),
			"err":       err.Error(),
		})
		return nil, nil, err
	}
	mylog.Log.Info(req.URL.Path, trace, mylog.DLTagHTTPSuccess, map[string]interface{}{
		"url":       urlString,
		"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
		"method":    "POST",
		"args":      Substr(urlParamEncode, 0, 1024),
		"result":    Substr(string(body), 0, 1024),
	})
	return resp, body, nil
}

// HTTPJSON 带有日志的 JSON 方法
func HTTPJSON(trace *mylog.TraceContext, urlString string, jsonContent string, msTimeout int, header http.Header) (*http.Response, []byte, error) {
	startTime := time.Now().UnixNano()
	client := http.Client{
		Timeout: time.Duration(msTimeout) * time.Millisecond,
	}
	req, err := http.NewRequest("POST", urlString, strings.NewReader(jsonContent))
	if len(header) > 0 {
		req.Header = header
	}
	req = addTrace2Header(req, trace)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		mylog.Log.Warn(req.URL.Path, trace, mylog.DLTagHTTPFailed, map[string]interface{}{
			"url":       urlString,
			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
			"method":    "POST",
			"args":      Substr(jsonContent, 0, 1024),
			"err":       err.Error(),
		})
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		mylog.Log.Warn(req.URL.Path, trace, mylog.DLTagHTTPFailed, map[string]interface{}{
			"url":       urlString,
			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
			"method":    "POST",
			"args":      Substr(jsonContent, 0, 1024),
			"result":    Substr(string(body), 0, 1024),
			"err":       err.Error(),
		})
		return nil, nil, err
	}
	mylog.Log.Info(req.URL.Path, trace, mylog.DLTagHTTPSuccess, map[string]interface{}{
		"url":       urlString,
		"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
		"method":    "POST",
		"args":      Substr(jsonContent, 0, 1024),
		"result":    Substr(string(body), 0, 1024),
	})
	return resp, body, nil
}

// description
func addTrace2Header(request *http.Request, trace *mylog.TraceContext) *http.Request {
	traceID := trace.Trace.TraceID
	cSpanID := mylog.NewSpanID()
	if traceID != "" {
		request.Header.Set("didi-header-traceID", traceID)
	}
	if cSpanID != "" {
		request.Header.Set("didi-header-spanid", cSpanID)
	}
	trace.CSpanID = cSpanID
	return request
}

// AddGetDataToURL xxx
func AddGetDataToURL(urlString string, data url.Values) string {
	if strings.Contains(urlString, "?") {
		urlString = urlString + "&"
	} else {
		urlString = urlString + "?"
	}
	return fmt.Sprintf("%s%s", urlString, data.Encode())
}

// Substr 截取字符串
func Substr(str string, start int64, end int64) string {
	length := int64(len(str))
	if start < 0 || start > length || end < 0 {
		return ""
	}

	if end > length {
		end = length
	}
	return string(str[start:end])
}
