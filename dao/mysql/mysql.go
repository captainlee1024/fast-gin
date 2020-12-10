package mysql

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"
	"unicode"

	mylog "github.com/captainlee1024/fast-gin/log"

	"github.com/captainlee1024/fast-gin/settings"
	// 初始化数据库驱动
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
)

// 全局变量
var (
	DBMapPool       map[string]*sqlx.DB
	GORMMapPool     map[string]*gorm.DB
	DBDefaultPool   *sqlx.DB
	GORMDefaultPool *gorm.DB
)

// Init 初始化 MySQL
func Init() {

}

// InitDBPool 初始化数据库连接池
func InitDBPool() error {
	// sqlx
	if len(settings.ConfMySQLMap.List) == 0 {
		fmt.Printf("[INFO] %s%s\n", time.Now().Format(settings.TimeFormat), "empty mysql config.")
	}

	DBMapPool = map[string]*sqlx.DB{}
	GORMMapPool = map[string]*gorm.DB{}

	var dsn string
	//
	for confName, DBConf := range settings.ConfMySQLMap.List {

		// sqlx
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
			DBConf.User,
			DBConf.Password,
			DBConf.Host,
			DBConf.Port,
			DBConf.DbName)
		dbsqlx, err := sqlx.Connect("mysql", dsn)
		if err != nil {
			return err
		}

		dbsqlx.SetMaxOpenConns(DBConf.MaxOpenConns)
		dbsqlx.SetMaxIdleConns(DBConf.MaxIdleConns)

		// gorm
		dbgorm, err := gorm.Open("mysql", dsn)
		if err != nil {
			return err
		}
		dbgorm.LogMode(true)
		// dbgorm.Log
		dbgorm.SetLogger(&GormLogger{Trace: mylog.NewTrace()})
		dbgorm.DB().SetMaxOpenConns(DBConf.MaxOpenConns)
		dbgorm.DB().SetMaxIdleConns(DBConf.MaxIdleConns)

		DBMapPool[confName] = dbsqlx
		GORMMapPool[confName] = dbgorm
	}

	// 手动配置连接
	if dbpool, err := GetDBPool("defaule"); err == nil {
		DBDefaultPool = dbpool
	}
	if dbpool, err := GetGormPool("default"); err == nil {
		GORMDefaultPool = dbpool
	}
	return nil
}

// GetDBPool 获取一个默认连接
func GetDBPool(name string) (*sqlx.DB, error) {
	if dbpool, ok := DBMapPool[name]; ok {
		return dbpool, nil
	}
	return nil, errors.New("get sqlxpool error")
}

// GetGormPool 从连接中获取一个默认连接
func GetGormPool(name string) (*gorm.DB, error) {
	if dbgorm, ok := GORMMapPool[name]; ok {
		return dbgorm, nil
	}
	return nil, errors.New("get gormpool error")
}

// Close 关闭 DB 连接
func Close() error {
	for _, dbpool := range DBMapPool {
		dbpool.Close()
	}

	for _, dbpool := range GORMMapPool {
		dbpool.Close()
	}
	return nil
}

// DBPoolLogQuery 获取日志
func DBPoolLogQuery(trace *mylog.TraceContext, sqlDB *sqlx.DB, query string,
	args ...interface{}) (*sql.Rows, error) {
	startExecTime := time.Now()
	rows, err := sqlDB.Query(query, args...)
	endExecTime := time.Now()
	if err != nil {
		mylog.Log.Error("sql", trace, mylog.DLTagMySQLFailed, map[string]interface{}{
			"sql":       query,
			"bind":      args,
			"proc_time": fmt.Sprintf("%f", endExecTime.Sub(startExecTime).Seconds()),
		})
	} else {
		mylog.Log.Info("sql", trace, mylog.DLTagMySQLSuccess, map[string]interface{}{
			"sql":       query,
			"bind":      args,
			"proc_time": fmt.Sprintf("%f", endExecTime.Sub(startExecTime).Seconds()),
		})
	}
	return rows, err
}

