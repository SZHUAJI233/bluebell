package snowflake

import (
	"time"
	"web/setting"

	sf "github.com/bwmarrin/snowflake"
)

var node *sf.Node

func Init(cfg *setting.AppConfig) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", cfg.StartTime)
	if err != nil {
		return
	}
	sf.Epoch = st.UnixNano() / 1000000
	node, err = sf.NewNode(cfg.MachineID)

	return
}

func GenID() int64 {
	return node.Generate().Int64()
}
