import * as React from 'react';
// import axios from 'axios';
// import Api from '../static/api/Api';
import '../static/css/login.css';
import store from "../store/Index"
import {login} from "../store/Reducers"
import { RouteComponentProps } from "react-router-dom";
import {Button,Row,Col,Input,Alert} from "antd"

interface UserInfo {
  account   : string,
  password  : string,
  status    : boolean,
  message   : string,
  type      : any,
  view      : boolean,
}

class Login extends React.Component<RouteComponentProps,UserInfo> {
  constructor(props:RouteComponentProps){
    super(props);
    
    this.state={
      account:"",
      password:"",
      status:false,
      type:"",
      message:"",
      view:false,
    }

    this.bindAccount=this.bindAccount.bind(this)
    this.bindPassWrod=this.bindPassWrod.bind(this)
    this.toLogin=this.toLogin.bind(this)
    this.closeAlert=this.closeAlert.bind(this)
  }
  public render() {
    return (
      <div className="login">
        <div className="login-main">
          <div className="login-logo">
            <img src="./logo.png"/>
          </div>
          <div className="login-info">
            
            <div className="login-input">
              <Row type="flex" align="middle">
                <Col span={5}>帐号</Col>
                <Col span={19}>
                  <Input placeholder="请输入帐号" allowClear={true} onChange={this.bindAccount} />
                </Col>
              </Row>
            </div>

            <div className="login-input">
              <Row type="flex" align="middle">
                <Col span={5}>密码</Col>
                <Col span={19}>
                  <Input placeholder="请输入帐号" allowClear={true} type="password" onChange={this.bindPassWrod}/>
                </Col>
              </Row>
            </div>

            <div className="login-btn">
              <Button type="primary" block={true} onClick={this.toLogin} loading={this.state.status}>确定</Button>
            </div>

          </div>
          <div className="login-alert" style={{display:this.state.view ? "block" : "none"}}>
            <Alert message={this.state.message} showIcon={true} type={this.state.type}  closable={true} onClose={this.closeAlert}/>
          </div>
        </div>
      </div>
    );
  };
  //帐号
  private bindAccount(elem: React.ChangeEvent<HTMLInputElement>) :void {
    this.setState({account:elem.target.value.trim()})
  };
  //密码
  private bindPassWrod(elem: React.ChangeEvent<HTMLInputElement>) :void {
    this.setState({password:elem.target.value.trim()})
  };
  //登录
  private toLogin(){
    let json=`{"status":"success","msg":"\u767b\u5f55\u6210\u529f","data":{"nickName":"( ^ω^)","id":"49709","token":"user74755","avatar":"./user.jpg"}}`;
    let result=JSON.parse(json);
    if(result.status === "success"){
      store.dispatch(login(result.data))
      this.props.history.push("/home");
    }

    // if(this.state.account.length <5 || this.state.password.length <5){
    //   this.setState({
    //     message:"帐号或密码格式错误",
    //     type:"error",
    //     view:true,
    //     status:false,
    //   })
    //   return;
    // }

    // this.setState({status:true,view:false});
    // axios.post(Api.login,this.state).then(res=>{
    //   let result=res.data;

    //   this.setState({
    //     message:result.msg,
    //     type:result.status,
    //     view:true,
    //     status:false,
    //   })

    //   if(result.status === "success"){
    //     this.props.history.push("/home");
    //   }
      
    //})
  };
  //关闭弹窗
  private closeAlert(){
    this.setState({view:false});
  }
}

export default Login;