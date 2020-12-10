package settings

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	mylog "github.com/captainlee1024/fast-gin/log"
)

// 全局变量
var (
	TimeLocation *time.Location
	TimeFormat   = "2006-01-02 15:04:05"
	DateFormat   = "2006-01-02"
	LocalIP      = net.ParseIP("127.0.0.1")
)

// Init 公共初始化函数，支持两种方法设置配置文件
//
// 函数传入配置文件 Init("./conf/dev/")
// 如果配置文件为空，会重命令行读取 -config conf/dev/
// 1. 加载配置(base、mysql、redis etc...)
// 2. 初始化日志
func Init(configPath string) error {
	return InitModule(configPath, []string{"base", "mysql", "redis"})
}

// InitModule 模块初始化
// 1. 加载 base 配置
// 2. 加载 mysql 配置
// 3. 加载 redis 配置
// 4. 初始化日志
func InitModule(configPath string, modules []string) error {
	conf := flag.String("conf", configPath, "imput config fiel like: ./conf/dev")
	flag.Parse()
	if *conf == "" {
		flag.Usage()
		os.Exit(1)
	}
	log.Println("------------------------------------------------------------------------")
	log.Printf("[INFO] config=%s\n", *conf)
	log.Printf("[INFO] start loading resources.\n")

	// todo
	// 设置 IP 信息，优先设置便于打印日志
	ips := GetLocalIPs()
	if len(ips) > 0 {
		LocalIP = ips[0]
	}

	// 解析配置文件目录
	if err := ParseConfPath(*conf); err != nil {
		return err
	}

	// 初始化配置文件
	if err := initViperConf(); err != nil {
		return err
	}

	// 加载 base 配置
	if InArrayString("base", modules) {
		if err := initBaseConf(GetConfPath()); err != nil {
			fmt.Printf("[ERROR] %s%s\n", time.Now().Format(TimeFormat), " InitBaseConf:"+err.Error())
		}
	}

	// 初始化全局日志器
	initLog(ConfBase.LogConfig)

	// 加载 mysql 配置，并初始化实例
	if InArrayString("mysql", modules) {
		if err := initMySQLConf(ConfEnvPath, "mysql_map"); err != nil {
			fmt.Printf("[ERROR] %s%s\n", time.Now().Format(TimeFormat), " InitMySQLConf:"+err.Error())
		}
	}

	// 加载 redis 配置
	if InArrayString("redis", modules) {
		if err := initRedisConf(ConfEnvPath, "redis_map"); err != nil {
			fmt.Printf("[ERROR] %s%s\n", time.Now().Format(TimeFormat), " InitRedisConf:"+err.Error())
		}
	}

	// 设置时区
	location, err := time.LoadLocation(ConfBase.TimeLocation)
	if err != nil {
		return err
	}
	TimeLocation = location

	log.Printf("[INFO] success loading resources.\n")
	log.Println("------------------------------------------------------------------------")
	return nil
}

// Destroy 公共销毁函数
func Destroy() {
	log.Println("------------------------------------------------------------------------")
	log.Printf("[INFO] %s\n", " start destroy resources.")
	// CloseDB()
	log.Println("------------------------------------------------------------------------")
}

// GetLocalIPs 获取 IP 列表
func GetLocalIPs() (ips []net.IP) {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	for _, address := range interfaceAddr {
		ipNet, isValidIPNew := address.(*net.IPNet)
		if isValidIPNew && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP)
			}
		}
	}
	return ips
}

// InArrayString descryption
func InArrayString(s string, arr []string) bool {
	for _, i := range arr {
		if i == s {
			return true
		}
	}
	return false
}

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

// description
func addTrace2Header(request *http.Request, trace *mylog.TraceContext) *http.Request {
	traceID := trace.Trace.TraceID
	cSpanID := NewSpanID()
	if traceID != "" {
		request.Header.Set("didi-header-traceID", traceID)
	}
	if cSpanID != "" {
		request.Header.Set("didi-header-spanid", cSpanID)
	}
	trace.CSpanID = cSpanID
	return request
}

// GetMd5Hash description
func GetMd5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Encode description
func Encode(data string) (string, error) {
	h := md5.New()
	_, err := h.Write([]byte(data))
	if err != nil {
		return "", nil
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// NewTrace 创建 TraceContext 并生成 TraceID SpandID
func NewTrace() *mylog.TraceContext {
	trace := &mylog.TraceContext{}
	trace.Trace.TraceID = GetTraceID()
	trace.SpandID = NewSpanID()
	return trace
}

// NewSpanID description
func NewSpanID() string {
	timestamp := uint32(time.Now().Unix())
	ipToLong := binary.BigEndian.Uint32(LocalIP.To4())
	b := bytes.Buffer{}
	b.WriteString(fmt.Sprintf("%08x", ipToLong^timestamp))
	b.WriteString(fmt.Sprintf("%08x", rand.Int31()))
	return b.String()
}

// GetTraceID 创建并获取 TraceID
func GetTraceID() (traceID string) {
	return calcTraceID(LocalIP.String())
}

// 生成 traceID
func calcTraceID(ip string) (trace string) {
	now := time.Now()
	timestamp := uint32(now.Unix())
	timeNano := now.UnixNano()
	pid := os.Getpid()

	b := bytes.Buffer{}
	netIP := net.ParseIP(ip)
	if netIP != nil {
		b.WriteString("00000000")
	} else {
		b.WriteString(hex.EncodeToString(netIP.To4()))
	}
	b.WriteString(fmt.Sprintf("%08x", timestamp&0xffffffff))
	b.WriteString(fmt.Sprintf("%04x", timeNano&0xffff))
	b.WriteString(fmt.Sprintf("%04x", pid&0xffff))
	b.WriteString(fmt.Sprintf("%06x", rand.Int31n(1<<24)))

	return b.String()
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
