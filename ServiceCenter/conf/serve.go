package conf

// import(
// 	"time"
// 	"fmt"
// 	"strconv"
// 	"ServiceCenter/file"
// 	"ServiceCenter/common"
// 	"ServiceCenter/request"
// )

//遍历用户列表 并使用传入的函数
func (list *ServeList)ForEach(handler func(*Serve)){
	item:=list.Head;
	for{

		if item==nil{
			break;
		}
		
		handler(item);
		item=item.Next;
	}
}

//查找对接客户数最少的客服
func (list *ServeList)MinConnect() *Serve {

	if list.Head==nil{
		return nil;
	}

	min:=&Serve{};
	//遍历找出连接数最少的客服
	list.ForEach(func(item *Serve){
		if min.ConnectNumber > item.ConnectNumber && item.ConnectNumber < Config.ServeConcurrent{
			min=item;
		}
	})
	return min;
}