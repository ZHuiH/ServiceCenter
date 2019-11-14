package controllers

import (
	"ServiceCenter/common"
	"ServiceCenter/conf"
	"fmt"
	"time"

	ws "github.com/fasthttp/websocket"
	http "github.com/valyala/fasthttp"
)

//客服的话就不要搞那么多事了 直接搞起
//客服必须登录验证
func Serve(ctx *http.RequestCtx) {
	//升级协议
	err := conf.ChatConf.Upgrade(ctx, handlerServe)

	if err != nil {
		return;
	}

}

//处理客服的数据
func handlerServe(conn *ws.Conn) {

	defer conn.Close();

	userAgent := make(chan *conf.Serve, 1)
	go handlerServeMessage(conn, userAgent)
	createTime := time.Now().Unix()
	user := <-userAgent
	go handlerServeWrite(user);
	for {
		time.Sleep(30 * time.Second)
		now := time.Now().Unix()
		//关闭
		if user.Info.Close {
			verfiy := (&user.Info).Sort()
			if verfiy {
				return
			}
			//没有登录才需要验证
		} else if user.Info.Id == "" {
			if (now - createTime) >= int64(30) {
				return;
			}
		}
		user.Info.Received<-common.ReturnFormat("success", conf.PING, "")
	}
}

func handlerServeMessage(conn *ws.Conn, userAgent chan *conf.Serve) {
	user := conf.Serve{}
	userAgent <- &user
	for {

		_, content, err := conn.ReadMessage()
		//读取失败 基本上就是已经断开了
		if err != nil {
			fmt.Println(err)
			return
		}

		result := Parse(content)

		token, ok := result["token"]
		//预设
		msg := common.ReturnFormat("fail", conf.TIPS, "检测到未登录")

		sendTo, _ := result["sendTo"];
		if ok {

			if result["type"] == "login" {
				msg = common.ReturnFormat("fail", conf.TIPS, "验证失败，token不存在")
				sign := login(token, conf.Config.ServeVerify)
				msg = sign
				//验证登录成功,开始注册信息
				if len(sign) > 0 && sign["status"] == "success" {
					user = buildServe(sign)
					//补充
					user.Info.Connect = conn
					//开始注册
					register := registerServe(&user)

					if register["status"] == "success" {
						//登录成功了
						msg = common.ReturnFormat("success", conf.TIPS, "登录验证成功");
					}
				}
				user.Info.Received<-msg;
			}else{
				go func(){
					switch result["type"] {
					case "msg":
						msg = sendToUser(sendTo, result)
						break
					default:
						msg = common.ReturnFormat("fail", conf.TIPS, "格式错误")
					}
					msg["sendTo"]=sendTo;
					user.Info.Received<-msg;
				}()
			}
		}
	}
}


//这里原则上来说是跟客户的函数一样的
//但是防止扩展 分开
func handlerServeWrite(user *conf.Serve){
	for{
		info:=user.Info;
		if info.Token !="" && info.Connect !=nil{
			msg:=<-info.Received;
			//转换格式
			contents:=toMessage(&info,msg);
			err:=info.Connect.WriteJSON(contents);
			if err!=nil{
				info.Close=true;
			}
		}
	}
}

//注册服务
func registerServe(user *conf.Serve) (msg map[string]string) {
	head := conf.ServiceList.Head
	msg = common.ReturnFormat("success", conf.TIPS, "注册成功")
	//没人登录
	if head == nil {
		conf.ServiceList.Total++
		conf.ServiceList.Head = user
		return
	}

	var target *conf.Serve
	conf.ServiceList.ForEach(func(item *conf.Serve) {
		if item.Info.Id == user.Info.Id {
			//更新token
			if item.Info.Token != user.Info.Token {
				item.Info.Token = user.Info.Token
				item.Info.Close = false
			}
			return
		}
		target = item
	})
	//防止目前登录的就是列表的最后一名
	verify := target.Info.Id == user.Info.Id

	//没有找到
	if target.Next == nil && !verify {
		target.Next = user
		conf.ServiceList.Total++
	} else {
		//登录过了
		msg["status"] = "fail"
		msg["msg"] = "该账号已登录"
	}

	return
}

//生成数据
func buildServe(data map[string]string) conf.Serve {
	createTime := time.Now().Unix()
	return conf.Serve{
		Info: conf.UserInfo{
			CreateTime: createTime,
			ActiveTime: createTime,
			Id:         data["id"],
			NickName:   data["nike_name"],
			Avatar:     "www.baidu.com",
			Close:      false,
			Received:	make(chan map[string]string,10),
			Token:      data["token"],
		},
		ConnectNumber: 0,
		//History				:	nil,
		AnswerNumber: 0,
		//FastestAnswerTime	:	0,
		//SlowestAnswerTime	:	0,
		Next: nil,
	}
}

//将信息转发到指定的客户
func sendToUser(token string, data map[string]string) (msg map[string]string) {
	if token == "" {
		msg = common.ReturnFormat("fail", conf.TIPS, "用户token不能为空")
		return
	}

	var target *conf.Client
	conf.UserList.ForEach(func(item *conf.Client) {
		if item.Info.Token == token {
			target = item
		}
	})

	if target == nil {
		msg = common.ReturnFormat("fail", conf.TIPS, "该用户已下线")
		return
	}

	//生成Message结构体的数据
	data["source"]="1";
	content := toMessage(&target.Info, data)
	if len(content.Content) <= 0 {
		return
	}

	history := conf.Message{
		CreateTime: content.CreateTime,
		Content:    content.Content,
		Source:     false,
	}

	//插入会话队列
	*target.History = append(*target.History, history)
	//就算是客服跟他说也要更新活跃时间 防止被标识为关闭
	target.Info.ActiveTime = time.Now().Unix()

	var feedback string = ""
	if _, ok := data["feedback"]; ok {
		feedback = data["feedback"]
	}

	//推送信息
	target.Info.Connect.WriteJSON(content)

	msg = common.ReturnFormat("success", conf.SENT, "发送成功", feedback);

	return
}
