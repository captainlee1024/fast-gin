// Package po PO（Persistent Object）：持久化对象
package po

import "time"

// FastGin fastgin 持久化实体
type FastGin struct {
	FastGinID  int64     `db:"fast_gin_id"`
	DemoName   string    `db:"demo_name"`
	Info       string    `db:"info"`
	CreateTime time.Time `db:"create_time"`
	UpdateTime time.Time `db:"update_time"`
}
