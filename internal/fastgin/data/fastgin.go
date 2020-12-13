package data

import (
	"github.com/captainlee1024/fast-gin/internal/fastgin/data/mysql"
	"github.com/captainlee1024/fast-gin/internal/fastgin/do"
	"github.com/captainlee1024/fast-gin/internal/fastgin/middleware"
	"github.com/captainlee1024/fast-gin/internal/fastgin/po"
	"github.com/captainlee1024/fast-gin/internal/fastgin/service"
	"github.com/captainlee1024/fast-gin/pkg/snowflake"
	"github.com/gin-gonic/gin"
	"time"
)

// 在编译的时候可以知道这个对象实现了这个 interface{}
var _ service.FastGinDoRepo = (service.FastGinDoRepo)(nil)

// NewFastGinRepo 创建一个 fastGinRepo ，它是 service.FastGinDoRepo 的实现
func NewFastGinRepo() service.FastGinDoRepo {
	return new(fastGinRepo)
}

type fastGinRepo struct{}

// SaveFastGin 保存 FastGinPo 至数据库
func (fg *fastGinRepo) SaveFastGin(fgDo *do.FastGinDo, c *gin.Context) (err error) {
	// do -> po
	currentTime := time.Now()
	fgPo := &po.FastGin{
		FastGinID:  snowflake.GenID(),
		DemoName:   fgDo.DemoName,
		Info:       fgDo.Info,
		CreateTime: currentTime,
		UpdateTime: currentTime,
	}

	db, err := mysql.GetDBPool("default")
	if err != nil {
		return err
	}

	sqlStr := `INSERT INTO fast_gin(fast_gin_id, demo_name, info, create_time, update_time)
			VALUES(:fast_gin_id, :demo_name, :info, :create_time, :update_time)`

	trace := middleware.GetGinTraceContext(c)
	_, err = mysql.SqlxLogNamedExec(trace, db, sqlStr, fgPo)
	//if err != nil {
	//	return err
	//}
	return
}

func (fg *fastGinRepo) GetFastGin(fgDo *do.FastGinDo, c *gin.Context) (fastGin *do.FastGinDo, err error) {
	db, err := mysql.GetDBPool("default")
	if err != nil {
		return nil, err
	}
	// do -> po

	trace := middleware.GetGinTraceContext(c)
	sqlStr := `SELECT fast_gin_id, demo_name, info
			FROM fast_gin
			WHERE demo_name = ?`
	fastGinPo := &po.FastGin{}
	err = mysql.SqlxLogGet(trace, db, fastGinPo, sqlStr, fgDo.DemoName)
	if err != nil {
		return nil, err
	}

	// 使用 NameQuery() 方法
	//fgPo := &po.FastGin{
	//	DemoName: fgDo.DemoName,
	//	Info:     fgDo.Info,
	//}
	//sqlStr := `SELECT demo_name, info
	//		FROM fast_gin
	//		WHERE demo_name = :demo_name OR info = :info`
	//rows, err := mysql.SqlxLogNamedQuery(trace, db, sqlStr, fgPo)
	//if rows != nil {
	//	defer rows.Close()
	//}
	//
	//for rows.Next() {
	//	err = rows.StructScan(fastGinPo)
	//	if err != nil {
	//		return nil, err
	//	}
	//}

	fastGin = &do.FastGinDo{
		FastGinID: fastGinPo.FastGinID,
		DemoName:  fastGinPo.DemoName,
		Info:      fastGinPo.Info,
	}

	return
}

func (fg *fastGinRepo) GetFastGinList(page, size int, c *gin.Context) (listFastGin []*do.FastGinDo, err error) {
	db, err := mysql.GetDBPool("default")
	if err != nil {
		return nil, err
	}

	// do -> po 这里不需要，省略

	listFastGinPo := make([]*po.FastGin, 0, 2)

	trace := middleware.GetGinTraceContext(c)
	sqlStr := `SELECT fast_gin_id, demo_name, info
		FROM fast_gin
		LIMIT ?,?`
	//sqlStr := `SELECT fast_gin_id, demo_name, info FROM fast_gin LIMIT 0, 10`
	err = mysql.SqlxLogSelect(trace, db, &listFastGinPo, sqlStr, (page-1)*size, size)
	if err != nil {
		return nil, err
	}

	// po -> do
	// 首先初始化返回值定义的变量，那里只是声明，并没有申请内存
	listFastGin = make([]*do.FastGinDo, 0, len(listFastGinPo))
	for _, fastGinPo := range listFastGinPo {
		fastGinDo := &do.FastGinDo{
			FastGinID: fastGinPo.FastGinID,
			DemoName:  fastGinPo.DemoName,
			Info:      fastGinPo.Info,
		}
		listFastGin = append(listFastGin, fastGinDo)
	}

	return
}

func (fg *fastGinRepo) UpdateFastGin(fgDo *do.FastGinDo, c *gin.Context) (err error) {
	db, err := mysql.GetDBPool("default")
	if err != nil {
		return err
	}

	// do -> po
	currentTime := time.Now()
	fgPo := &po.FastGin{
		FastGinID:  fgDo.FastGinID,
		DemoName:   fgDo.DemoName,
		Info:       fgDo.Info,
		UpdateTime: currentTime,
	}

	trace := middleware.GetGinTraceContext(c)
	sqlStr := `UPDATE fast_gin
			SET demo_name=?, info=?, update_time=?
			WHERE fast_gin_id=?`
	_, err = mysql.SqlxLogExec(trace, db, sqlStr, fgPo.DemoName, fgPo.Info, fgPo.UpdateTime, fgPo.FastGinID)
	if err != nil {
		return err
	}
	return
}

func (fg *fastGinRepo) DeleteFastGin(fgDo *do.FastGinDo, c *gin.Context) (err error) {
	db, err := mysql.GetDBPool("default")
	if err != nil {
		return err
	}

	// do -> po
	fgPo := &po.FastGin{
		FastGinID: fgDo.FastGinID,
	}
	trace := middleware.GetGinTraceContext(c)
	sqlStr := `DELETE FROM fast_gin
			WHERE fast_gin_id = ?`
	_, err = mysql.SqlxLogExec(trace, db, sqlStr, fgPo.FastGinID)
	if err != nil {
		return err
	}
	// 返回受影响行数
	//_, err = ret.RowsAffected()
	return
}
