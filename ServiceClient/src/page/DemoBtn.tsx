import * as React from 'react';
import store from '../store/Index';
import {UserParam} from '../store/Active';
import {Button,Menu,Dropdown } from "antd"
import {get} from '../store/Reducers';

interface View {
    list:UserParam[]
}

class DemoBtn extends React.Component<any,View>{
  constructor(props:any){
    super(props)
    this.state={
      list:store.getState().list,
    }
    this.userList=this.userList.bind(this)
    this.update=this.update.bind(this)
    this.item=this.item.bind(this)
    store.subscribe(this.update)
  }

  private update(){
    this.setState({
      list:store.getState().list
    })
  }
  private userList(){
    const list=this.state.list.map(item=>this.item(item))
    return <Menu>{list}</Menu>
  }


  private item(item:UserParam){
    const callBack=()=>{
      let data={
        source:true,
        content:"模拟用户发送信息",
        token:item.token,
        createTime :123456789
      }
      store.dispatch(get(data))
    }
    return <Menu.Item key={`demo${item.id}`} onClick={callBack} >{item.nickName}</Menu.Item>
  }

  public render() {
    return (
      <Dropdown overlay={this.userList} placement="topCenter">
        <Button size="small">模拟发送信息</Button>
      </Dropdown>
    );
  };
}

export default DemoBtn;