import * as React from "react"
import {Drawer,Tabs,Input } from "antd";
import store from "@/store/Index"
import {suggestGoods} from "@/store/Reducers"
import Goods from "./goods/Goods"

interface GoodsParam {
    confirm     ?:  (action:string,viewUI:boolean,value:string)=>void
}
interface GoodsInfo {
    name    :   string
    code    :   string
    price   :   number
    picture :   string
    id      :   string
    status  :   boolean
}
interface View {
    view        :   boolean
    list        :   GoodsInfo[]
    collection  :   GoodsInfo[]
}

class GoodsList extends React.Component<GoodsParam,View>{
    constructor(props:GoodsParam) {
        super(props)
        this.state={view:false,list:[],collection:[]};
        for(let i=0;i<20;i++){
            this.state.list.push({
                name    :   "测试商品名字",
                code    :   "code123456",
                price   :   999.99,
                id      :   (i*100).toString(),
                status  :   false,
                picture :   "https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=1834850066,1792480418&fm=26&gp=0.jpg",
            })
            this.migrate=this.migrate.bind(this);
        }
        this.goodsList=this.goodsList.bind(this)
        this.collectionList=this.collectionList.bind(this)
        this.viewLayer=this.viewLayer.bind(this);
        this.send=this.send.bind(this);
    }

    private migrate(id:string,index:number,status:boolean){
        let data=this.state.list;
        let collections=this.state.collection;
        if(!status){
            data[index].status=true;
            collections.push(data[index]);
        }else{
            data.forEach((item,key)=>{
                if(item.id === id){
                    data[key].status=false;
                    collections.splice(index,1);
                    return;
                }
            })
        }
        this.setState({list:data,collection:collections})
    }
    //返回商品的标签
    private goodsTab(item:GoodsInfo,index:number){
        return(
            <Goods 
            name={item.name} 
            code={item.code} 
            id={item.id} 
            index={index}
            picture={item.picture} 
            status={item.status} 
            price={item.price}  
            migrate={this.migrate}
            send={this.send}
            key={"goods"+index.toString()}/>
        );
    }
    //商品列表
    private goodsList():JSX.Element {
        const list=this.state.list.map((item,index)=>this.goodsTab(item,index));
        return <div className="goods-list">{list}</div>;
    };
    //收藏列表
    private collectionList(){
        const list=this.state.collection.map((item,index)=>this.goodsTab(item,index));
        return <div  className="goods-list">{list}</div>;
    }
    //发送商品
    private send(goods:any){
        let user=store.getState().target.token;
        if(user){
            //默认是发送完毕之后关闭列表
            this.viewLayer();
            store.dispatch(suggestGoods(user,goods));
        }
    }
    //组件的显示与隐藏
    private viewLayer():void{
        let val=!this.state.view;
        this.setState({view:val});
    }
    //抽屉在那里显示 父节点
    private main:() => HTMLElement=()=>document.getElementsByClassName("home")[0] as HTMLElement;
    public render(){
        return(
            <div className="menu-btn" >
                <span className="menu-icon icon-shangpin iconfont" onClick={this.viewLayer} style={{fontWeight:"bold"}} />
                <Drawer 
                    bodyStyle={{padding:"6px"}}
                    mask={false} 
                    maskClosable={false} 
                    title="商品列表"
                    width={300} 
                    getContainer={this.main}
                    onClose={this.viewLayer}
                    visible={this.state.view} 
                    style={{position:"absolute"}}
                >
                    <Input.Search   placeholder="目前只能搜索款号"  allowClear={true}  style={{height:"30px"}}/>
                    <Tabs defaultActiveKey="1">
                        <Tabs.TabPane key="goodsTab"  tab="商品列表">
                            {this.goodsList()}
                        </Tabs.TabPane>

                        <Tabs.TabPane key="collectionTab"  tab="收藏列表">
                            {this.collectionList()}
                        </Tabs.TabPane>
                    </Tabs>
                </Drawer> 
            </div>
            
        )
    }
}

export default GoodsList;