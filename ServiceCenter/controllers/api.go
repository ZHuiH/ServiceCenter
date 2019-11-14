package controllers

import(
	"ServiceCenter/common"
	"ServiceCenter/conf"
	"os"
	"io"
	//"fmt"
	"strings"
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
		data:=common.ReturnFormat("error",conf.MSG,"找不到该用户");
		json=common.JsonEncode(data);
	}

	ctx.WriteString(string(json));
}

//给用户上传图片的api
func Upload(ctx *http.RequestCtx){
	
	Core(ctx);
	//先检查 token
	token:=string(ctx.Request.Header.Peek("token"));
	if len(token) <=5 {
		ctx.WriteString(failJson("用户不存在"));
		return;
	}
	//接受 表单新的 图片
	fileHeader,err := ctx.FormFile("picture");
	//获取图片失败
	if err != nil{
        ctx.WriteString(failJson("图片不存在"))
        return;
	}

	//检查文件格式
	verfiy:=strings.Split(fileHeader.Header["Content-Type"][0],"/");
	//必须是image格式
	if len(verfiy) <=1 || verfiy[0] !="image"{
		ctx.WriteString(failJson("文件格式错误"));
		return;
	}

	file,err:=fileHeader.Open();
	//打开文件失败
	if err != nil{
        ctx.WriteString(failJson("打开文件失败"))
        return;
	}
	defer file.Close();

	//先查看一下文件夹是否存在
	current:=common.CurrentTime();
	path:=conf.Config.UploadPath+current;
	exist:=common.PathExist(path);
	if !exist{
		ctx.WriteString(failJson("创建目录失败"));
		return;
	}
	name:=common.RandStr(64)+"."+verfiy[1];
	//新开一个文件保存一下
	newFile:=path+"/"+name;

	nf,err := os.OpenFile(newFile,os.O_APPEND|os.O_CREATE|os.O_RDWR,0666)
    if err != nil{
        ctx.WriteString(failJson("上传文件失败"))
        return
	}
	defer nf.Close();
	//复制文件下去
	_,err = io.Copy(nf,file)
    if err != nil{
        ctx.WriteString(failJson("保存文件失败"))
        return;
	}
	picture_path:="/static/"+current+"/"+name;
    ctx.WriteString(successJson("上传文件成功",picture_path));

}

//输出的静态文件
func Static(ctx *http.RequestCtx){
	//请求的目录 截取路由的
	temp:=string(ctx.RequestURI()[7:]);

	//防止文件带参数 把参数去掉
	url:=strings.Split(temp,"?");

	path:="./upload/"+url[0];
	
	//检测文件是否存在
	fileInfo,err:=os.Stat(path);
	//没有文件 404
	if err!=nil{
		ctx.NotFound();
		return;
	}

	file,err:=os.Open(path);
	//打开文件失败
	if err !=nil{
		ctx.WriteString("打开文件失败");
		return;
	}

	defer file.Close();

	data:=make([]byte,fileInfo.Size());
	_,err=file.Read(data);

	//读取文件失败
	if err !=nil{
		ctx.WriteString("读取文件失败");
		return;
	}
	//输出
	ctx.Write(data);
}

//允许跨域
func Core(ctx *http.RequestCtx){
	ctx.Response.Header.Set("Access-Control-Allow-Origin","*");
	ctx.Response.Header.Set("Access-Control-Allow-Headers","Origin,Content-Type,Accept,token,x-requested-with");
	ctx.Response.Header.Set("Access-Control-Allow-Methods","POST,GET,OPTIONS,PUT,DELETE");
	ctx.Response.Header.Set("Content-Type","application/json; charset=utf-8;");
}

//返回失败的信息 json
func failJson(msg string) string {
	restult:=common.ReturnFormat("error",conf.TIPS,msg);
	return string(common.JsonEncode(restult));
}
//返回成功的信息 json
func successJson(msg string,data string) string {
	restult:=common.ReturnFormat("success",conf.TIPS,msg,data);
	return string(common.JsonEncode(restult));
}