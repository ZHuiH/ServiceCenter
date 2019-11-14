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
	//预留10个位置超过就要等

	//接收用户的信息 并返回相应的操作
	go handlerUserMessage(conn,user);
	
	//接收协程发过来的的信息
	target:=<-user;

	//处理用户发出的信息
	go handlerUserWrite(target);
	for{
		check:=headlerHeartbeat(conf.Config.Heartbeat ,&target.Info ,conn);
		
		if !check{
			//把信息推送出去
			target.PushMessage();
			return;
		}
	}
}

//处理用户发来的信息
func handlerUserMessage(conn *ws.Conn,userAgent chan *conf.Client){
	user:=conf.Client{};

	//计算短期之间是否多次频繁发送 如果超过50次就把他踢掉
	var count int = 0;
	var responseTime int64=time.Now().Unix();
	//给主协程发送注册的信号
	userAgent<-&user;
	for{
		
		_, content, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return;
		}
		//检测是不是ddos攻击
		verify:=PreventDdos(&count,&responseTime);
		//确认是ddos 直接返回 并且设置活跃时间 直接关闭
		if !verify{
			user.Info.ActiveTime=0;
			continue;
		}
		result:=Parse(content);
		//登录是单独的
		if result["type"]==conf.LOGIN{
			UserLogin(conn,&user,result["token"]);
		}else{

			//一接到信息就开协程
			go func(){
				if user.Info.Token=="" && user.Info.Connect !=nil {
					return;
				}
			//最基本的处理 
			//将内容进行json解码 并且进行信息的类型归类
				var msg map[string]string;
				switch  result["type"] {
					case	"msg"	: 	msg=broadcastToServe(&user,result);break;
					default	: 	msg=common.ReturnFormat("fail",conf.TIPS,"格式错误");
				}
	
				if  user.Serve !=nil && user.Serve.Info.Token !=""{
					msg["sendTo"]=user.Serve.Info.Token;
				}

				user.Info.Received<-msg;
			}();

		}
	}
}

func UserLogin(conn *ws.Conn,user *conf.Client,token string){
	sign:=make(map[string]string);
	msg:=common.ReturnFormat("fail",conf.LOGIN,"验证失败，token不存在");			
	//用户登录
	if len(token) > 0{
		sign=login(token,conf.Config.VerifyUrl);
	//游客登录
	}else if len(conf.Config.VerifyUrl) <= 0{
		sign=guestsLogin();
	}

	if len(sign) > 0{
		//创建用户数据
		*user=buildUserData(sign);

		message:=new([]conf.Message);
		//补充
		user.History=message;

		//尝试去插入到用户列表
		msg=appendUser(user);
		
		if msg["status"]=="success"{
			//绑定客服
			user.Serve=conf.ServiceList.MinConnect();

			if user.Serve==nil || user.Serve.Info.Connect==nil{
				if conf.ServiceList.Head !=nil{
					msg=common.ReturnFormat("success",conf.TIPS,"客服繁忙，请耐心等待");
				}else{
					msg=common.ReturnFormat("success",conf.TIPS,"当前并未有客服在线！");
				}
			}else{
				user.Serve.ConnectNumber++;
				//分配客服成功了 通知给全部客服
				userInfo:=map[string]string{
					"avatar":user.Info.Avatar,
					"nickName":user.Info.NickName,
				}
				loginTips:=map[string]string{
					"data":string(common.JsonEncode(userInfo)),
					"type":"access",
					"status":"success",
					"content":"用户登录",
				}
				broadcastToServe(user,loginTips);
			}

			jsonStr:=common.JsonEncode(sign);
			msg["data"]=string(jsonStr);
		}else{
			*user=conf.Client{};
			msg=common.ReturnFormat("fail",conf.LOGIN,"登录失败");
		}
	}
	user.Info.Connect=conn;
	user.Info.Received<-msg;
}

//由此控制写入
func handlerUserWrite(user *conf.Client){
	for{
		//为防止同一个conn 同时写入 设定为单线程
		info:=user.Info;

		if info.Connect !=nil{

			msg:=<-info.Received;
			if len(msg) > 0{
				//转换格式
				contents:=toMessage(&info,msg);
				err:=info.Connect.WriteJSON(contents);
				if err !=nil{
					user.Info.Close=true;
				}
			}
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
	msg=common.ReturnFormat("success",conf.LOGIN,"登录成功");

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
			msg["status"]="fail";
			msg[conf.MSG]="该账号已登录";
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
		NickName		:	data["nickName"],
		Token			:	data["token"],
		Close			:	false,
		Received		:	make(chan map[string]string,10),
		Avatar			:	data["avatar"],
	};
	return;
}

//将信息转发到客服
//广播到客服
func broadcastToServe(user *conf.Client,data map[string]string)(msg map[string]string){

	if user.Serve == nil || user.Serve.Info.Connect==nil{
		msg=common.ReturnFormat("fail",conf.TIPS,"发送失败,当前并未有客服");
		return;
	}else{
		msg=common.ReturnFormat("fail",conf.MSG,"发送失败");
	}

	
	//生成Message结构体的数据
	data["sentTo"]=user.Serve.Info.Token;
	content:=toMessage(&user.Info,data);

	if len(content.Content) <= 0{
		return;
	}

	history:=conf.Message{
		CreateTime	:content.CreateTime,
		Content		:content.Content,
		Source		:true,
	}
	//插入会话队列
	*user.History=append(*user.History,history);
	//更新活跃时间 防止被标识为关闭
	user.Info.ActiveTime=time.Now().Unix();
	//是时候推送信息了
	if conf.Config.ServeUnion{
		//开启联合服务模式
		conf.ServiceList.ForEach(func(item *conf.Serve){
			item.Info.Received<-data;
		})
	}else{
		user.Serve.Info.Received<-data;
	}
	msg=common.ReturnFormat("success",conf.TIPS,"发送成功");
	return;
}
