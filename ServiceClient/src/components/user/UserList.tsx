import * as React from 'react';
import './user.css';
import User from "./User"
import {List,Input} from "antd"
import {UserParam} from "../../store/Active"
import store from "../../store/Index"
import {appendUser,reset} from "../../store/Reducers"

interface CallBack {
    changeUser  ?: (user:UserParam)=>void
}

interface UserListParam{
    list    :   UserParam[]
    view    :   boolean
    tipsList:   object
}


class UserList extends React.Component<CallBack,UserListParam>{

    constructor(props:CallBack) {
        super(props);
        
        this.buildFakeData=this.buildFakeData.bind(this);
        this.userList=this.userList.bind(this);
        this.messageTips=this.messageTips.bind(this);
        let data=this.buildFakeData();
        this.state={list:data,view:false,tipsList:{}};

        //专门用来登记用户发来的信息 显示提示的红点 
        store.subscribe(this.messageTips)
    }

    private messageTips(){
        this.setState({
            list:store.getState().list,
            tipsList:store.getState().toView
        })
    }

    private userList(item:UserParam){
        //更换用户
        const changeUser:()=>void=()=>{
            if(this.props.changeUser){
                this.props.changeUser(item);
                store.dispatch(reset(item))
            }
        }
        let target=this.state.tipsList[item.token];
        const tips=target ?  target.length : 0;

        return(
            <List.Item key={`user:${item.id}`}  onClick={changeUser} className="component-user-hover">
                <User 
                    token={item.token} 
                    tips={tips}
                    nickName={item.nickName} 
                    message={item.message.length > 0 ? item.message[item.message.length-1].content : ""} 
                    avatar={item.avatar}
                    id={item.id} />
            </List.Item>
        )
    }

    private buildFakeData():UserParam[] {
        let result:UserParam[]=[];
        for(let i=0;i<20;i++){
            let temp={
                avatar:"./logo.png",
                id:i.toString(),
                nickName:"user"+i.toString(),
                token:"token"+i.toString(),
                message:[{
                        content:`message for ${i.toString()}`,
                        source:true,
                        createTime:154156156,
                    }]
                };

            store.dispatch(appendUser(temp))
            result.push(temp)
        }
        return result;
    }

    public render(){
        return( 
            <div>
                <div className="home-user-list-search">
                    <Input.Search  placeholder="搜索" />
                </div>
                <List dataSource={this.state.list}  renderItem={this.userList} className="component-user-list"/>
            </div>
        );
    }
}   

export default UserList;