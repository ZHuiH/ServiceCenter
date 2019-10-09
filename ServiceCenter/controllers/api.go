package controllers

import(
	"ServiceCenter/common"
	"ServiceCenter/conf"
	http "github.com/valyala/fasthttp"
)

/*
*
*	API
*
*/

//获取用户的信息 
func GetMessageList(ctx *http.RequestCtx){
	token:=ctx.UserValue("token").(string);

	//找用户
	var target *conf.Client;
	conf.UserList.ForEach(func(item *conf.Client){
		if item.Info.Token==token{
			target=item
			return;
		}
	});
	var json []byte;
	//找到用户
	if target !=nil{
		history:=target.HandlerChatHistory();
		json=common.JsonEncode(history);
	}else{
		data:=common.ReturnFormat("error","msg","找不到该用户");
		json=common.JsonEncode(data);
	}

	ctx.WriteString(string(json));
}