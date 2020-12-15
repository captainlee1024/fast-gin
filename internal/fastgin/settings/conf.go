package settings

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	mylog "github.com/captainlee1024/fast-gin/internal/fastgin/log"
)

// 配置信息全局变量
var (
	ConfBase = new(BaseConfig)
	//ConfMySQL    = new(MySQLConfig)
	ConfMySQLMap = new(MySQLMapConfig)
	//ConfRedis    = new(RedisConfig)
	ConfRedisMap = new(RedisMapConfig)
	ViperConfMap map[string]*viper.Viper
)

// BaseConfig 应用程序配置信息
type BaseConfig struct {
	Name         string `mapstructure:"name"`
	Mode         string `mapstructure:"mode"`
	Version      string `mapstructure:"version"`
	StartTime    string `mapstructure:"start_time"`
	TimeLocation string `mapstructure:"time_location"`
	MachineID    int64  `mapstructure:"machine_id"`
	Port         int    `mapstructure:"port"`

	*LogConfig `mapstructure:"log"`
}

// LogConfig Zap 配置信息
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backup"`
}

// MySQLMapConfig 数据库列表
type MySQLMapConfig struct {
	List map[string]*MySQLConfig `mapstructure:"list"`
}

// MySQLConfig MySQL 配置信息
type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DbName       string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idel_conns"`
}

// RedisMapConfig Redis 列表
type RedisMapConfig struct {
	List map[string]*RedisConfig `mapstructure:"list"`
}

// RedisConfig Redis 配置信息
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// GetBaseConf 获取基本配置信息
func GetBaseConf() *BaseConfig {
	return ConfBase
}

// 读取基本配置
func initBaseConf(path string) (err error) {
	err = ParseConfig(path, "base", ConfBase)
	if err != nil {
		return err
	}

	if ConfBase.Mode == "" {
		ConfBase.Mode = "debug"
	}

	if ConfBase.TimeLocation == "" {
		ConfBase.TimeLocation = "Asiz/Shanghai"
	}

	return

}

// initRedisConf 初始化 Redis 配置信息
func initRedisConf(path string, fileName string) (err error) {
	err = ParseConfig(path, fileName, ConfRedisMap)
	if err != nil {
		return err
	}
	//ConfRedis = ConfRedisMap.List["default"]
	return nil
}

// initMySQLConf 初始化数据库信息
func initMySQLConf(path string, fileName string) (err error) {
	err = ParseConfig(path, fileName, ConfMySQLMap)
	if err != nil {
		return err
	}
	//ConfMySQL = ConfMySQLMap.List["default"]
	return
}

// initLog 初始化 Zap logger
// 配置自己的 logger　并替换 zap 中定义的全局变量 logger
// func initLog(cfg *LogConfig, mode string) (err error) {
func initLog(cfg *LogConfig) (err error) {
	writerSyncer := mylog.GetLogWriter(
		cfg.Filename,
		cfg.MaxSize,
		cfg.MaxBackups,
		cfg.MaxAge,
	)

	encoder := mylog.GetEncoder()

	// 把 yaml 配置文件中的 string 类型的 level 配置，解析成 zap 中的 level 类型
	var l = new(zapcore.Level)
	//err = l.UnmarshalText([]byte(viper.GetString("log.level")))
	err = l.UnmarshalText([]byte(viper.GetString(cfg.Level)))
	if err != nil {
		return
	}

	var core zapcore.Core
	if cfg.Level == "debug" {
		// 开发模式，输出日志到终端
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee( // 指定两个输出位置
			// 第一个输出　和下面的配置一样，以 JSON 方式写入到日志文件里面
			zapcore.NewCore(encoder, writerSyncer, l),

			// 第二个输出
			// consoleEncoder 指定console编码器
			// zapcore.Lock(os.Stdout) 指定输出位置是标准输出，给它转换成符合条件的 WriteSyncer
			// zapcore.DebugLevel 指定输出日志级别
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)

	} else {
		//
		core = zapcore.NewCore(encoder, writerSyncer, zapcore.DebugLevel)

	}

	// 生成配置的 logger
	mylog.Log = new(mylog.Logger)
	mylog.Log.L = zap.New(core, zap.AddCaller())
	return
}

// initViperConf 初始化配置文件
func initViperConf() error {
	f, err := os.Open(ConfEnvPath + "/")
	if err != nil {
		return err
	}

	fileList, err := f.Readdir(1024)
	if err != nil {
		return err
	}

	for _, f0 := range fileList {
		if !f0.IsDir() {
			bts, err := ioutil.ReadFile(ConfEnvPath + "/" + f0.Name())
			if err != nil {
				return err
			}
			v := viper.New()
			v.SetConfigType("yaml")
			v.ReadConfig(bytes.NewBuffer(bts))
			pathArr := strings.Split(f0.Name(), ".")
			if ViperConfMap == nil {
				ViperConfMap = make(map[string]*viper.Viper)
			}
			ViperConfMap[pathArr[0]] = v
		}
	}
	// fmt.Printf("%v\n%v\n%v\n", ViperConfMap["base"],
	// 	ViperConfMap["mysql"], ViperConfMap["redis"])
	return nil
}

// GetStringConf 获取 string 类型配置信息
func GetStringConf(key string) string {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return ""
	}
	v, ok := ViperConfMap[keys[0]]
	if !ok {
		return ""
	}

	confString := v.GetString(strings.Join(keys[1:len(keys)], "."))
	return confString
}

// GetStringMapConf 获取 string 为 key 的 map 类型信息
func GetStringMapConf(key string) map[string]interface{} {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return nil
	}

	v := ViperConfMap[keys[0]]
	conf := v.GetStringMap(strings.Join(keys[1:len(keys)], "."))
	return conf
}

// GetBoolConf 获取 bool 类型配置信息
func GetBoolConf(key string) bool {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return false
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetBool(strings.Join(keys[1:len(keys)], "."))
	return conf
}

// GetFloat64Conf 获取 float64 类型配置信息
func GetFloat64Conf(key string) float64 {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return 0
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetFloat64(strings.Join(keys[1:len(keys)], "."))
	return conf
}

// GetIntConf 获取 int 类型配置信息
func GetIntConf(key string) int {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return 0
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetInt(strings.Join(keys[1:len(keys)], "."))
	return conf
}

// GetStringMapStringConf 获取 map[string]string 类型配置信息
func GetStringMapStringConf(key string) map[string]string {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return nil
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetStringMapString(strings.Join(keys[1:len(keys)], "."))
	return conf
}

// GetStringSliceConf 获取  []string 类型配置信息
func GetStringSliceConf(key string) []string {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return nil
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetStringSlice(strings.Join(keys[1:len(keys)], "."))
	return conf
}

// GetTimeConf 获取 time.Time 类型配置信息
func GetTimeConf(key string) time.Time {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return time.Now()
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetTime(strings.Join(keys[1:len(keys)], "."))
	return conf
}

// GetDurationConf 获取时间阶段长度
func GetDurationConf(key string) time.Duration {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return 0
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetDuration(strings.Join(keys[1:len(keys)], "."))
	return conf
}

// IsSetConf 是否设置了key
func IsSetConf(key string) bool {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return false
	}
	v := ViperConfMap[keys[0]]
	conf := v.IsSet(strings.Join(keys[1:len(keys)], "."))
	return conf
}
