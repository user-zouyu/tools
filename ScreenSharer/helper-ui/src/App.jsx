import React, {useState} from 'react'
import './App.css'
import {FloatButton, Form, Input, Layout, List, message, Modal, Switch} from "antd";
import Show from "./show/index.jsx";
import Message from "./message/index.jsx";
import {SettingOutlined} from "@ant-design/icons";


function App() {
    const [username, setUsername] = useState("zou yu")
    const [btnName, setBtnName] = useState("连接")
    const [list, setList] = useState([]);
    const [showIdx, setShowIdx] = useState(-1)
    const [ws, updateWs] = useState(null);
    const [connected, setConnected] = useState(false);
    const [host, setHost] = useState(window.location.host)
    const [siderDisplay, setSiderDisplay] = useState("block")

    const [modalShow, setModalShow] = useState(false)

    const connect = () => {
        if (username.length < 5) {
            message.error("用户名必须超过5个字符").then(() => {
            })
            return
        }
        if (!connected) {
            const ws = new WebSocket("ws://" + host + "/api/connect?username=" + username);
            ws.addEventListener('open', () => {
                message.info("连接成功").then(() => {
                })
                setConnected(true)
                updateWs(ws)
                setBtnName("断开连接")
            });

            ws.addEventListener("error", () => {
                message.error("连接出错啦").then(() => {
                })
            })
            ws.addEventListener('message', (event) => {
                const data = JSON.parse(event.data);
                console.log(data)
                message.info(data.msg).then(() => {
                })

                if (data.code === 1) {
                    setList(prevState => {
                        return [...prevState, ...data.data]
                    })

                }
                if (data.code === 2) {
                    setList(data.data)
                }

                if (data.code === 4) {
                    setShowIdx(pre => {
                        return pre + parseInt(data.data)
                    })
                }
            });

        } else {
            ws.close()
            message.info("连接关闭").then(() => {
            })
            updateWs(null)
            setConnected(false)
            setBtnName("连接")
        }
    }


    return (
        <div>
            <Modal title="Settings"
                   open={modalShow}
                   onOk={() => connect()}
                   onCancel={() => setModalShow(false)}
                   okText={btnName}
            >
                <Form labelCol={{
                    span: 8
                }} wrapperCol={{
                    span: 16
                }}>
                    <Form.Item
                        label="聊天记录"
                    >
                        <Switch
                            checkedChildren="开启"
                            unCheckedChildren="关闭"
                            defaultChecked
                            onChange={(e) => {
                                console.log(e)
                                e ? setSiderDisplay("block"): setSiderDisplay("none")
                            }}
                        />
                    </Form.Item>
                    <Form.Item
                        label="服务器地址"
                    >
                        <Input placeholder="服务器地址"
                               onChange={(e) => {
                                   setHost(e.target.value)
                               }}
                               value={host}
                               disabled={connected}/>
                    </Form.Item>
                    <Form.Item
                        label="用户名"
                    >
                        <Input placeholder="用户名"
                               onChange={(e) => {
                                   setUsername(e.target.value)
                               }}
                               value={username}
                               disabled={connected}/>
                    </Form.Item>
                </Form>
            </Modal>
            <FloatButton
                icon={<SettingOutlined />}
                onClick={() => {
                setModalShow(true);
            }}/>
            <Layout className="site-layout">
                <Layout.Sider style={{
                    display: siderDisplay,
                    padding: "5px",
                    borderLeft: "1px",
                    overflow: 'auto',
                    height: '100vh',
                    left: 0,
                    top: 0,
                    bottom: 0,
                }} theme="light"
                >
                    <List
                        header={<div>聊天记录</div>}
                        bordered
                        dataSource={list}
                        renderItem={(item, idx) => (
                            <List.Item>
                                <Message key={item.id} idx={idx} showIdx={showIdx} setShowIdx={setShowIdx} data={item}/>
                            </List.Item>
                        )}
                    />
                </Layout.Sider>

                <Layout.Content style={{height: "100vh", width: "100vw"}}>
                    <Show list={list} index={showIdx}/>
                </Layout.Content>
            </Layout>

        </div>
    )
}

export default App
