import{TYPES,Store,UpdateStore,UserParam,SelfParam,SuggestGoodsParam} from "./Active"
const {ADD,DELETE,SEND,RESET,LOGIN,GET,GOODS}=TYPES;

const defultTalk:Store={
    list:[],
    target:{
        avatar       :   "", 
        nickName     :   "",
        message      :   [],
        token        :   "",
        id        :   "",
    },
    self:{
        avatar       :   "", 
        nickName     :   "",
        token           :   "",
        id        :   "",
    },
    suggestGoods:{
        picture      :   "", 
        name         :   "", 
        price        :   0, 
        id           :   "", 
        code         :   "", 
    },
    toView           :   {}
};
//store 的数据操作
export function storeActive(store:Store=defultTalk,action?:any){
    //复制一个新的对象 不更改原store
    let NewStore=Object.assign({},store)
    if(action){
        switch(action.type){
            case LOGIN   : handleLogin(NewStore,action.content);break;
            case ADD     : handleAdd(NewStore,action.content);break;
            case RESET   : handleReset(NewStore,action.content);break;
            case DELETE  : break;
            case SEND    : handleSend(NewStore,action);break;
            case GOODS   : handleGoods(NewStore,action);break;
            case GET     : handleGet(NewStore,action.content);break;
        }
    }
    return NewStore;
}

//登录
function handleLogin(store:Store,self:SelfParam){
    //注意因为是保存的对象里面 没有做其他处理
    //刷新之后就没了登录状态
    //由于现在是demo 所以不做任何多余的操作
    store.self=self;
}

//添加用户
function handleAdd(store:Store,content:UserParam){
    store.list.push(content)
}

//现在目标的用户
function handleReset(store:Store,content:UserParam){
    store.target=content;
}

//现在需要发送的商品
function handleGoods(store:Store,action:any){
    store.suggestGoods=action.content;

    store.list.forEach(item=>{
        //找一下 到底发给谁
        if(item.token===action.token){
            item.message.push({
                content:`suggest_goods=${action.content.id}`,
                source:false,
                createTime:new Date().valueOf()/1000,
            })
        }
    })
}

//处理用户需要推送的信息
function handleSend(store:Store,action:any){
    store.list.forEach(item=>{
        //找一下 到底发给谁
        if(item.token===action.token){
            item.message.push({
                content:action.content,
                source:false,
                createTime:new Date().valueOf()/1000,
            })
        }
    })
}

//服务器推送过来的信息
function handleGet(store:Store,action:any){
    
    store.list.forEach(item=>{
        //找一下 到底发给谁
        if(item.token===action.token){
            item.message.push({
                content:action.content,
                source:action.source,
                createTime:action.createTime,
            })
            //是用户发的 而且不是当前的用户
            if(action.source && action.token !== store.target.token){
                if(store.toView[action.token]){
                    store.toView[action.token].push(action.content)
                }else{
                    store.toView[action.token]=[action.content]
                }
            }
        }
    })
}

//添加用户
export function login(contents:SelfParam):UpdateStore{
    let result:UpdateStore={
        type    :   LOGIN,
        token      :   "",
        content :   contents
    }
    return result;
}
//添加用户
export function appendUser(contents:UserParam):UpdateStore{
    let result:UpdateStore={
        type    :   ADD,
        token   :   "",
        content :   contents
    }
    return result;
}
//自己发送信息
export function send(targetToken:string,contents:string):UpdateStore{
    let result:UpdateStore={
        type    :   SEND,
        token   :   targetToken,
        content :   contents
    }
    return result;
}

//用户给你发信息
export function get(data:any):UpdateStore{
    let result:UpdateStore={
        type    :   GET,
        token   :   "",
        content :   data
    }
    return result;
}

//点击用户
export function reset(contents:UserParam):UpdateStore{
    let result:UpdateStore={
        type    :   RESET,
        token   :   "",
        content :   contents
    }
    return result;
}

//点击用户
export function suggestGoods(targetToken:string,goods:SuggestGoodsParam):UpdateStore{
    let result:UpdateStore={
        type    :   GOODS,
        token   :   targetToken,
        content :   goods
    }
    return result;
}

