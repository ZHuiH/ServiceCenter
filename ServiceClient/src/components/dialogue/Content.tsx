import * as React from 'react';
import {MsgList,SelfParam} from "../../store/Active"
import {Avatar,Row,Col} from "antd"
import store from "../../store/Index"
import Goods from "./Goods"

interface View {
    message     :   MsgList
    avatar      :   string
    name        :   string
}

class Content extends React.Component<View,{}>{
    constructor(props:View) {
        super(props)
        this.discern=this.discern.bind(this);
        this.template=this.template.bind(this);
        this.sender=this.sender.bind(this);
        this.speaker=this.speaker.bind(this);
        this.content=this.content.bind(this);
        this.self=store.getState().self;
    }
    private self:SelfParam;
    //分辨出是自己发送的还是客户发送的
    private discern():JSX.Element {
        if(this.props.message.source){
            return this.sender()
        }
        return this.speaker()
    }

    private sender():JSX.Element{
        const avatar=this.props.avatar
        const name=this.props.name
        return(
            <Row type="flex">
                <Col span={2}>
                    <Avatar icon="user"  size={32} src={avatar}/>
                </Col>
                <Col span={22}>
                {this.template(name,"user")}
                </Col>
            </Row>
        );
    }

    private speaker(){
        const avatar=this.self.avatar
        const name=this.self.nickName
        return(
            <Row type="flex">
                <Col span={22}>
                    {this.template(name,"self")}
                </Col>
                <Col span={2}>
                    <Avatar icon="user"  size={32} src={avatar}/>
                </Col>
            </Row>
        );
    }

    private template(name:string,role:string){
        return(
            <div className={`dialogue-${role}-message`}>
            <div className="dialogue-name">{name}</div>
            {this.content(role)}
        </div>
        )
    }

    private content(role:string):JSX.Element{
        const content=this.props.message.content;
        let index=content.includes("suggest_goods");
        let className=`dialogue-content dialogue-${role}-hangings`;
        if(index){
            const goods=store.getState().suggestGoods;
            return <div className={className}><Goods name={goods.name} id={goods.id}  picture={goods.picture} price={goods.price} code={goods.code}/></div>;
        }else{
            return <div className={className} dangerouslySetInnerHTML={{__html:content}}/>;
        }
       
    }

    public render(){
        return( 
            <div style={{marginBottom:"10px"}}><this.discern/></div>
        );
    }
}   

export default Content;