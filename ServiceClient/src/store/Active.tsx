//基本上是后端传回来的信息模版
export interface MsgList {
    source      : boolean
    content     : string
    createTime  : number
}

//用户的数据
export interface UserParam{
    avatar       :   string //头像
    nickName     :   string //昵称
    message      :   MsgList[] //信息
    id           :   string //用户id
    token        :   string //用户id
}
//当前登录的客服自己的信息
export interface SelfParam {
    avatar       :   string //头像
    nickName     :   string //昵称
    token        :   string //用户token
    id           :   string //用户id
}
//当前推荐的商品
export interface SuggestGoodsParam {
    picture      :   string 
    price        :   number 
    name         :   string 
    id           :   string 
    code         :   string 
}
//真正需要共享的数据仓库
export interface Store {
    list                :   UserParam[]
    target              :   UserParam
    self                :   SelfParam
    suggestGoods        :   SuggestGoodsParam
    toView              :   object //token=>content
}

//仓库数据 变换定义
export enum TYPES{
    LOGIN   =   "LOGIN",
    ADD     =   "ADD",
    DELETE  =   "DELETE",
    SEND    =   "SEND",
    GET     =   "GET",
    RESET   =   "RESET",
    GOODS   =   "GOODS",
}

export interface UpdateStore{
    type        :   TYPES
    token       :   string
    content     :   any
}
