package controllers

import(
	"ServiceCenter/common"
	"ServiceCenter/request"
	"ServiceCenter/conf"
	"time"
	//"fmt"
	ws "github.com/fasthttp/websocket"
)
//解析
func Parse(content []byte)(result map[string]string){

	//json解码
	result=common.JsonDecode(content);

	if len(result) <= 0 {
		result=make(map[string]string);
	}

	_,nothink:=result["type"];
	if !nothink {
		result["type"]="nothink";
	}
	return ;
}

//用户登录
func login(token string,url string)(msg map[string]string){

	args:=make(map[string]string);
	args["token"]=token;
	result:=request.Post(url,args);

	if result["status"]=="success"{

		if _,ok:=result["id"];ok{
			msg=result;
		}
	}
	

	return;
}

//处理心跳的逻辑
func headlerHeartbeat(duration int,target *conf.UserInfo,conn  *ws.Conn) bool {
	heartbeat:=time.Duration(duration);
	//心跳的间隔
	time.Sleep(heartbeat * time.Second);

	if target.Close {
		return false;
	}else{
		now:=time.Now().Unix();
		//连接超过30秒 没有注册自动 干掉
		if (now - target.ActiveTime) >= int64(30){
			
			msg:=common.ReturnFormat("error","连接失效，请及时注册");
			conn.WriteJSON(msg);
			return false;
		}
	}
	
	return true;
}

//将map数据变更为为json string并进行最低程度的检测
func toMessage(info *conf.UserInfo,data map[string]string,source bool) (content conf.Message) {
	//必须有内容
	if _,ok:=data["content"];!ok{
		return;
	}

	data["id"]=info.Id;
	data["nike_name"]=info.NickName;
	data["profile_picture"]=info.ProfilePicture;
	json:=common.JsonEncode(data);
	//不管是什么只要转发就行了
	message:=string(json);

	content=conf.Message{
		CreateTime:time.Now().Unix(),
		Content:message,
		Source:source,
	};
	return content;
}