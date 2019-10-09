package router
import(
	route "github.com/buaazp/fasthttprouter"
	ctrl "ServiceCenter/controllers"
)

var Routers *route.Router=route.New();
//路由的初始化 全部在这里
func init(){
	//客户聊天入口
	Routers.GET("/chat",ctrl.Chat)

	//客服登录入口
	Routers.GET("/serve",ctrl.Serve)

	//获取用户的信息列表 
	Routers.GET("/message/:token",ctrl.GetMessageList)
}