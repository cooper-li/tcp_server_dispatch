package server

import (
	"errors"
)

type FuncHandlerMap struct {
	ServerList map[string]func(ctx *ServerContext) error
}

var ServerHandler *FuncHandlerMap

// 注册服务
func init() {
	ServerHandler = &FuncHandlerMap{
		ServerList: map[string]func(ctx *ServerContext) error{},
	}
}

// 注册服务
func (f *FuncHandlerMap) RegisterService(serverName string, fn func(ctx *ServerContext) error) {
	f.ServerList[serverName] = fn
}

// 查找服务
func (f *FuncHandlerMap) IsServer(serverName string) bool {
	_, ok := f.ServerList[serverName]
	return ok
}

// 获取服务
func (f *FuncHandlerMap) getServer(serverName string) (func(ctx *ServerContext) error, error) {
	fn, ok := f.ServerList[serverName]
	if !ok {
		return nil, errors.New("No such server func: " + serverName)
	}
	return fn, nil
}

//func Call(m map[string]interface{}, name string, params ...interface{}) ([]reflect.Value, error) {
//	f := reflect.ValueOf(m[name])
//	if len(params) != f.Type().NumIn() {
//		return nil, errors.New("the number of input params not match!")
//	}
//	in := make([]reflect.Value, len(params))
//	for k, v := range params {
//		in[k] = reflect.ValueOf(v)
//	}
//	return f.Call(in), nil
//}

// 执行
func (f *FuncHandlerMap) Call(ctx *ServerContext) error {

	// 查找方法
	fn, err := f.getServer(ctx.Body.M)
	if err != nil {
		return err
	}

	return fn(ctx)
}
