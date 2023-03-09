import React from 'react';
import {Image} from "antd";

function Message(props) {
    const data = props.data;
    const bg = props.idx !== props.showIdx ? "#000" : "#ff0000"
    return (
        <div onClick={() => {
            props.setShowIdx(props.idx)
        }} style={{margin: "5px 10px"}}>
            <div style={{color: bg}}>{data.username}</div>
            <hr/>
            <Image src={data.url} preview={false}/>
        </div>
    );
}

export default Message;

export function List(props) {
    const data = props.data;
    return (
        <div>
            {data.map((item, idx) => {
                return (
                    <Message key={item.id} idx={idx} showIdx={props.showIdx} setShowIdx={props.setShowIdx} data={item}/>
                )
            })}
        </div>
    )
}