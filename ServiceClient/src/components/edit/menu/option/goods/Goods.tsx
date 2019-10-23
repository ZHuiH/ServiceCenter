import * as React from "react"
import {Col,Row,Button} from "antd";
interface GoodsInfo {
    name    :   string
    id      :   string
    code    :   string
    price   :   number
    picture :   string
    status  :   boolean
    index   :   number
    migrate ?:  (id:string,index:number,status:boolean)=>void
    send    ?:  (data:any)=>void
}

class Goods extends React.Component<GoodsInfo,{}>{
    constructor(props:GoodsInfo) {
        super(props)
        this.migrate=this.migrate.bind(this);
        this.send=this.send.bind(this);
    }
    private migrate(){
        if(this.props.migrate){
            this.props.migrate(this.props.id,this.props.index,this.props.status);
        }
    }
    private send(){
        if(this.props.send){
            this.props.send({
                id:this.props.id,
                name:this.props.name,
                code:this.props.code,
                picture:this.props.picture,
                price:this.props.price,
            });
        }
    }

    public render(){
        return(
            <div className="goods-item">
                <Row>
                    <Col span={10}>
                        <div className="goods-item-img"><img src={this.props.picture} style={{width:"100%",height:"100%"}} /></div>
                    </Col>
                    <Col span={14}>
                        <div style={{width:"92%",margin:"auto"}}>
                            <b className="goods-item-info-title">{this.props.name}</b>
                            <p className="goods-item-info">货号：{this.props.code}</p>
                            <p className="goods-item-info">售价：{this.props.price}</p>
                        </div>
                    </Col>
                </Row>
                <div className="goods-item-btn-list">
                    <Button onClick={this.send} size="small" type="primary" ghost={true} className="goods-item-btn goods-item-btn-send">发送</Button>
                    <Button 
                        type={this.props.status ? "danger" : "primary"} 
                        onClick={this.migrate} size="small" 
                        ghost={true} className="goods-item-btn" >
                            {this.props.status ? "取消标记" : "标记"}
                    </Button>
                </div>
            </div>
        )
    }
}

export default Goods;