// GormLogger MySQL 日志打印类
// Logger default logger
type GormLogger struct {
	gorm.Logger
	Trace *mylog.TraceContext
}

// Print 日志输出
func (logger *GormLogger) Print(values ...interface{}) {
	message := logger.LogFormatter(values...)
	if message["level"] == "sql" {
		mylog.Log.Info("sql", logger.Trace, mylog.DLTagMySQLSuccess, message)
	} else {
		mylog.Log.Info("sql", logger.Trace, mylog.DLTagMySQLFailed, message)
	}
}

// TODO:
// CtxPrint LogCtx(true) 时调用的方法
// func (logger *GormLogger) CtxPrint(s *gorm.DB, values ...interface{}) {
// 	ctx, ok := s.GetCtx()
// 	trace := settings.NewTrace()
// 	if ok {
// 		trace = ctx.(*TraceContext)
// 	}
// 	message := logger.LogFormatter(values...)
// 	if message["level"] == "sql" {
// 		mylog.Log.Info("sql", trace, mylog.DLTagMySQLSuccess, message)
// 	} else {
// 		mylog.Log.Info("sql", trace, mylog.DLTagMySQLFailed, message)
// 	}
// }

// LogFormatter 格式化日志格式
func (logger *GormLogger) LogFormatter(values ...interface{}) (messages map[string]interface{}) {
	if len(values) > 1 {
		var (
			sql             string
			formattedValues []string
			level           = values[0]
			currentTime     = logger.NowFunc().Format("2006-01-02 15:04:05")
			source          = fmt.Sprintf("%v", values[1])
		)

		messages = map[string]interface{}{
			"level":   level,
			"source":  source,
			"current": currentTime}

		if level == "sql" {
			// duration
			messages["proc_time"] = fmt.Sprintf("%fs", values[2].(time.Duration).Seconds())

			// sql
			for _, value := range values[4].([]interface{}) {
				indirectValue := reflect.Indirect(reflect.ValueOf(value))
				if indirectValue.IsValid() {
					value = indirectValue.Interface()
					if t, ok := value.(time.Time); ok {
						formattedValues = append(formattedValues, fmt.Sprintf("%v", t.Format("2006-01-02 15:04:05")))
					} else if b, ok := value.([]byte); ok {
						if str := string(b); logger.isPrintable(str) {
							formattedValues = append(formattedValues, fmt.Sprintf("%v", value))

						} else {
							formattedValues = append(formattedValues, "'<binary>'")
						}
					} else if r, ok := value.(driver.Valuer); ok {
						if value, err := r.Value(); err != nil && value != nil {
							formattedValues = append(formattedValues, fmt.Sprintf("%v", value))
						} else {
							formattedValues = append(formattedValues, "NULL")
						}
					} else {
						formattedValues = append(formattedValues, fmt.Sprintf("%v", value))
					}

				} else {
					formattedValues = append(formattedValues, "NULL")
				}
			}

			// differentiate between $n placeholders or  else treat like?
			if regexp.MustCompile(`\$\d+`).MatchString(values[3].(string)) {
				sql = values[3].(string)
				for index, value := range formattedValues {
					placeholder := fmt.Sprintf(`\$%d([^\d]|$)`, index+1)
					sql = regexp.MustCompile(placeholder).ReplaceAllString(sql, value+"$1")
				}
			} else {
				formattedValuesLength := len(formattedValues)
				for index, value := range regexp.MustCompile(`\?`).Split(values[3].(string), -1) {
					sql += value
					if index < formattedValuesLength {
						sql += formattedValues[index]
					}
				}
			}

			messages["sql"] = sql
			if len(values) > 5 {
				messages["affected_row"] = strconv.FormatInt(values[5].(int64), 10)
			}
		} else {
			messages["ext"] = values
		}
	}
	return
}

// NowFunc 获取当前时间
func (logger *GormLogger) NowFunc() time.Time {
	return time.Now()
}

func (logger *GormLogger) isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}
