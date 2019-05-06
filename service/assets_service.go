package service

import (
	"errors"
	"assets_server/depend/precision"

	. "assets_server/depend/comm"
	"os_go_comm/redis_client"
	"github.com/garyburd/redigo/redis"
	"fmt"
)

type (
	AssetsService struct{}

	// 冻结用户资产
	FrozenAssestParam struct {
		Coin       string  `json:"coin"`
		Uid        int     `json:"uid"`
		Number     float64 `json:"number"`
		AssestType int     `json:"assest_type"`
	}
)

const (
	ASSEST_FLOW_MOD = 64

	REDIS_KEY_USER_ASSETS_INFO_HASH = "user:assets:%d:hash" // 用户资产流水hash, %s: uid,   key: coin  val: (over):(lock)
	REDIS_KEY_QUEUE_ASSEST_FLOW     = "queue:asset:flow:%d" // 用户资产流水队列, %d: uid % mod
)

func (s *AssetsService) FrozenAssets(p *FrozenAssestParam) error {

	if precision.Compare(0, p.Number) != 0 {
		return errors.New(Spf("coin_num lt 0, coin=%s, num=%f", p.Coin, p.Number))
	}

	r := redis_client.SelectDB("assets")
	defer r.Close()

	userAssestKey := Spf(REDIS_KEY_USER_ASSETS_INFO_HASH, p.Uid)


	assetsInfo, err := redis.StringMap(r.Do("HGETALL", userAssestKey))
	// todo(。。。)
	fmt.Println(assetsInfo, err)

	return nil
}
