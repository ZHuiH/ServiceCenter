package main
import(
	"fmt"
	"ServiceCenter/conf"
	http "github.com/valyala/fasthttp"
	router "ServiceCenter/router"
	"ServiceCenter/file"
)

func init(){
	//绑定flag
	conf.BindFlag();
}

func main(){
	//初始化
	addr:=conf.StartUp();
	//启动之前先搞一波错误日志
	if len(conf.Config.ErrorLog) > 0{
		defer func(){
			fmt.Println("error")
			if err:=recover();err !=nil{
				file.WriteFile(conf.Config.ErrorLog,err.(string),true);
			}
		}();
	}
	//用户列表开启自动清除
	go conf.UserList.ClearUp();
	//启动服务
	err:=http.ListenAndServe(addr,router.Routers.Handler);

	if(err !=nil){
		fmt.Println("启动失败：",err);
	}
}