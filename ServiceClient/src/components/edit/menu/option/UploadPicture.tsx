import * as React from "react"
import {Upload,Icon} from "antd";



class UploadPicture extends React.Component<any,{}>{
    constructor(props:any) {
        super(props)
    }

    public render(){
        return(
            <div className="menu-btn" >
                <Upload>
                    <Icon type="picture"  className="menu-touch" />
                </Upload>
            </div>
        )
    }
}

export default UploadPicture;