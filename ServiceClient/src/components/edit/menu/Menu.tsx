import * as React from 'react';
import {Row,Col,Tooltip,Icon} from "antd";
import MenuData from "./menuData";
import  "../edit.css";

interface MenuEvent {
    data:Event
}

interface Event{
    confirm       :   (action:string,viewUI:boolean,value:string)=>void
    clickMenu       :   (type:string)=>void
}

interface CallBack {
    clickMenu   :   (action:string,viewUI:boolean,value:string)=>void
}

interface Commands {
    actions  :   string
    view    :   boolean
    value   :   string
}

class Menu extends React.Component<CallBack,MenuEvent>{
    constructor(props:CallBack){
        super(props)
        this.handlerExec=this.handlerExec.bind(this);
        this.state={
            data:{
                clickMenu:this.clickMenu,
                confirm:this.handlerExec
            }
        }
        this.getComponents=this.getComponents.bind(this);
        //执行命令
        this.execute=this.execute.bind(this);
    };
    private MenuList() :JSX.Element{
        const list =MenuData(this.state.data).map((item,index)=>
            <Col key={index.toString()} span={2} className="menu-touch" id={`menu-${item.action}`} >
                <Tooltip placement="bottom" title={item.title} trigger="hover">
                    {this.getComponents(item)}
                </Tooltip>
            </Col>
        )
        return <Row type="flex">{list}</Row>
    };
    private getComponents(item:any){
        if(item.components){
            //使用自定义的组件 不会马上执行命令
            return item.components;
        }else{
            //如果是null那就是使用默认 直接使用命令
            const handler:()=>void=()=>this.handlerExec(item.action,item.view,item.value);
            return <Icon type={item.icon} onClick={handler} className="menu-icon" />
        }
    }
    private clickMenu(type:string){
        //先过滤掉复杂的步骤 简单直接执行
        switch(type){
            //case "CreateLink":this.viewLink();break;
            // case "viewGoods":this.viewGoods();break;
            // case "FontSize":this.viewFont();break;
            //default:this.handlerExec(item.action,item.viewUI,item.value);
        }
    }

    //让组件调用 不需要再重新定义太多东西
    private handlerExec(action:string,show:boolean,val:string){
        let result:Commands={
            actions:action,
            view:show,
            value:val,
        }
        this.execute(result);
    }
    //点击回调 进行操作
    private execute(data:Commands){
        this.props.clickMenu(data.actions,data.view,data.value)
    }
    public render(){
        return(
            <div>
                {this.MenuList()}
            </div>
        )
    }
}

export default Menu;