package controllers
import(
	"fmt"
	http "github.com/valyala/fasthttp"
	ws "github.com/fasthttp/websocket"
	"ServiceCenter/conf"
	"ServiceCenter/common"
	"ServiceCenter/file"
	"time"
	"strconv"
)
//客户端 入口
func Chat(ctx *http.RequestCtx){

	//如果设置了连接日志 那就写入日志
	if len(conf.Config.AccessLog) > 0{
		//请求的url包括参数
		content	:=	"RequestURI:"+string(ctx.RequestURI());
		//客户端的ip
		content+="; CilentIp:"+fmt.Sprintf("%s",ctx.RemoteIP());
		//客户端的信息
		content+="; UserAgent:"+string(ctx.UserAgent());
		file.WriteFile(conf.Config.AccessLog,content,true);
	}
	//升级协议
	err:=conf.ChatConf.Upgrade(ctx,handlerUser);
	
	if err !=nil{
		return;
	}

}
//处理客户连接信息
func handlerUser(conn *ws.Conn){
	defer conn.Close();
	user:=make(chan *conf.Client, 1);
	//createTime:=time.Now().Unix();
	//接收用户的信息 并返回相应的操作
	go handlerUserMessage(conn,user);
	//接收协程发过来的的信息
	temp:=<-user;

	for{

		check:=headlerHeartbeat(conf.Config.Heartbeat ,&temp.Info ,conn);

		if !check{
			//把信息推送出去
			temp.PushMessage();
			return;
		}
	}
}

//处理用户发来的信息
func handlerUserMessage(conn *ws.Conn,userAgent chan *conf.Client){
	user:=conf.Client{};

	//给主协程发送注册的信号
	userAgent<-&user;
	for{
		index, content, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return;
		}
		//最基本的处理 
		//将内容进行json解码 并且进行信息的类型归类
		result:=Parse(content);

		//返回信息
		msg:=common.ReturnFormat("error","login","验证失败，token不存在");
		//登录是单独的
		if result["type"]=="login"{
			sign:=make(map[string]string);
			
			//用户登录
			if len(result["token"]) > 0{
				sign=login(result["token"],conf.Config.VerifyUrl);
			//游客登录
			}else if len(conf.Config.VerifyUrl) <= 0{
				sign=guestsLogin();
			}

			if len(sign) > 0{
				//创建用户数据
				user=buildUserData(sign);

				message:=new([]conf.Message);
				//补充
				user.History=message;
				user.Info.Connect=conn;
				
				//尝试去插入到用户列表
				msg=appendUser(&user);
				
				if msg["status"]=="success"{
					//绑定客服
					user.Serve=conf.ServiceList.MinConnect();
					if user.Serve==nil{
						if conf.ServiceList.Head !=nil{
							msg=common.ReturnFormat("success","tips","客服繁忙，请耐心等待");
						}else{
							msg=common.ReturnFormat("success","tips","当前并未有客服在线！");
						}
						
					}

					jsonStr:=common.JsonEncode(sign);
					msg["data"]=string(jsonStr);
				}else{
					user=conf.Client{};
					msg=common.ReturnFormat("error","login","登录失败");
				}
			}
		}

		switch  result["type"] {
			case	"login"	: 	break;
			case	"msg"	: 	msg=broadcastToServe(&user,result);break;
			default	: msg=common.ReturnFormat("error","tips","格式错误");
		}

		json:=common.JsonEncode(msg);

		if err := conn.WriteMessage(index, json); err != nil {
			fmt.Println(err);
			return;
		}
	}
}

//游客登录
func guestsLogin()(info map[string]string){
	if conf.UserList.Head==nil{
		conf.GuestsTotal=1;
	}else{
		conf.GuestsTotal++;
	}
	
	//伪造数据
	info=make(map[string]string);
	info["nick_name"]=conf.Config.GuestsName;
	info["id"]=strconv.Itoa(conf.GuestsTotal);
	info["token"]=common.RandStr(32);
	return;
}

//添加用户
func appendUser(user *conf.Client)(msg map[string]string){
	head:=conf.UserList.Head;
	msg=common.ReturnFormat("success","login","登录成功");

	//首个 列表的开头
	if head == nil{
		conf.UserList.Head=user;
		conf.UserList.Total++;
	}else{
		var target *conf.Client;
		//遍历查找
		conf.UserList.ForEach(func(item *conf.Client){

			if user.Info.Id == item.Info.Id{
				//id正确 token对不上那就更新token
				if user.Info.Token != item.Info.Token{
					item.Info.Token=user.Info.Token;
				}
				return;
			};
			target=item;
		})
		//防止列表当前登录的就是最后一个用户
		verify:=target.Info.Id==user.Info.Id;

		//如果到了尾部 就是没有 注册
		if target.Next ==nil && !verify{
			target.Next=user;
			conf.UserList.Total++;
		}else{
			//登录过了
			msg["status"]="error";
			msg["msg"]="该账号已登录";
		}
	}

	return;
}

//返回创建的数据
func buildUserData(data map[string]string)(user conf.Client){
	user=conf.Client{
		Info:conf.UserInfo{},
		History:nil,
		Next:nil,
	};
	timestamp:=time.Now().Unix();
	user.Info=conf.UserInfo{
		CreateTime		:	timestamp,
		ActiveTime		:	timestamp,
		Id				:	data["id"],
		NickName		:	data["nick_name"],
		Token			:	data["token"],
		Close			:	false,
		ProfilePicture	:	"https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=806928800,2548687982&fm=26&gp=0.jpg",
	};
	return;
}

//将信息转发到客服
//广播到客服
func broadcastToServe(user *conf.Client,data map[string]string)(msg map[string]string){

	if user.Serve == nil{
		msg=common.ReturnFormat("error","tips","发送失败,当前并未有客服");
		return;
	}else{
		msg=common.ReturnFormat("error","msg","发送失败");
	}
	
	//生成Message结构体的数据
	content:=toMessage(&user.Info,data,true);
	if len(content.Content) <= 0{
		return;
	}
	//插入会话队列
	*user.History=append(*user.History,content);
	//更新活跃时间 防止被标识为关闭
	user.Info.ActiveTime=time.Now().Unix();
	
	//是时候推送信息了
	if conf.Config.ServeUnion{
		//开启联合服务模式
		conf.ServiceList.ForEach(func(item *conf.Serve){
			item.Info.Connect.WriteJSON(content);
		})
	}else{
		user.Serve.Info.Connect.WriteJSON(content);
	}
	msg=common.ReturnFormat("success","msg","发送成功");
	return;
}