package conf

import(
	"time"
	"fmt"
	"strconv"
	"ServiceCenter/file"
	"ServiceCenter/common"
	"ServiceCenter/request"
)

//遍历用户列表 并使用传入的函数
func (list ClientList)ForEach(handler func(*Client)){
	item:=list.Head;
	for{

		if item==nil{
			break;
		}
		
		handler(item);
		item=item.Next;
	}
}

// 到最底部
//为了兼容  可以用其他位置继续往下找 不会用 UserList
func (head *Client)ToTail()(result *Client){
	result=head;
	for{
		if result==nil{
			break;
		}
		result=result.Next;
	}
	return;
}

//查找插入用户
func (list ClientList)AppendUser(user *Client)  (exists bool) {

	if list.Head==nil{
		list.Head=user;
		exists=true;
		list.Total++;
		return;
	}

	//先找有没有这个用户
	var target *Client;
	list.ForEach(func(item *Client){
		if user.Info.Id == item.Info.Id{
			if user.Info.Token != item.Info.Token{
				item.Info.Token=user.Info.Token;
			}
			fmt.Println(item.Info)
			target=item;
			return;
		};
	})

	//如果没有找到就可以注册了
	if target == nil{
		target=user;
		list.Total++;
		exists=true;
	}
	return;
}

//处理聊天历史的数据
func  (target *Client)HandlerChatHistory() ChatHistory {
	serve:=make(map[string]string);
	if target.Serve !=nil{
		serve["id"]=target.Serve.Info.Id;
		serve["name"]=target.Serve.Info.NickName;
	}else{
		serve["id"]="0";
		serve["name"]="未分配到客服";
	}
	//转换
	data:=ChatHistory{
		UserId 		:	target.Info.Id,
		UserName	:	target.Info.NickName,
		ServeId		:	serve["id"],
		ServeName	:	serve["name"],
		CreateTime	:	target.Info.CreateTime,
		Content		:	*target.History,
	}
	return data;
}

//推送信息
func (target *Client)PushMessage() bool {

	if len(Config.PushUrl) > 0{
		//转换数据
		data:=target.HandlerChatHistory();

		json:=common.JsonEncode(data);
		result:=string(json);
		param:=common.ReturnFormat("success","tips","用户超时",result);

		msg:=request.Post(Config.PushUrl,param);

		
		if msg["status"]=="error"{
			file.WriteFile(Config.ErrorLog,result,true);
		}else{
			return true;
		}
	}
	return false;
}

//自动检查 是否有脏数据 并删掉
func (list *ClientList)ClearUp(){

	for{

		time.Sleep(2 * time.Second);

		if list.Head == nil {
			time.Sleep(2 * time.Second);
			continue;
		}

		now:=time.Now().Unix();
		//需要删除的列表 装有下标的数组
		var deleteList []int;
		//记录当前循环的下标
		var index int;
		//先查询有没有超时的用户
		list.ForEach(func(item *Client){
			val:= now - item.Info.ActiveTime;
			str:=strconv.FormatInt(val,10);
			value,_:=strconv.Atoi(str);
			if value >= Config.ExpiresTime{
				deleteList=append(deleteList,index);
			}
			//下标进一
			index++;
		})

		if len(deleteList) > 0{
			//删除超时的用户
			list.Reset(deleteList);
		}else{
			time.Sleep(2 * time.Second);
		}

	}
}
//根据一组(有序)下标删除用户 并排序
func (list *ClientList)Reset(deleteList []int){
	count:=len(deleteList);

	if count <=0{
		return;
	}

	//下标
	var index int=0;
	//偏移量 防止删除了之后 下标不对
	//var offset int=0;
	//执行删除的时候 保存上一个 以免断开了
	var pre 	*Client;

	// if deleteList[0]==0{
	// 	target:=list.Head;
	// 	if list.Head.Next !=nil{
	// 		list.Head=list.Head.Next;
	// 	}else{
	// 		list.Head=&Client{};
	// 	}
	// 	target.PushMessage();
	// 	target.Recall();
	// 	deleteList=deleteList[1:];
	// 	count--;
	// 	list.Total--;
	// 	//所有下标减一
	// 	offset--;
	// }
	
	// if count <=0{
	// 	return;
	// }

	list.ForEach(func(item *Client){
		//防止超标
		if count<=0{
			return;
		}
		
		//加上偏移量
		key:=deleteList[0];
		//不能让链表失联
		if list.Head.Info.Id == item.Info.Id &&  index ==key{
			list.Head=item.Next;
		}

		if index ==key{
			if pre !=nil{
				//切断
				pre.Next=item.Next;         
			}
			//删除操作
			item.Recall();

			//更新偏移
			deleteList=deleteList[1:];
			count--;
			//offset++;
			list.Total--;
		}else{
			//保存上一位
			pre=item;
		}
		index++;
	})
}

//关闭指定的用户连接 注意这个并没有拼接删除之前的上下关系
func (target *Client)Recall(){
	//临死前的遗言
	msg:=common.ReturnFormat("success","tips","您长时间未提问，系统已自动断开");
	target.Info.Connect.WriteJSON(msg);

	//设置关闭标识
	target.Info.Close=true;

	return;
}
