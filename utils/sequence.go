package utils

import "github.com/bwmarrin/snowflake"

var (
	global *snowflake.Node
)

func InitSnowflake() error {
	g, err := snowflake.NewNode(1)
	if err != nil {
		return err
	}

	global = g
	return nil
}

func GetId() string {
	return global.Generate().String()
}
