package conf

import(
	"fmt"
	"os"
	"flag"
	"strconv"
	"strings"
	"ServiceCenter/file"
)

//绑定flag
func BindFlag(){
	flag.StringVar(&Config.Port,"p","91","访问的端口号");
	flag.StringVar(&Config.Address,"ip","0.0.0.0","访问的ip地址");
	flag.StringVar(&Config.Path,"c","","配置文件的路径");
	flag.IntVar(&Config.ExpiresTime,"t",300,"用户会话的过期时间（单位：秒）");
	flag.StringVar(&Config.ErrorLog,"error_log","","错误日志的路径");
	flag.StringVar(&Config.AccessLog,"access_log","","错误日志的路径");
	flag.StringVar(&Config.PushUrl,"push_url","","信息的推送地址");
	flag.StringVar(&Config.SessionLog,"session_log","","保存连接日志的路径");
	flag.IntVar(&Config.PushInterval,"push_interval",300,"信息推送的间隔（单位：秒）");
	flag.StringVar(&Config.VerifyUrl,"verify_url","","验证登录的地址");
	flag.StringVar(&Config.ServeVerify,"serve_verify","","客服登录验证的地址");
	flag.StringVar(&Config.GuestsName,"guests_name","游客","游客的昵称");
	flag.IntVar(&Config.Heartbeat,"heartbeat",30,"心跳的时间");
	flag.BoolVar(&Config.ServeUnion,"serve_union",true,"开启联合服务模式");
	flag.IntVar(&Config.ServeConcurrent,"serve_concurrent",10,"客服最大连接数量");
	flag.BoolVar(&Config.Env,"env",false,"查看配置");
}

//启动
func StartUp() (addr string) {
	flag.Parse();

	//读取配置的优先级最高
	if len(Config.Path) > 0{
		ReadConf();
	}
	//最后才是查看配置
	if Config.Env {
		PrintEnv();
		os.Exit(0);
	}

	addr=Config.Address+":"+Config.Port;
	return;
}

//打印配置
func PrintEnv(){
	fmt.Println("配置文件的路径:",Config.Path);
	fmt.Println("IP地址：",Config.Address);
	fmt.Println("端口号：",Config.Port);
	fmt.Println("用户会话的过期时间:",Config.ExpiresTime);
	fmt.Println("保存会话记录的路径:",Config.SessionLog);
	fmt.Println("保存连接日志的路径:",Config.AccessLog);
	fmt.Println("保存错误日志的路径:",Config.ErrorLog);
	fmt.Println("验证登录的地址:",Config.VerifyUrl);
	fmt.Println("客服登录验证的地址:",Config.ServeVerify);
	fmt.Println("信息的推送地址:",Config.PushUrl);
	fmt.Println("信息推送的间隔:",Config.PushInterval);
	fmt.Println("游客的昵称:",Config.GuestsName);
	fmt.Println("心跳时间:",Config.Heartbeat);
	fmt.Println("联合服务模式:",Config.ServeUnion);
	fmt.Println("客服最大连接数:",Config.ServeConcurrent);
	
}

//读取文件的每一行
func ReadConf(){
	err:=file.ReadLine(Config.Path,func(line string){
		//处理参数
		index:=strings.Index(line,"=");
		if index >0{
			value:=line[index+1:len(line)];
			key:=line[:index];
			switch key {
				case "address" 			: Config.Address=value;break;
				case "port"				: Config.Port=value;break;
				case "expires"			: Config.ExpiresTime,_=strconv.Atoi(value);break;
				case "error_log"		: Config.ErrorLog=value;break;
				case "access_log"		: Config.AccessLog=value;break;
				case "push_url"			: Config.PushUrl=value;break;
				case "push_interval"	: Config.PushInterval,_=strconv.Atoi(value);break;
				case "verify_url"		: Config.VerifyUrl=value;break;
				case "guests_name"		: Config.GuestsName=value;break;
				case "serve_verify"		: Config.ServeVerify=value;break;
				case "heartbeat"		: Config.Heartbeat,_=strconv.Atoi(value);break;
				case "serve_concurrent"	: Config.ServeConcurrent,_=strconv.Atoi(value);break;
				case "serve_union"		: Config.ServeUnion,_=strconv.ParseBool(value);break;
			}
		}
	})

	if err !=nil{
		fmt.Println("读取配置文件错误",err);
		os.Exit(0);
	}
}

