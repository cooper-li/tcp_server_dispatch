package conf

import (
	"fmt"
	"time"
	"os"
	"encoding/json"
	"sync"
	"runtime"

	"os_go_comm/comm_log"
	"os_go_comm/comm_cfg"

	"github.com/parnurzeal/gorequest"
	"github.com/labstack/gommon/log"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type (
	// 检查信息
	ApiConf struct {
		CheckInterval int    `json:"check_interval"`
		Api           string `json:"api"`
	}

	// 响应信息
	ReponseConf struct {
		Status int    `json:"status"`
		Msg    string `json:"msg"`
		Data   []struct {
			CfgID     string      `json:"cfgId"`
			Namespace string      `json:"namespace"`
			GroupName string      `json:"groupName"`
			DataID    string      `json:"dataId"`
			TplName   string      `json:"tplName"`
			CreateAt  int         `json:"createAt"`
			Content   interface{} `json:"content"`
		} `json:"data"`
		Time      int    `json:"time"`
		Microtime int64  `json:"microtime"`
		Source    string `json:"source"`
	}

	// 通用信息
	CommConf struct {
		// 配置信息
		ConfName string `json:"conf_name,omitempty"` // 配置名称
		CreateAt int    `json:"create_at,omitempty"` // 配置时间, 用来校验是否需要更新
	}

	// mysql 配置
	MysqlConf struct {
		CommConf
		UserName      string `json:"user_name"`
		Password      string `json:"password"`
		ServerAddress string `json:"server_address"`
		DBName        string `json:"db_name"`
		MaxIdleConn   int    `json:"max_idle_conn"`
		MaxOpenConn   int    `json:"max_open_conn"`
	}

	// redis配置
	RedisConf struct {
		CommConf
		// 模板信息
		ServerAddr        string `json:"server_addr"`
		Auth              string `json:"auth,omitempty"`
		ConnTimeout       int    `json:"conn_timeout"`
		ReadTimeout       int    `json:"read_timeout"`
		WriteTimeout      int    `json:"write_timeout"`
		MaxIdle           int    `json:"max_idle"`
		MaxActive         int    `json:"max_active"`
		IdleTimeout       int    `json:"idle_timeout"`
		HeartBeatMin      int    `json:"heart_beat_min"`
		HeartBeatMax      int    `json:"heart_beat_max"`
		HeartBeatInterval int    `json:"heart_beat_interval"`
		MaxWait           bool   `json:"max_wait"`
	}

	DBContainer struct {
		lock    sync.Mutex
		MysqlDB map[string]*gorm.DB // db容器
		DBToken map[string]int      // dataId => timestamp, 校验是否需要更新
		RedisDB map[string]*redis.Pool
	}
)

var (
	conf         *ApiConf
	DBContainers *DBContainer
)

func init() {
	DBContainers = &DBContainer{
		lock:    sync.Mutex{},
		MysqlDB: map[string]*gorm.DB{},
		DBToken: map[string]int{},
		RedisDB: map[string]*redis.Pool{},
	}

	conf = &ApiConf{
		CheckInterval: comm_cfg.Int("database_api", "check_interval"),
		Api:           comm_cfg.GetValue("database_api", "assets_api"),
	}

}

// 初始化db配置
func InitDatabase() {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 200)
			runtime.Stack(buf, false)
			comm_log.Error("sync_database_panic", "err", err, "stack", buf)
			os.Exit(1)
		}
	}()

	syncDBConf()
	timer := time.NewTicker(time.Duration(conf.CheckInterval) * time.Second)
	for {
		select {
		case <-timer.C:
			syncDBConf()
		}
	}

}

