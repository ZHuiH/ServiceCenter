import * as React from 'react';
//import axios from 'axios';
//import Api from '../static/api/Api';
import '../static/css/home.css';
import { RouteComponentProps } from "react-router-dom";
import {Layout} from "antd"
import UserList from '../components/user/UserList';
import User from '../components/user/User';
import {UserParam} from '../store/Active';
import Edit from '../components/edit/Edit';
import Dialogue from '../components/dialogue/Index';


const { Header, Footer, Sider, Content } = Layout;

interface View{
  user:UserParam | null
}


class Home extends React.Component<RouteComponentProps,View>{
  constructor(props:RouteComponentProps){
    super(props);
    this.changeUser=this.changeUser.bind(this)

    this.state={user:null};
    this.currentCustomer=this.currentCustomer.bind(this)

  } 

  private currentCustomer() :JSX.Element {
    const user=this.state.user;
    if(user){
      let len=user.message.length;
      let msg=len > 0 ? user.message[len-1].content : "";
      return(
        <User token={user.token} nickName={user.nickName} message={msg} avatar={user.avatar} id={user.id} />
      );
    }else{
      return  <p style={{color:"#333"}}>当前并无客户对接</p>;
    }
  };

  private changeUser(target:UserParam){
    this.setState({user:target})
  }
  
  public render() {
    return (
        <div className="home">
          <Layout>
            {/* 左边用户列表 */}
            <Sider className="home-user-list">
                <UserList changeUser={this.changeUser} />
            </Sider>
            {/* 头部 当前的用户信息 */}
            <Layout>
              <Header className="home-header">
                <this.currentCustomer/>
              </Header>

              {/* 中间的对话内容 */}
              <Content className="home-user-message">
                <Dialogue />
              </Content>

              {/* 底部编辑 */}
              <Footer className="home-edit">
                  <Edit userToken={this.state.user ? this.state.user.token : ""}/>
              </Footer>
            </Layout>
          </Layout>
        </div>


    );
  };
}



export default Home;