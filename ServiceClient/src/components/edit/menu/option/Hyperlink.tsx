import * as React from "react"
import {Modal,Input,Icon,Tooltip} from "antd";

interface HyperlinkParam{
    confirm     ?:  (action:string,view:boolean,value:string)=>void
    close       ?:  ()=>void
    clickMenu   ?:  (type:string)=>void
}

interface View {
    view        :   boolean
}

class Hyperlink extends React.Component<HyperlinkParam,View>{
    constructor(props:HyperlinkParam) {
        super(props)
        this.state={view:false}
        this.saveLink=this.saveLink.bind(this)
        this.viewLayer=this.viewLayer.bind(this)
        this.sendLink=this.sendLink.bind(this)
    }
    //填写url
    private Link:string;
    //每次填写都保存
    private saveLink(elem:React.ChangeEvent<HTMLInputElement>):void{
        this.Link=elem.target.value
        //window.getSelection().anchorNode.parentNode
    }
    //确定回调
    private sendLink(){
        this.viewLayer();
        //成功之后的回调函数
        if(this.props.confirm){
            this.props.confirm("CreateLink",false,this.Link);
        }
    }
    //组件的显示与隐藏
    private viewLayer():void{
        let val=!this.state.view;
        this.setState({view:val});
    }
    public render(){
        return(
            <div className="menu-btn" >
                <Tooltip placement="bottom" title="超链接">
                    <Icon type="link" onClick={this.viewLayer} className="menu-icon" />
                    <Modal title="创建超链接" visible={this.state.view} onCancel={this.viewLayer} onOk={this.sendLink}  mask={false} maskClosable={false}>
                        <Input placeholder="请填写连接地址" onChange={this.saveLink}  />
                    </Modal>
                </Tooltip>
            </div>
        )
    }
}

export default Hyperlink;