func syncDBConf() {

	_, res, _ := gorequest.New().Get(conf.Api).End()

	var confResp *ReponseConf
	err := json.Unmarshal([]byte(res), &confResp)
	if err != nil {
		log.Fatal("sync database conf fail :", err)
	}

	for _, dbInfo := range confResp.Data {

		// 是否需要更新
		isUp := DBContainers.compareToken(dbInfo.DataID, dbInfo.CreateAt)
		if isUp {
			continue
		}

		// 通用配置
		commConf := CommConf{
			ConfName: dbInfo.DataID,
			CreateAt: dbInfo.CreateAt,
		}

		switch item := dbInfo.Content.(type) {
		case *RedisConf:
			poolCfg := &RedisConf{
				CommConf:          commConf,
				ServerAddr:        dbInfo.Content.(*RedisConf).ServerAddr,
				Auth:              dbInfo.Content.(*RedisConf).Auth,
				ConnTimeout:       dbInfo.Content.(*RedisConf).ConnTimeout,
				ReadTimeout:       dbInfo.Content.(*RedisConf).ReadTimeout,
				WriteTimeout:      dbInfo.Content.(*RedisConf).WriteTimeout,
				MaxIdle:           dbInfo.Content.(*RedisConf).MaxIdle,
				MaxActive:         dbInfo.Content.(*RedisConf).MaxActive,
				IdleTimeout:       dbInfo.Content.(*RedisConf).IdleTimeout,
				HeartBeatMin:      dbInfo.Content.(*RedisConf).HeartBeatMin,
				HeartBeatMax:      dbInfo.Content.(*RedisConf).HeartBeatMax,
				HeartBeatInterval: dbInfo.Content.(*RedisConf).HeartBeatInterval,
				MaxWait:           dbInfo.Content.(*RedisConf).MaxWait,
			}
			//DBContainers.RedisDB[dbInfo.DataID] = NewRedisPool(poolCfg)
			DBContainers.setDBContainer(dbInfo.DataID, NewRedisPool(poolCfg), dbInfo.CreateAt)
			fmt.Printf("init database redis success, name: %s, addr: %s \n", dbInfo.DataID, item.ServerAddr)
		case *MysqlConf:
			poolCfg := &MysqlConf{
				CommConf:      commConf,
				UserName:      dbInfo.Content.(*MysqlConf).UserName,
				Password:      dbInfo.Content.(*MysqlConf).Password,
				ServerAddress: dbInfo.Content.(*MysqlConf).ServerAddress,
				DBName:        dbInfo.Content.(*MysqlConf).DBName,
				MaxIdleConn:   dbInfo.Content.(*MysqlConf).MaxIdleConn,
				MaxOpenConn:   dbInfo.Content.(*MysqlConf).MaxOpenConn,
			}
			//DBContainers.MysqlDB[dbInfo.DataID] = NewMysqlPool(poolCfg)
			DBContainers.setDBContainer(dbInfo.DataID, NewMysqlPool(poolCfg), dbInfo.CreateAt)
			fmt.Printf("init database mysql success, name: %s, addr: %s \n", dbInfo.DataID, item.ServerAddress)
		default:
			fmt.Println("无效配置: ", dbInfo)

		}
	}
}

// 保存DB
func (p *DBContainer) setDBContainer(dataId string, DB interface{}, token int) {
	p.lock.Lock()
	defer p.lock.Unlock()

	switch DB.(type) {
	case *gorm.DB:
		p.MysqlDB[dataId] = DB.(*gorm.DB)
	case *redis.Pool:
		p.RedisDB[dataId] = DB.(*redis.Pool)
	default:
		panic(errors.New(fmt.Sprintf("error db type, dataId=%s, token=%d", dataId, token)))
	}
	p.DBToken[dataId] = token
}

// 获取mysql数据库
func (p *DBContainer) GetMysql(dataId string) *gorm.DB {
	p.lock.Lock()
	defer p.lock.Unlock()

	if mdb, ok := p.MysqlDB[dataId]; ok {
		return mdb
	}
	return nil
}

// 获取redis数据库
func (p *DBContainer) GetRedis(dataId string) *redis.Pool {
	p.lock.Lock()
	defer p.lock.Unlock()

	if rdb, ok := p.RedisDB[dataId]; ok {
		return rdb
	}
	return nil
}

// 是否需要更新 false 需要  true 不需要
func (p *DBContainer) compareToken(dataId string, token int) bool {
	p.lock.Lock()
	defer p.lock.Unlock()

	if oldT, ok := p.DBToken[dataId]; ok {
		return oldT == token
	}
	return false
}
