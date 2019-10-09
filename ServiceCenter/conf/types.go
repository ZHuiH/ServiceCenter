package conf
import(
	ws "github.com/fasthttp/websocket"
	http "github.com/valyala/fasthttp"
)

/*
*
	结构体的定义
*
*/

//程序的配置
type Configure struct{
	Env				bool	//查看配置
	Path			string	//配置文件的路径
	SessionLog		string	//会话日志的路径
	ErrorLog		string	//错误日志的路径
	AccessLog		string	//连接日志的路径
	Port  			string	//端口号
	GuestsName  	string	//游客的昵称
	Address			string	//ip地址
	ExpiresTime		int		//设置用户的会话过期时间 单位：秒
	PushUrl			string	//推送对话信息的地址  
	PushInterval	int		//推送间隔(再推送之前会临时保存信息，但是在推送之后会把信息删除掉)	设定这个时间之后 每隔这个时间 自动推送到设置的地址 单位：秒
	VerifyUrl		string	//验证登录的地址
	ServeVerify		string	//客服登录验证的地址
	Heartbeat		int		//心跳
	ServeUnion		bool	//是否开启联合服务模式 开启之后任何客服都能看到所有客户的发言 但是其他人回复依然是绑定的客服的名下
	ServeConcurrent int		//单个客服的并发数量;
}

//通用的信息
type UserInfo struct{
	CreateTime 		int64 		"生成时间"
	ActiveTime		int64		"该用户上次的活跃时间 根据这个时间去关闭连接"
	Id 				string		"用户的id"
	NickName		string		"用户的昵称"
	Connect			*ws.Conn	"连接的标识符"
	ProfilePicture	string		"头像"
	Close			bool		"是否已关闭连接并推送"
	Token			string		"验证使用的token 尽量所有验证都是用token 每次登录随机生成"
}

//信息列表
//客户跟客服都是指向同一个sclice
type Message struct{
	CreateTime	int64
	Content		string
	Source		bool	//通过这个来判断是不是客服发送的 true就是客服发起的
	//Customer	*Client
	//Service		*Serve
}

//客户的信息
type Client struct{
	Info		UserInfo
	Serve		*Serve
	History		*[]Message
	Next		*Client
}

//客服的信息
type Serve struct{
	Info				UserInfo
	ConnectNumber		int			"连接数量，随着客户的增加减少变化 比较重要根据这个数量来平均分配对接的客户数量"
	//History				*[]Message	"客服是没有自己的消息队列的"
	AnswerNumber		int			"回答的数量"
	//FastestAnswerTime	int			"最快的回答时间"
	//SlowestAnswerTime	int			"最慢的回答时间 单位是毫秒 因为会设置成5分钟不回复就断开 所以最大的时间不会超过5分钟"
	Next				*Serve		"下一个客服 按照登录的先后循序来排"
}

//客服列表的信息
type ServeList struct{
	Head 	*Serve		"第一个客服 随着首个客服的推出而变化"
	Total 	int			"客服总数 随着客服的增加减少变化"
}

//客户列表的信息
type ClientList struct{
	Head 	*Client		"第一个客户 随着首个客户的推出而变化"
	Total 	int			"客户总数 随着客户的增加减少变化"
}

//聊天记录 用于推送或者保存到本地
type  ChatHistory struct{
	UserId 		string
	UserName	string
	ServeId		string
	ServeName	string
	CreateTime	int64
	Content		[]Message
}

/*
*
	下面开始声明这些结构体
*
*/

//websocket 升级配置 允许跨域
var ChatConf = ws.FastHTTPUpgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:func(ctx *http.RequestCtx) bool {
		return true;
	},
}

//用户列表
var UserList *ClientList=new(ClientList);

//客服列表
var ServiceList *ServeList=new(ServeList);

//配置数据
var Config *Configure=new(Configure);

//游客数量
var GuestsTotal int;

