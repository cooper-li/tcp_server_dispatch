package control

import (
	"sync"
	"assets_server/depend/conf"
)

var (
	wg *sync.WaitGroup
)

func init()  {
	wg = &sync.WaitGroup{}
}

func MainControl()  {

	// 初始化DB配置
	go conf.InitDatabase()

	// 启动tcp服务
	InitServer()

	wg.Wait()
}
