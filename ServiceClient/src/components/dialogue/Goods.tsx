import * as React from 'react';
import {Row,Col} from "antd"

interface View {
    name    :   string
    id      :   string
    code    :   string
    price   :   number
    picture :   string
}

class Goods extends React.Component<View,{}>{
    constructor(props:View) {
        super(props)
    }

    public render(){
        return( 
            <div  className="dialogue-goods">
                <p className="dialogue-goods-tips">分享商品</p>
                <Row>
                    <Col span={8}><img src={this.props.picture} className="dialogue-goods-picture"/></Col>
                    <Col span={16}>
                        <div className="dialogue-goods-info">
                            <b className="dialogue-goods-name">{this.props.name}</b>
                            <div>货号:  {this.props.code}</div>
                            <div>售价:  {this.props.price}</div>
                        </div>
                    </Col>
                </Row>
            </div>
        );
    }
}   

export default Goods;