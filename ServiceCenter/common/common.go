package common
import(
	"encoding/json"
	"math/rand"
	"time"
	"os"
)


//随机字符串素材
var Letters = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

//json解密
func JsonDecode(content []byte) (result map[string]string ) {
	err:=json.Unmarshal(content,&result);
	if(err !=nil){
		return nil;
	}
	return;
}

//json加密
func JsonEncode(content interface{}) (result []byte ) {
	result,err:=json.Marshal(content);
	if(err !=nil){
		return nil;
	}
	return;
}

//生成随机字符串
func RandStr(length int) string {

	if length <=0{
		return "";
	}

	b:=make([]byte,length);
	for i:=range b{
		b[i]=Letters[rand.Intn(len(Letters))];
	}
	return string(b);
}

//设置返回格式
// 依次是 状态(必填) 信息(必填) 以及json的数据(选填)
func ReturnFormat(data ...string)(msg map[string]string){
	msg=make(map[string]string);

	if len(data) <3{
		return nil;
	}
	msg["status"]=data[0];
	msg["type"]=data[1];
	msg["content"]=data[2];
	
	if len(data) > 3{
		msg["data"]=data[3];
	}
	
	return;
}

//获取当前时间的年月日
func CurrentTime() string {
	now:=time.Now();
	result:=now.Format("20060102");
	return string(result);
}

//检测是否存在 如果不存在就创建
func PathExist(path string) bool {
	if _,err:= os.Stat(path); err != nil {
		err:=os.MkdirAll(path,os.ModePerm);
		if err!=nil{
			return false;
		}
	}
	return true;
}