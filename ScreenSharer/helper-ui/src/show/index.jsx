import React from 'react';
import {Image} from "antd";
import SyntaxHighlighter from "react-syntax-highlighter";
import {vs2015} from "react-syntax-highlighter/src/styles/hljs/index.js";


const Show = (props) => {
    let item = {}
    console.log(props)
    if (props.id !== -1) {
        for (let i = 0; i < props.list.length; i++) {
            if (props.list[i].id === props.id) {
                item = props.list[i]
            }
        }
    }
    if (props.id < 0) {
        return (
            <div style={{
                display: "flex",
                marginTop: "30px",
                justifyContent: "space-around",
                alignItems: "center",
                fontSize: "30px"
            }}>
                欢迎使用
            </div>
        )
    }


    if (item.type === "text") {
        const data = JSON.parse(item.data)
        return (
            <div style={{
                width: "calc(100%)",
                minHeight: "100%",
                display: "inline-block",
                padding: "0 14px"
            }}>
                <SyntaxHighlighter language={data.language} style={vs2015}>
                    {data.text}
                </SyntaxHighlighter>
            </div>
        )
    }

    if (item.type === "image") {
        return (
            <div style={{
                display: "flex",
                justifyContent: "space-around",
                alignItems: "center"
            }}>
                <Image style={{maxHeight: "100vh"}} src={item.data}/>
            </div>
        )
    }

    return (
        <div style={{
            display: "flex",
            marginTop: "30px",
            justifyContent: "space-around",
            alignItems: "center",
            fontSize: "30px"
        }}>
            暂不支持该数据类型
        </div>
    )

}

export default Show;