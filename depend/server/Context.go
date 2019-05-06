package server

import (
	"net"
	"encoding/json"
)

type (
	ServerContext struct {
		Conn *net.TCPConn
		Body *RequestMsgBody
	}
)

// JSON 返回
func (ctx *ServerContext) ReJson(i interface{}) error {
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}

	_, err = ctx.Conn.Write(b)
	return err
}

// 字符串返回
func (ctx *ServerContext) ReString(s string) error {
	_, err := ctx.Conn.Write([]byte(s))
	return err
}
