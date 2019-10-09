package controllers

import(
	"fmt"
	"time"
	http "github.com/valyala/fasthttp"
	ws "github.com/fasthttp/websocket"
	"ServiceCenter/conf"
	"ServiceCenter/common"
)
//客服的话就不要搞那么多事了 直接搞起
//客服必须登录验证
func Serve(ctx *http.RequestCtx){
		//升级协议
		err:=conf.ChatConf.Upgrade(ctx,handlerServe);
		
		if err !=nil{
			return;
		}
	
}

//处理客服的数据
func handlerServe(conn *ws.Conn){

	defer conn.Close();
	
	userAgent:=make(chan *conf.Serve,1);
	go handlerServeMessage(conn,userAgent)
	createTime:=time.Now().Unix();
	user:=<-userAgent;

	for{
		now:=time.Now().Unix();
		//关闭
		if user.Info.Close{
			return;
		//没有登录才需要验证
		}else if user.Info.Id == ""{
			if (now-createTime) <= int64(30){
				return;
			}
		}
	}
}

func handlerServeMessage(conn *ws.Conn,userAgent chan *conf.Serve){
	user :=conf.Serve{};
	userAgent<-&user
	for{

		index,content,err:=conn.ReadMessage();

		if err != nil{
			fmt.Println(err)
			return;
		}
		
		result:=Parse(content);

		token,ok:=result["token"];
		
		msg:=common.ReturnFormat("error","tips","检测到未登录");
		if ok{

			if result["type"]=="login"{
				msg=common.ReturnFormat("error","tips","验证失败，token不存在");
				sign:=login(token,conf.Config.ServeVerify);

				//验证登录成功,开始注册信息
				if len(sign) > 0 {
					user=buildServe(sign);
					//补充
					user.Info.Connect=conn;
					//开始注册
					register:=registerServe(&user);

					if register["status"]=="success"{
						//登录成功了
						msg=common.ReturnFormat("success","tips","登录验证成功");
					}
				}
			}
			clientToken,_:=result["client_token"];
			switch result["type"] {
				case 	"login"	:	break;
				case 	"msg"	:	msg=sendToUser(clientToken,result);break;
				default	: msg=common.ReturnFormat("error","tips","格式错误");
			}
		}

		json:=common.JsonEncode(msg);
		if err := conn.WriteMessage(index, json); err != nil {
			fmt.Println(err);
			return;
		}
	}
}

//注册服务
func registerServe(user *conf.Serve)(msg map[string]string){
	head:=conf.ServiceList.Head;
	msg=common.ReturnFormat("success","tips","注册成功");
	//没人登录
	if head ==nil{
		conf.ServiceList.Total++;
		conf.ServiceList.Head=user;
		return
	}

	var target *conf.Serve;
	conf.ServiceList.ForEach(func(item *conf.Serve){
		if item.Info.Id == user.Info.Id{
			//更新token
			if item.Info.Token !=  user.Info.Token{
				item.Info.Token=user.Info.Token;
			}
			return;
		}
		target=item;
	})
	//防止目前登录的就是列表的最后一名
	verify:=target.Info.Id==user.Info.Id;

	//没有找到
	if target.Next == nil && !verify{
		target.Next=user;
		conf.ServiceList.Total++;
	}else{
		//登录过了
		msg["status"]="error";
		msg["msg"]="该账号已登录";
	}

	return;
}

//生成数据
func buildServe(data map[string]string)(conf.Serve){
	createTime:=time.Now().Unix();
	return conf.Serve{
		Info:conf.UserInfo{
			CreateTime		:	createTime,
			ActiveTime		:	createTime,
			Id				:	data["id"],
			NickName		:	data["nike_name"],
			ProfilePicture	:	"www.baidu.com",
			Close			:	false,
			Token			:	data["token"],
		},
		ConnectNumber		:	0,
		//History				:	nil,
		AnswerNumber		:	0,
		//FastestAnswerTime	:	0,
		//SlowestAnswerTime	:	0,
		Next				:	nil,
	}
}

//将信息转发到指定的客户
func sendToUser(token string,data map[string]string)(msg map[string]string){
	if token ==""{
		msg=common.ReturnFormat("error","tips","用户token不能为空")
		return;
	}

	var target *conf.Client
	conf.UserList.ForEach(func(item *conf.Client){
		if item.Info.Token== token{
			target=item;
		}
	})

	if target ==nil{
		msg=common.ReturnFormat("error","tips","该用户已下线")
		return;
	}
	
	//生成Message结构体的数据
	content:=toMessage(&target.Info,data,false);
	if len(content.Content) <= 0{
		return;
	}
	//插入会话队列
	*target.History=append(*target.History,content);
	//就算是客服跟他说也要更新活跃时间 防止被标识为关闭
	target.Info.ActiveTime=time.Now().Unix();
	
	//推送信息
	target.Info.Connect.WriteJSON(content);
	
	msg=common.ReturnFormat("success","msg","发送成功");
	return;
}



