import * as React from "react"
import Menu from "./menu/Menu"
import store from "../../store/Index"
import {send} from "../../store/Reducers"
import {Button} from "antd";
import DemoBtn from "@/page/DemoBtn"

interface User {
    userToken: string | null
}

class Edit extends React.Component<User,{}>{
    constructor(porps:User){
        super(porps);
        this.execute=this.execute.bind(this)
        this.send=this.send.bind(this)
    }
    public componentDidMount(){
        let element=document.getElementsByTagName("iframe")[0];
        let doc=element.contentDocument;

        if(doc){
            let css=document.getElementsByTagName("style")[1].innerHTML;
            let style=doc.createElement("style");
            style.type = 'text/css';
            style.innerHTML=css;
            doc.head.appendChild(style)
            this.doc=doc;
            doc.designMode="on";
        }
    }
    //document 对象 任何菜单栏的操作都是对他进行操作
    private doc:Document;
    //执行操作都是在这里 从menu组件处理完的数据传回来
    private execute(action:string,view:boolean,value:string):void {
        this.doc.execCommand(action,view,value);
    }
    //发送信息
    private send(){
        const body=this.doc.body.innerHTML;
        if(this.props.userToken && body){
            //发送
            store.dispatch(send(this.props.userToken,body))
            //重置
            this.doc.body.innerHTML="";
        }

    }
    public render(){
        return(
            <div className="edit">
                <div className="edit-menu"><Menu clickMenu={this.execute} /></div>
                <iframe  id="edit-main" contentEditable={true} suppressContentEditableWarning={true}/>
                <div className="edit-send">
                    <Button.Group>
                        <Button onClick={this.send} size="small">发送</Button>
                        <Button icon="down" size="small"/>
                        <DemoBtn />
                    </Button.Group>
                </div>
            </div>
        )
    }
}

export default Edit;