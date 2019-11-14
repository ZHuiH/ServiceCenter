package controllers

import (
	"ServiceCenter/common"
	"ServiceCenter/conf"
	"ServiceCenter/request"
	"time"
	//"fmt"
	ws "github.com/fasthttp/websocket"
)

//解析
func Parse(content []byte) (result map[string]string) {

	//json解码
	result = common.JsonDecode(content)

	if len(result) <= 0 {
		result = make(map[string]string);
	}

	_, nothink := result["type"];
	if !nothink {
		result["type"] = "nothink";
	}
	return
}

//用户登录
func login(token string, url string) (msg map[string]string) {

	args := make(map[string]string);
	args["token"] = token;
	msg = request.Post(url, args);

	// if result["status"]=="success"{

	// 	if _,ok:=result["id"];ok{
	// 		msg=result;
	// 	}
	// }

	return
}

//处理心跳的逻辑
func headlerHeartbeat(duration int, target *conf.UserInfo, conn *ws.Conn) bool {
	heartbeat := time.Duration(duration)

	//心跳的间隔
	time.Sleep(heartbeat * time.Second)

	if target.Close {
		return false;
	} else {
		now := time.Now().Unix()
		//连接超过30秒 没有注册自动 干掉
		if (now - target.ActiveTime) >= int64(30) {

			msg := common.ReturnFormat("error",conf.TIPS,"连接失效，请及时注册");
			target.Received<-msg;
			return false;
		}
	}

	target.Received<-common.ReturnFormat("success",conf.PING, "");
	return true;
}

//将map数据变更为为json string并进行最低程度的检测
func toMessage(info *conf.UserInfo, data map[string]string) (content conf.SendFormat) {
	//必须有内容
	if _, ok := data["content"]; !ok {
		return
	}

	msgType := "msg"
	//类型如果没有就默认msg
	if _, check := data["type"]; check {
		msgType = data["type"]
	}

	otherData := ""
	//是否存在附加的数据
	if _, check := data["data"]; check {
		otherData = data["data"];
	}
	status := "error"
	if _, check := data["status"]; check {
		status = data["status"];
	}
	//发送给谁
	sendTo := "";
	if _, check := data["sendTo"]; check {
		sendTo = data["sendTo"];
	}

	source := true;
	if _, check := data["source"]; check {
		source = false;
	}

	//生成指定的格式
	content = conf.SendFormat{
		CreateTime: time.Now().Unix(),
		Content:    data["content"],
		Source:     source,
		From:      	info.Token,
		To			:sendTo,
		Status:     status,
		Type:       msgType,
		Data:       otherData,
	}
	return content
}

//阻止ddos
func PreventDdos(count *int, responseTime *int64) bool {
	//超过50次直接踢出去
	if *count >= 50 {
		return false
	}

	now := time.Now().Unix()
	//重置
	if now-*responseTime >= 30 {
		*count = 0
		*responseTime = now
	} else {
		*count++
	}

	return true

}
