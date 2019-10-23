import * as React from "react"
import {Menu,Dropdown,Icon} from "antd";

interface FontSizeParam{
    confirm     ?:  (action:string,viewUI:boolean,value:string)=>void
}

class FontSize extends React.Component<FontSizeParam,{}>{
    constructor(props:FontSizeParam) {
        super(props)
        this.fontMenu=this.fontMenu.bind(this)
        this.changeFont=this.changeFont.bind(this);
    }

    //private fontList:string[]=["10px","11px","12px","14px","16px","18px","20px","24px","36px"];
    private fontList:string[]=["一号","二号","三号","四号","五号","六号","七号"];

    private fontMenu():JSX.Element {
        const list=this.fontList.map((item,index)=>
                <Menu.Item key={item}>
                    <a style={{fontSize:`${(index+1)*4}px`}}>{item}</a>
                </Menu.Item>
        );
        return <Menu onClick={this.changeFont}>{list}</Menu>;
    };
    private changeFont({key}:any){
        
        if(this.props.confirm){
            let size:number=this.fontList.indexOf(key);
            this.props.confirm("fontSize",false,size.toString());
        }
    };
    //确定回调
    // private sendLink(){
    //     //先执行
    //    // document.execCommand("CreateLink",false,this.Link);
    //     //成功之后的回调函数
    //     if(this.props.confirm){
    //         //this.props.confirm("");
    //     }
    // }
    //private dom:()=>HTMLElement=()=>document.getElementById("menu-FontSize") as HTMLElement;
    public render(){
        return(
            <div className="menu-btn" >
                <Dropdown overlay={this.fontMenu()}  placement="topCenter">
                    <Icon type="font-size"  className="menu-icon" />
                </Dropdown>
            </div>
            
        )
    }
}

export default FontSize;