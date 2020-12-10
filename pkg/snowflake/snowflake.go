// Package snowflake 基于雪花算法的 ID 生成器，可用于分布式系统
package snowflake

import (
	"time"

	sf "github.com/bwmarrin/snowflake"
)

// node 初始化一个全局的 node 实例
var node *sf.Node

// Init 初始化一个雪花算法实例
func Init(startTime string, machineID int64) (err error) {
	// 设置时间因子 起始时间 因为是分布式的，所以要指定机器的ID
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	// 初始化开始的时间
	sf.Epoch = st.UnixNano() / 1000000
	// 设置机器ID，并拿到 node 节点
	node, err = sf.NewNode(machineID)
	return
}

// GenID 生成ID
func GenID() int64 {
	return node.Generate().Int64()
}

/*
func main() {
	// 传入参数指定起始时间和机器ID
	if err := Init("2020-07-01", 1); err != nil {
		fmt.Printf("init failed, err:%v\n", err)
		return
	}
	id := GenID()
	fmt.Println(id)
}
*/

// 索尼开源的雪花算法
/*
import (
	"fmt"
	"github.com/sony/sonyflake"
	"time"
)

var (
	sonyFlake     *sonyflake.Sonyflake
	sonyMachineID uint16
)

func getMachineID() (uint16, error) {
	return sonyMachineID, nil
}

// 需传⼊当前的机器ID
func Init(startTime string, machineId uint16) (err error) {
	sonyMachineID = machineId
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return err
	}
	settings := sonyflake.Settings{
		StartTime: st,
		MachineID: getMachineID,
	}
	sonyFlake = sonyflake.NewSonyflake(settings)
	return
}

// GenID ⽣成id
func GenID() (id uint64, err error) {
	if sonyFlake == nil {
		err = fmt.Errorf("snoy flake not inited")
		return
	}
	id, err = sonyFlake.NextID()
	return
}
func main() {
	if err := Init("2020-07-01", 1); err != nil {
		fmt.Printf("Init failed, err:%v\n", err)
		return
	}
	id, _ := GenID()
	fmt.Println(id)
}
*/
