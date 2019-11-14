package conf

import(
	"time"
// 	"fmt"
// 	"strconv"
// 	"ServiceCenter/file"
// 	"ServiceCenter/common"
// 	"ServiceCenter/request"
)

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

	min:=list.Head;
	//遍历找出连接数最少的客服
	list.ForEach(func(item *Serve){
		//条件一：记录中最小连接数量大于当前的连接数量
		conditions1:= min.ConnectNumber > item.ConnectNumber;
		//条件2 当前的连接数量 小于等于 系统设置的最高连接数量
		conditions2:= item.ConnectNumber < Config.ServeConcurrent;
		//条件3 客服不处于关闭状态
		if conditions1 &&  conditions2 &&  !item.Info.Close {
			min=item;
		}
	})
	return min;
}

//有客服断开或者退出 重新排序
func (target *UserInfo)Sort() bool {
	if ServiceList.Head==nil{
		return true;
	}
	//防止客服断开了设置30秒作为重连时间
	time.Sleep(30)
	//重新连接了
	if !target.Close {
		return false;
	}

	//查看是否第一个
	if ServiceList.Head.Info.Token==target.Token{
		ServiceList.Head=ServiceList.Head.Next;
	}

	var pre *Serve=ServiceList.Head;
	ServiceList.ForEach(func(item *Serve){
		
		if item.Info.Token==target.Token{
			pre.Next=item.Next;
		}

		if item.Next==nil {
			return;
		}
		pre=item;
	});
	return true;
}