import * as React from 'react';
import store from "../../store/Index"
import Content from "./Content"
import  "./dialogue.css"
import {UserParam} from "../../store/Active"

interface View {
    user    : UserParam | null
}

class Dialogue extends React.Component<{},View>{

    constructor(props:{}) {
        super(props);
        this.state={user:null};
        this.changeUser=this.changeUser.bind(this);
        this.MessageList=this.MessageList.bind(this);
        store.subscribe(this.changeUser);
    }
    
    private MessageList():JSX.Element {
        const user=this.state.user;
        let list=null;
        if(user && user.message.length>0){
            list=user.message.map((item,index)=>
                <Content key={`dialogue${index}`} name={user.nickName} avatar={user.avatar} message={item}  />
            )
        }

        return <div>{list}</div>
    }

    private changeUser(){
        const data=store.getState();
        if(data.target){
            this.setState({user:data.target},()=>{
                const target=document.getElementsByClassName("home-user-message")[0];
                target.scrollTop=target.scrollHeight;
            })
            
        }
    }

    public render(){
        return( 
            <div className="dialogue">
                {this.MessageList()}
            </div>
        );
    }
}   

export default Dialogue;