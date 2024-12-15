package snowflake_id

import (
	"github.com/bwmarrin/snowflake"
	"strconv"
	"time"
)

var node *snowflake.Node

func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	// 格式化 1月2号下午3时4分5秒  2006年
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	snowflake.Epoch = st.UnixNano() / 1e6
	node, err = snowflake.NewNode(machineID)
	if err != nil {
		return
	}
	return
}

// GenID 生成 64 位的 雪花 ID
func GenID() string {
	return strconv.FormatInt(node.Generate().Int64(), 10)
}

// NewSnowFlakeNode ， snowflake id generate
// 参数：
// 返回值：
//
//	*snowflake.Node ：desc
func NewSnowFlakeNode() *snowflake.Node {
	// 使用固定的 Epoch，避免冲突
	fixedEpoch := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	snowflake.Epoch = fixedEpoch.UnixNano() / 1000000
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	return node
}

// IIdGenerate id生成器接口
type IIdGenerate interface {
	GenStringId() string
	GenInt64Id() int64
}

type SnowIdGen struct {
	sf *snowflake.Node
}

func NewSnowIdGen() IIdGenerate {
	return &SnowIdGen{
		sf: NewSnowFlakeNode(),
	}
}

func (g *SnowIdGen) GenStringId() string {
	return strconv.FormatInt(abs(g.sf.Generate().Int64()), 10)
}

func (g *SnowIdGen) GenInt64Id() int64 {
	return abs(g.sf.Generate().Int64())
}
func abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}
