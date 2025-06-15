package snowflake

import "github.com/nortoo/utils-go/generator/snowflake"

var cli *snowflake.Worker

func GetSnowWorker() *snowflake.Worker {
	return cli
}

func Init(id int64) error {
	var err error
	cli, err = snowflake.NewWorker(id)
	return err
}
