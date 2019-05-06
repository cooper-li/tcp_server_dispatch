package handler

import (
	"assets_server/depend/server"
)

type Demo struct {
	Name string `json:"name"`
	Attr string `json:"attr"`
}

// 测试方法
func SayHelloWorld(ctx *server.ServerContext) error {
	returnBody := "接收到参数: " + ctx.Body.P
	return ctx.ReString(returnBody)
}

func IsMe(ctx *server.ServerContext) error {
	me := &Demo{
		Name: "cooper",
		Attr: "特点就是帅",
	}
	return ctx.ReJson(me)

}
