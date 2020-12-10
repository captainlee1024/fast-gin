package redis

import (
	"errors"
	"fmt"
	"time"

	mylog "github.com/captainlee1024/fast-gin/log"
	"github.com/captainlee1024/fast-gin/settings"
	"github.com/garyburd/redigo/redis"
)

/*
// go-redis
var (
	rdb *redis.Client
	// rdbMap map[string]*redis.Client
)

// Init 初始化 redis 默认连接
func Init(name string) (err error) {
	if settings.ConfRedisMap != nil && settings.ConfRedisMap.List != nil {
		for confName, cfg := range settings.ConfRedisMap.List {
			if name == confName {
				rdb = redis.NewClient(&redis.Options{
					Addr: fmt.Sprintf("%s:%d",
						cfg.Host,
						cfg.Port,
					),
					Password: cfg.Password,
					DB:       cfg.DB,
					PoolSize: cfg.PoolSize,
				})
			}
		}
	}
	_, err = rdb.Ping().Result()
	return
}

// Close 关闭 redis 连接
func Close() {
	_ = rdb.Close()
}
*/

// redigo

// ConnFactory 获取连接
func ConnFactory(name string) (redis.Conn, error) {
	if settings.ConfRedisMap != nil && settings.ConfRedisMap.List != nil {
		for confName, cfg := range settings.ConfRedisMap.List {
			if name == confName {
				// randHost := cfg.ProxyList[rand.Intn(len(cfg.ProxyList))]
				randHost := cfg.Host + fmt.Sprintf(":%d", cfg.Port)

				// if cfg.ConnTimeout == 0 {
				// 	cfg.ConnTimeout = 50
				// }
				// if cfg.ReadTimeout == 0 {
				// 	cfg.ReadTimeout = 100
				// }
				// if cfg.WriteTimeout == 0 {
				// 	cfg.WriteTimeout = 100
				// }
				c, err := redis.Dial(
					"tcp",
					randHost,
					// redis.DialConnectTimeout(time.Duration(cfg.ConnTimeout)*time.Millisecond),
					// redis.DialReadTimeout(time.Duration(cfg.ReadTimeout)*time.Millisecond),
					// redis.DialWriteTimeout(time.Duration(cfg.WriteTimeout)*time.Millisecond),
				)
				if err != nil {
					return nil, err
				}
				if cfg.Password != "" {
					if _, err := c.Do("AUTH", cfg.Password); err != nil {
						c.Close()
						return nil, err
					}
				}
				if cfg.DB != 0 {
					if _, err := c.Do("SELECT", cfg.DB); err != nil {
						c.Close()
						return nil, err
					}
				}
				return c, nil
			}
		}
	}
	return nil, errors.New("create redis conn fail")
}

// LogDo 带有日志的 Do 方法
func LogDo(trace *mylog.TraceContext, c redis.Conn, commandName string, args ...interface{}) (interface{}, error) {
	startExecTime := time.Now()
	reply, err := c.Do(commandName, args...)
	endExecTime := time.Now()
	if err != nil {
		mylog.Log.Error("redis", trace, "_com_redis_failure", map[string]interface{}{
			"method":    commandName,
			"err":       err,
			"bind":      args,
			"proc_time": fmt.Sprintf("%fs", endExecTime.Sub(startExecTime).Seconds()),
		})
	} else {
		replyStr, _ := redis.String(reply, nil)
		mylog.Log.Info("redis", trace, "_com_redis_success", map[string]interface{}{
			"method":    commandName,
			"bind":      args,
			"reply":     replyStr,
			"proc_time": fmt.Sprintf("%fs", endExecTime.Sub(startExecTime).Seconds()),
		})
	}
	return reply, err
}

// ConfDo 通过配置 执行redis
func ConfDo(trace *mylog.TraceContext, name string, commandName string, args ...interface{}) (interface{}, error) {
	c, err := ConnFactory(name)
	if err != nil {
		mylog.Log.Error("redis", trace, "_com_redis_failure", map[string]interface{}{
			"method": commandName,
			"err":    errors.New("RedisConnFactory_error:" + name),
			"bind":   args,
		})
		return nil, err
	}
	defer c.Close()

	startExecTime := time.Now()
	reply, err := c.Do(commandName, args...)
	endExecTime := time.Now()
	if err != nil {
		mylog.Log.Error("redis", trace, "_com_redis_failure", map[string]interface{}{
			"method":    commandName,
			"err":       err,
			"bind":      args,
			"proc_time": fmt.Sprintf("%fs", endExecTime.Sub(startExecTime).Seconds()),
		})
	} else {
		replyStr, _ := redis.String(reply, nil)
		mylog.Log.Info("redis", trace, "_com_redis_success", map[string]interface{}{
			"method":    commandName,
			"bind":      args,
			"reply":     replyStr,
			"proc_time": fmt.Sprintf("%fs", endExecTime.Sub(startExecTime).Seconds()),
		})
	}
	return reply, err
}
