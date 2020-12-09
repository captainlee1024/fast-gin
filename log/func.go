package log

// // HTTPGET 封装日志信息
// func HTTPGET(trace *TraceContext, urlString string, urlParams url.Values, msTimeout int,
// 	header http.Header) (*http.Response, []byte, error) {
// 	startTime := time.Now().UnixNano()
// 	client := http.Client{
// 		Timeout: time.Duration(msTimeout) * time.Millisecond,
// 	}

// 	urlString = AddGetDataToURL(urlString, urlParams)
// 	req, err := http.NewRequest("GET", urlString, nil)
// 	if err != nil {
// 		Log.Warn(trace, DLTagHTTPFailed, map[string]interface{}{
// 			"url":       urlString,
// 			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
// 			"method":    "GET",
// 			"args":      urlParams,
// 			"err":       err.Error(),
// 		})
// 		return nil, nil, err
// 	}

// 	if len(header) > 0 {
// 		req.Header = header
// 	}

// 	req = addTrace2Header(req, trace)
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		Log.Warn(trace, DLTagHTTPFailed, map[string]interface{}{
// 			"url":       urlString,
// 			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
// 			"method":    "GET",
// 			"args":      urlParams,
// 			"err":       err.Error(),
// 		})
// 		return nil, nil, err
// 	}
// 	body, err := ioutil.ReadAll(resp.Body)
// 	defer resp.Body.Close()
// 	if err != nil {
// 		Log.Warn(trace, DLTagHTTPFailed, map[string]interface{}{
// 			"url":       urlString,
// 			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
// 			"method":    "GET",
// 			"result":    Substr(string(body), 0, 1024),
// 			"err":       err.Error(),
// 		})
// 		return nil, nil, err
// 	}
// 	Log.Info(trace, DLTagHTTPFailed, map[string]interface{}{
// 		"url":       urlString,
// 		"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
// 		"method":    "GET",
// 		"args":      urlParams,
// 		"err":       err.Error(),
// 	})
// 	return resp, body, nil
// }

// // description
// func addTrace2Header(request *http.Request, trace *TraceContext) *http.Request {
// 	traceID := trace.Trace.TraceID
// 	cSpanID := NewSpanID()
// 	if traceID != "" {
// 		request.Header.Set("didi-header-traceID", traceID)
// 	}
// 	if cSpanID != "" {
// 		request.Header.Set("didi-header-spanid", cSpanID)
// 	}
// 	trace.CSpanID = cSpanID
// 	return request
// }

// // GetMd5Hash description
// func GetMd5Hash(text string) string {
// 	hasher := md5.New()
// 	hasher.Write([]byte(text))
// 	return hex.EncodeToString(hasher.Sum(nil))
// }

// // Encode description
// func Encode(data string) (string, error) {
// 	h := md5.New()
// 	_, err := h.Write([]byte(data))
// 	if err != nil {
// 		return "", nil
// 	}
// 	return hex.EncodeToString(h.Sum(nil)), nil
// }

// // NewTrace 创建 TraceContext 并生成 TraceID SpandID
// func NewTrace() *TraceContext {
// 	trace := &TraceContext{}
// 	trace.Trace.TraceID = GetTraceID()
// 	trace.SpandID = NewSpanID()
// 	return trace
// }

// // NewSpanID description
// func NewSpanID() string {
// 	timestamp := uint32(time.Now().Unix())
// 	ipToLong := binary.BigEndian.Uint32(settings.LocalIP.To4())
// 	b := bytes.Buffer{}
// 	b.WriteString(fmt.Sprintf("%08x", ipToLong^timestamp))
// 	b.WriteString(fmt.Sprintf("%08x", rand.Int31()))
// 	return b.String()
// }

// // GetTraceID 创建并获取 TraceID
// func GetTraceID() (traceID string) {
// 	return calcTraceID(settings.LocalIP.String())
// }

// // 生成 traceID
// func calcTraceID(ip string) (trace string) {
// 	now := time.Now()
// 	timestamp := uint32(now.Unix())
// 	timeNano := now.UnixNano()
// 	pid := os.Getpid()

// 	b := bytes.Buffer{}
// 	netIP := net.ParseIP(ip)
// 	if netIP != nil {
// 		b.WriteString("00000000")
// 	} else {
// 		b.WriteString(hex.EncodeToString(netIP.To4()))
// 	}
// 	b.WriteString(fmt.Sprintf("%08x", timestamp&0xffffffff))
// 	b.WriteString(fmt.Sprintf("%04x", timeNano&0xffff))
// 	b.WriteString(fmt.Sprintf("%04x", pid&0xffff))
// 	b.WriteString(fmt.Sprintf("%06x", rand.Int31n(1<<24)))

// 	return b.String()
// }

// // GetLocalIPs 获取 IP 列表
// func GetLocalIPs() (ips []net.IP) {
// 	interfaceAddr, err := net.InterfaceAddrs()
// 	if err != nil {
// 		return nil
// 	}
// 	for _, address := range interfaceAddr {
// 		ipNet, isValidIPNew := address.(*net.IPNet)
// 		if isValidIPNew && !ipNet.IP.IsLoopback() {
// 			if ipNet.IP.To4() != nil {
// 				ips = append(ips, ipNet.IP)
// 			}
// 		}
// 	}
// 	return ips
// }

// // AddGetDataToURL xxx
// func AddGetDataToURL(urlString string, data url.Values) string {
// 	if strings.Contains(urlString, "?") {
// 		urlString = urlString + "&"
// 	} else {
// 		urlString = urlString + "?"
// 	}
// 	return fmt.Sprintf("%s%s", urlString, data.Encode())
// }

// // Substr 截取字符串
// func Substr(str string, start int64, end int64) string {
// 	length := int64(len(str))
// 	if start < 0 || start > length || end < 0 {
// 		return ""
// 	}

// 	if end > length {
// 		end = length
// 	}
// 	return string(str[start:end])
// }
