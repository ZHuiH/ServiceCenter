package router

import (
	ctrl "ServiceCenter/controllers"

	route "github.com/buaazp/fasthttprouter"
)

var Routers *route.Router = route.New()

//路由的初始化 全部在这里
func init() {
	//客户聊天入口
	Routers.GET("/chat", ctrl.Chat)

	//客服登录入口
	Routers.GET("/serve", ctrl.Serve)

	//获取用户的信息列表
	Routers.GET("/message/:token", ctrl.GetMessageList)

	//上传图片
	Routers.POST("/upload", ctrl.Upload)

	//跨域
	Routers.OPTIONS("/upload", ctrl.Core)

	//上传图片
	Routers.GET("/static/:date/:file", ctrl.Static)
}
