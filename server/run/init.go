package run

import (
	"slgserver/db"
	"slgserver/net"
	"slgserver/server/controller"
	"slgserver/server/entity"
	"slgserver/server/static_conf"
)

var MyRouter = &net.Router{}

func Init() {
	db.TestDB()
	initRouter()

	static_conf.FPRC.Load()

	entity.BCMgr.Load()
	entity.NMMgr.Load()
	entity.RCMgr.Load()
	entity.RBMgr.Load()
	entity.RFMgr.Load()

}

func initRouter() {
	controller.DefaultAccount.InitRouter(MyRouter)
	controller.DefaultRole.InitRouter(MyRouter)
	controller.DefaultMap.InitRouter(MyRouter)
	controller.DefaultCity.InitRouter(MyRouter)
}