package file

import(
	"os"
	"time"
	"fmt"
)

//写入文件
//这里只支持 string的写入 
//模式也就2种 写入日志与普通写入 日志模式下 自动\n写入时间
func WriteFile(path string,content string,model bool) error {
	file,err:=os.OpenFile(path,os.O_CREATE|os.O_APPEND|os.O_RDWR,0766);
	if err !=nil{
		return err;
	}
	//自动关闭文件
	defer file.Close();

	var str string;
	//日志模式
	if model{
		str+=fmt.Sprintf("\n%s\t",time.Now());
	}

	str+=content;
	file.WriteString(str);
	return nil;
}