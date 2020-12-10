package settings

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
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

	// 初始化配置文件 ? => 作用
	// if err := initViperConf(); err != nil {
	// 	return err
	// }

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
