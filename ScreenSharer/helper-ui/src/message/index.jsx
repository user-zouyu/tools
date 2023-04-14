import React from 'react';
import {Image} from "antd";

function Message(props) {
    const data = props.data;
    const bg = data.id !== props.showId ? "#000" : "#ff0000"

    return (
        <div onClick={() => {
            props.setShowId(data.id)
            if (props.ws !== null) {
                props.ws.send(JSON.stringify({
                    command: "setup",
                    id: data.id.toString(),
                }))
            }
        }} style={{margin: "5px 10px"}}>
            <div style={{color: bg}}>{data.id}:{data.username}</div>
            <hr/>
            {
                data.type === "text" ?
                    <div style={{maxWidth: "100%", wordBreak: "break-all"}}>
                        {JSON.parse(data.data).text.substring(0, 100)}
                    </div>
                    :
                    <Image src={data.data} preview={false}/>
            }
        </div>
    );
}

export default Message;