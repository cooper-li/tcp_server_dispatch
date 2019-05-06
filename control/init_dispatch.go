package control

import (
	"sync"
	"fmt"
	"net"
		"encoding/json"
			"assets_server/depend/comm"
	"assets_server/depend/server"
)

type (
	Dispatcher struct {
		// 连接信息
		lock                 sync.RWMutex
		requestChan          chan *requestInfo
		activeUserProcessMap sync.Map    // uid —> processId
		activeUserCountMap   map[int]int // uid -> count
	}

	requestInfo struct {
		conn *net.TCPConn
		msg  []byte
	}
)

var (
	NewDispatcher *Dispatcher
)

func init() {
	NewDispatcher = &Dispatcher{
		lock:                 sync.RWMutex{},
		requestChan:          make(chan *requestInfo, 100),
		activeUserProcessMap: sync.Map{},
		activeUserCountMap:   map[int]int{},
	}

	go NewDispatcher.dispatcher()

	fmt.Println("init dispatch success")
}

func (dispatcher *Dispatcher) dispatcher() {

	for {
		select {
		case msg := <-dispatcher.requestChan:
			dispatcher.parseMsg(msg)
		}
	}
}

func (dispatcher *Dispatcher) parseMsg(info *requestInfo) {

	// 解析body数据
	requestBody := &server.RequestMsgBody{}
	err := json.Unmarshal(info.msg, &requestBody)
	if err != nil {
		writeinfo := fmt.Sprintf("解析发生错误: %+v, 消息丢弃: %s", err, string(info.msg))
		write, _ := comm.GbkToUtf8([]byte(writeinfo))
		info.conn.Write(write)
		return
	}

	// 如果买家ID 与 卖家ID 有一个任务正在撮合,则分配到同一个进程里
	var processId int
	if val, buyOk := dispatcher.activeUserProcessMap.Load(requestBody.Md.Uid); buyOk {
		processId = val.(int)
	} else if val, saleOk := dispatcher.activeUserProcessMap.Load(requestBody.Md.BakUid); saleOk {
		processId = val.(int)
	} else {
		processId = NewTaskProcess.GetIdleProcess()
	}

	// 将数据分发
	taskData := &server.ServerContext{
		Conn: info.conn,
		Body: requestBody,
	}

	// 发送数据分发消息
	ch, ok := NewTaskProcess.goPool.Load(processId)
	if ok {
		ch.(chan *server.ServerContext) <- taskData
	}

	//fmt.Println("分配的进程: ", processId)
	// 发送成功，保存map & 计数
	if requestBody.Md.Uid > 0 {
		dispatcher.activeUserProcessMap.Store(requestBody.Md.Uid, processId)
		dispatcher.activeUserCount(requestBody.Md.Uid, 1)
	}

	if requestBody.Md.BakUid > 0 {
		dispatcher.activeUserProcessMap.Store(requestBody.Md.BakUid, processId)
		dispatcher.activeUserCount(requestBody.Md.BakUid, 1)
	}

	//dispatcher.activeUserProcessMap.Range(func(key, value interface{}) bool {
	//	fmt.Println("key: ", key, " val: ", value)
	//	return true
	//})

}

// 维护uid在 process中统计
func (dispatcher *Dispatcher) activeUserCount(uid, num int) int {
	dispatcher.lock.Lock()
	defer dispatcher.lock.Unlock()

	val, ok := dispatcher.activeUserCountMap[uid]
	if !ok {
		dispatcher.activeUserCountMap[uid] = 1
		return 1
	}

	saveVal := val + num
	dispatcher.activeUserCountMap[uid] = saveVal
	return saveVal
}
