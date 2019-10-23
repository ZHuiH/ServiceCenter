import * as React from "react"
import {Icon,Popover,Row,Col} from "antd";
import Monochrome from "./FacebreadList/Monochrome";

interface FacebreadParam{
    confirm     ?:  (action:string,view:boolean,value:string)=>void
    close       ?:  ()=>void
    clickMenu   ?:  (type:string)=>void
}

class Facebread extends React.Component<FacebreadParam,{}>{
    constructor(props:FacebreadParam) {
        super(props)
        this.state={view:false}
        this.toInit();
        this.insertFacebread=this.insertFacebread.bind(this);
    }
    
    private FacebreadList:JSX.Element;
    private toInit(){
        //const handler:()=>void=
        const list=Monochrome.map((item,index)=>
            <Col span={4} key={index.toString()}>
                <p className="facebread-site" onClick={this.insertFacebread.bind(this,item.icon)}>
                    <span className={["iconfont","facebread-icon",item.icon].join(" ")}/>
                </p>
                <p className="facebread-name">{item.name}</p>
            </Col>
        );
        this.FacebreadList=<div className="facebread"><Row type="flex">{list}</Row></div>
    };
    //原本是使用图片的 目前先使用icon代替下
    //插入表情
    private insertFacebread(icon:string){
        if(this.props.confirm){
            let val:string=`<span class="iconfont ${icon} facebread-icon"></span>`;
            //let img:string=`<img src="http://img.baidu.com/hi/jx2/j_0021.gif">`;
            this.props.confirm("insertHTML",false,val)
        }
    }
    public render(){
        return(
            <div className="menu-btn" >
                <Popover placement="topLeft" title="表情" content={this.FacebreadList}>
                    <Icon type="frown" className="menu-icon" />
                </Popover>
            </div>
        )
    }
}

export default Facebread;