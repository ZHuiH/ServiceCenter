package file

import(
	"io"
	"os"
	"strings"
	"bufio"
)

//逐行读取文件
func ReadLine(path string,handler func(string)) error {
	file,err:=os.Open(path);
	//自动关闭文件
	defer file.Close();

	if err !=nil{
		return err;
	}
	buf:=bufio.NewReader(file);
	for{
		line,err:=buf.ReadString('\n');
		line=strings.TrimSpace(line);

		//err里面包含了文件结束符与错误
		if err != nil{
			//查看是否文件尾部
			if err !=io.EOF{
				return err;
			}
			return nil;
		}
		handler(line);
	}
}