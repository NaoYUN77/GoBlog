package snowflake


import (
	"github.com/bwmarrin/snowflake"
	"time"
	"Blog/settings"


)

var node *snowflake.Node
func Init() (err error) {
	var st time.Time
	//开始时间
	startTime := settings.Cnf.App.StartTime

	//机器id
	MachineID := settings.Cnf.App.MechineId

	st, err = time.Parse("2006-01-01 15:04:05", startTime)
	if err != nil {
		return
	}
	snowflake.Epoch = st.UnixNano() / 1000000
	node, err = snowflake.NewNode(MachineID)
	return
}


//用雪化算法生成文章唯一id
func GenId() int64 {
	return node.Generate().Int64()
}