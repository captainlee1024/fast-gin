package settings

import (
	"fmt"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 配置文件全局变量
var (
	ConfEnvPath string // 配置文件路径
	ConfEnv     string // 模式
)

// 默认读取配置文件路径
const (
	DefaultConfEnvPath = "./configs/dev/"
)

// ParseConfPath 解析配置文件目录
//
// 配置文件必须放到一个文件夹中
// 如：config=conf/dev/base.json 	ConfEnvPath=conf/dev	ConfEnv=dev
// 如：config=conf/base.json		ConfEnvPath=conf		ConfEnv=conf
func ParseConfPath(config string) error {
	path := strings.Split(config, "/")
	prefix := strings.Join(path[:len(path)-1], "/")
	ConfEnvPath = prefix
	ConfEnv = path[len(path)-2]
	return nil
}

// GetConfEnv 获取配置环境名
func GetConfEnv() string {
	return ConfEnv
}

// GetConfPath 获取配置目录路径
func GetConfPath() string {
	return ConfEnvPath
}

// GetConfFilePath 获取文件路径
func GetConfFilePath(fileName string) string {
	return ConfEnvPath + "/" + fileName
}

// GetConfFile 获取带后缀名的路径
func GetConfFile(fileName string) string {
	return ConfEnvPath + "/" + fileName + ".yaml"
}

// ParseLocalConfig 解析本地文件
func ParseLocalConfig(fileName string, st interface{}) error {
	path := GetConfPath()
	return ParseConfig(path, fileName, st)
}

// ParseConfig 解析本地配置文件
func ParseConfig(path string, fileName string, conf interface{}) (err error) {
	v := viper.New()
	v.SetConfigName(fileName)
	v.AddConfigPath(path)

	if err = v.ReadInConfig(); err != nil {
		return fmt.Errorf("Read config %v failed, err: %v", fileName, err)
	}

	if err = v.Unmarshal(conf); err != nil {
		return fmt.Errorf("unmarshal failed, err: %v", err)
	}

	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println(time.Now().Format(TimeFormat)+" ------------------------------------------------------------------------")
		fmt.Printf("%s [INFO] %v config file changed...\n", time.Now().Format(TimeFormat), fileName)
		if err = v.Unmarshal(conf); err != nil {
			err = fmt.Errorf("unmarshal changed conf failed, err: %v", err)
		}
		fmt.Printf("%s [INFO] ==> %v\n", time.Now().Format(TimeFormat), ConfBase)
		// fmt.Printf("===>%#v\n", ConfBase)
		fmt.Printf("%s [INFO] ===> %#v\n", time.Now().Format(TimeFormat),ConfRedisMap.List["default"])
		//fmt.Printf("===>%#v\n", ConfRedis)
		fmt.Printf("%s [INFO] ===> %#v\n", time.Now().Format(TimeFormat),ConfMySQLMap.List["default"])
		//fmt.Printf("===>%#v\n", ConfMySQL)
		fmt.Println(time.Now().Format(TimeFormat) + " [INFO] 配置已同步...")
		fmt.Println(time.Now().Format(TimeFormat)+" ------------------------------------------------------------------------")
	})

	return nil
}
