import * as React  from 'react';
import FontSize from "./option/FontSize";
import Hyperlink from "./option/Hyperlink";
import GoodsList from "./option/GoodsList";
import Facebread from "./option/Facebread";
import UploadPicture from "./option/UploadPicture";
interface MenuData {
    title       :   string
    action      ?:   string
    view        ?:   boolean
    icon        ?:   string
    value       ?:   string
    components  :   JSX.Element | null
}

interface MenuEvents {
    clickMenu       : (type:string)=>void
    confirm         : (action:string,viewUI:boolean,value:string)=>void
}

function MenuEvent(data:MenuEvents):MenuData[]{
    const envent=data;
    return [{
        title       :   "表情",
        components  :   <Facebread confirm={envent.confirm} />
    },{
        title       :   "加粗",
        action      :   "Bold",
        view        :   false,
        value       :   "",
        icon        :   "bold",
        components  :   null
    },{
        title       :   "斜体",
        action      :   "Italic",
        view        :   false,
        value       :   "",
        icon        :   "italic",
        components  :   null
    },{
        title       :   "下划线",
        action      :   "Underline",
        view        :   false,
        value       :   "",
        icon        :   "underline",
        components  :   null
    },{
        title       :  "删除线",
        action      :   "StrikeThrough",
        view        :   false,
        value       :   "",
        icon        :   "strikethrough",
        components  :    null
    },{
        title       :   "字体大小",
        components  :   <FontSize confirm={envent.confirm}  />
    },{
        title       :   "图片",
        components  :   <UploadPicture />
    },{
        title       :   "超链接",
        components  :   <Hyperlink confirm={envent.confirm}/>
    },{
        title       :   "商品列表",
        components  :  <GoodsList confirm={envent.confirm}/>
    }];
}
export default MenuEvent;