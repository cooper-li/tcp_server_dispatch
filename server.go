package main

import (
	"assets_server/control"
	"os_go_comm/comm_builder"
	"os_go_comm/comm_base"
	"os_go_comm/comm_log"
	"os_go_comm/comm_cfg"
)

var (
	PROGRAM_VERSION  string
	COMPILER_VERSION string
	BUILD_TIME       string
	AUTHOR           string
)

func main()  {
	comm_builder.Banner_show(PROGRAM_VERSION, COMPILER_VERSION, BUILD_TIME, AUTHOR)

	main_name := comm_base.Get_main_name()
	comm_log.Init(main_name, comm_cfg.GetValue("log", "level"))
	defer comm_log.Sync()

	control.MainControl()
}