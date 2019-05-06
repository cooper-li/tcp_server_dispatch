package control

import (
	"fmt"
	"sync"
		"assets_server/depend/server"
	"os_go_comm/comm_cfg"
	"math/rand"
)

type (
	taskProcess struct {
		maxGo  int
		goPool sync.Map
	}
)

var (
	NewTaskProcess *taskProcess
)

func init() {
	NewTaskProcess = &taskProcess{
		maxGo:  comm_cfg.Int("system","process_num "),
		goPool: sync.Map{},
	}
	NewTaskProcess.createGoPool()

	fmt.Println("init task pool success")
}

// 创建process
func (process *taskProcess) createGoPool() {
	for i := 0; i < process.maxGo; i++ {
		go func(processId int) {

			ch := make(chan *server.ServerContext, 128)
			process.goPool.Store(processId, ch)

			for {
				select {
				case data, ok := <-ch:
					if !ok {
						fmt.Println("process is not ok ", processId)
					}
				 	//fmt.Println("执行进程: ", processId)
					process.dealData(data)
				}
			}
		}(i)
	}
}

// 处理具体任务
// 维护计数器
func (process *taskProcess) dealData(data *server.ServerContext) {

	// 执行方法
	err := server.ServerHandler.Call(data)
	fmt.Println("执行完毕，错误信息: ", err)

	// 维护计数
	if data.Body.Md.Uid > 0 {
		uidCount := NewDispatcher.activeUserCount(data.Body.Md.Uid, -1)
		if uidCount <= 0 {
			NewDispatcher.activeUserProcessMap.Delete(data.Body.Md.Uid)
		}
	}

	if data.Body.Md.BakUid > 0 {
		bakUidCount := NewDispatcher.activeUserCount(data.Body.Md.BakUid, -1)
		if bakUidCount <= 0 {
			NewDispatcher.activeUserProcessMap.Delete(data.Body.Md.BakUid)
		}
	}

}

// 获取随机协程ID
func (process *taskProcess) GetRandProcess() int {
	return rand.Intn(process.maxGo)
}

// 获取空闲进程
func (process *taskProcess) GetIdleProcess() int {

	lenProcess := 0
	idleProcess := 0
	for i:=0;i<process.maxGo;i++ {

		ch, _ := NewTaskProcess.goPool.Load(i)
		l := len(ch.(chan *server.ServerContext))
		if l == 0 {
			return i
		}

		if l < lenProcess {
			idleProcess = i
			lenProcess = l
		}
	}

	return idleProcess

}
