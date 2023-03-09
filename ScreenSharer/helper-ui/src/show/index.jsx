import React from 'react';
import {Image} from "antd";

const Show = (props) => {
    const index = props.index
    const data = props.list[index];
    if (index < 0) {
        return (
            <div style={{
                display: "flex",
                marginTop: "30px",
                justifyContent: "space-around",
                alignItems: "center",
                fontSize: "30px"
            }}>
                shift+s 切换到下一张照片
            </div>
        )
    }
    if (index >= props.list.length) {
        return (
            <div style={{
                display: "flex",
                marginTop: "30px",
                justifyContent: "space-around",
                alignItems: "center",
                fontSize: "30px"
            }}>
                shift+w 切换到上一张照片
            </div>
        )
    }

    return (
        <div style={{
            display: "flex",
            justifyContent: "space-around",
            alignItems: "center"
        }}>
            <Image style={{maxHeight: "100vh"}} src={data.url}/>
        </div>
    )
}

export default Show;