package request
import(
	http "github.com/valyala/fasthttp"
	"ServiceCenter/common"
)

//请求操作
func Request(url string,method string,args string)(string,error){
	//请求
	req := http.AcquireRequest();
	//释放资源
	defer http.ReleaseRequest(req) 

	req.SetRequestURI(url)
	//统一设置成表单 不然 post无法接收 必须
	req.Header.SetContentType("application/x-www-form-urlencoded")
	req.Header.SetMethod(method);

	//添加参数
	if len(args) > 0{
		//格式：  parm=value
		req.AppendBodyString(args)
	}
	
	resp := http.AcquireResponse()
	//释放资源
    defer http.ReleaseResponse(resp) 
	client := &http.Client{}

    if err := client.Do(req, resp);err != nil {
		return "请求失败",err;
    }

	result := resp.Body()

	return string(result),nil;
}

//get请求
func Get(url string)(result map[string]string){
	//设置客户端
	data,err:=Request(url,"GET","");
	result=HandlerleResponse(data,err);
	return;
}

//post请求
func Post(url string,args map[string]string)(result map[string]string){
	var data string;
	//设置参数
	for key,item:= range args{
		data+=key+"="+item+"&";
	}
	//截取
	data=data[:len(data)-1];

	data,err:=Request(url,"POST",string(data));
	result=HandlerleResponse(data,err);
	return;
}

//处理http请求回来的数据
func HandlerleResponse(data string,err error)(msg map[string]string){
	if err !=nil{
		message:="请求失败:"+err.Error();
		msg=common.ReturnFormat("error","msg",message);
		return;
	}

	if  len(data) <=0{
		msg=common.ReturnFormat("success","msg","请求成功，无返回结果");
	}else{
		msg=common.JsonDecode([]byte(data));
		if len(msg) <=0{
			msg=common.ReturnFormat("error","msg","json解析失败");
		}
	}

	return;
}


// //处理http请求回来的数据
// func HandlerleResponse(code int,data []byte)(result map[string]string){
// 	msg:=HttpStatusCode(code);
// 	result=make(map[string]string)
// 	if msg["status"]=="success"{
// 		if len(data) > 0{
// 			result=JsonDecode(data);
// 		}else{
// 			result=msg;
// 			result["msg"]="请求成功，无返回数据";
// 		}

// 		if len(result) <= 0 {
// 			result=make(map[string]string);
// 			result["status"]="error";
// 			result["msg"]="json解析错误";
// 		}

// 	}else{
// 		result=msg;
// 	}
// 	return;
// }

// //处理http状态码
// func HttpStatusCode(code int)(result map[string]string){
// 	result=make(map[string]string)
// 	//默认是服务器请求 5XX
// 	result["status"]="error";
// 	result["msg"]="服务器请求失败！";

// 	if code ==200{
// 		result["status"]="success";
// 		result["msg"]="请求成功";
// 	}else if code < 300{
// 		result["status"]="success";
// 		result["msg"]="请求成功但无返回";
// 	}else if code < 400{
// 		result["msg"]="服务器请求被重定向！";
// 	}else if code < 500{
// 		result["msg"]="服务器请求错误！";
// 	}

// 	return;
// }