import * as React from 'react';
import './user.css';
import {Avatar,Row,Col,Badge} from "antd"

interface UserParam {
    avatar       :   string //头像
    nickName     :   string //昵称
    message      :   string //信息
    id           :   string //用户id
    token        :   string //用户id
    tips         ?:  number
    selection    ?:  (data:object)=>void
};

class User extends React.Component<UserParam,object>{

    constructor(props:UserParam) {
        super(props);
        this.selection=this.selection.bind(this);
       // this.viewTips=this.viewTips.bind(this);
    }
    // private viewTips(){
    //     if(this.props.tips){
    //         return <Badge count={this.props.tips}/>
    //     }
    //     return <div/>
    // }
    public render(){
        return( 
            <div onClick={this.selection} className="component-user">
                <Row type="flex">
                    <Col span={7}><Avatar icon="user" src={this.props.avatar} size="large"/></Col>
                    <Col span={17}>
                        <Badge count={this.props.tips} style={{position:"absolute",top:"10px"}}>
                            <div className="component-user-info">
                                <span className="component-user-info-name">{this.props.nickName}</span>
                                    <br/>
                                <span className="component-user-view-msg">{this.props.message}</span>
                            </div>
                        </Badge>
                    </Col>
                </Row>
            </div>
        );
    }

    private selection():void{
        let data= {
                    id:this.props.id,
                    token:this.props.token
                }
        if(this.props.selection){
            this.props.selection(data);
        }
        
    }
}   

export default User;