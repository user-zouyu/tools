import React, {useState} from 'react'
import './App.css'
import {FloatButton, Form, Input, Layout, List, message, Modal, Select, Switch} from "antd";
import Show from "./show/index.jsx";
import Message from "./message/index.jsx";
import {SettingOutlined, UploadOutlined} from "@ant-design/icons";
const {TextArea} = Input


function App() {
    const [username, setUsername] = useState("zou yu")
    const [group, setGroup] = useState("test")
    const [host, setHost] = useState(window.location.host)
    const [btnName, setBtnName] = useState("连接")
    const [list, setList] = useState([]);
    const [ws, setWs] = useState(null);
    const [showID, setShowID] = useState(-1)
    const [connected, setConnected] = useState(false);
    const [siderDisplay, setSiderDisplay] = useState("block")
    const [modalShow, setModalShow] = useState(false)
    const [updateCode, setUpdateCode] = useState(false)
    const [code, setCode] = useState("")
    const [language, setLanguage] = useState("Java")

    const connect = () => {
        if (username.length < 5) {
            message.error("用户名必须超过5个字符").then(() => {
            })
            return
        }
        if (group.length < 5) {
            message.error("房间号必须超过5个字符").then(() => {
            })
            return
        }
        if (!connected) {
            const ws = new WebSocket(`ws://${host}/api/connect?group=${group}&username=${username}`);
            ws.addEventListener('open', () => {
                message.info("连接成功").then(() => {
                })
                setWs(_ => ws)
                setConnected(true)
                setBtnName("断开连接")
            });

            ws.addEventListener("error", () => {
                message.error("连接出错啦").then(() => {
                })
            })
            ws.addEventListener('message', (event) => {
                const data = JSON.parse(event.data);
                message.info(data.msg).then(() => {
                })

                if (data.code === 1) {
                    setList(prevState => {
                        return [...prevState, ...data.data]
                    })
                }

                if (data.code === 2) {
                    setShowID(data.data["currentID"])
                    setList(data.data.list)
                }

                if (data.code === 4) {
                    console.log(data)
                    const cmd = data.data.command
                    switch (cmd) {
                        case "next":
                        case "prev":
                            setShowID(_ => data.data.data)
                            break
                    }
                }
            });

        } else {
            message.info("连接关闭").then(() => {
            })
            setWs(_ => null)
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
                    <Form.Item label="房间号">
                        <Input placeholder="房间号"
                               onChange={(e) => {
                                   setGroup(e.target.value)
                               }}
                               value={group}
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

            <Modal title="Update Code"
                   open={updateCode}
                   onOk={() => {
                       if (connected) {
                          ws.send(JSON.stringify({
                              "command": "updateCode",
                              "text": code,
                              "language": language
                          }))
                       } else {
                           message.error("未连接").then(_=>{})
                       }
                       setUpdateCode(false)
                   }}
                   onCancel={() => setUpdateCode(false)}
                   okText={"提交"}
            >
                <Form labelCol={{
                    span: 4
                }} wrapperCol={{
                    span: 20
                }}>
                    <Form.Item
                        label="语言"
                    >
                        <Select
                            defaultValue="Java"
                            onChange={(v) => {
                                setLanguage(v)
                            }}
                            options={[
                                { value: 'Java', label: 'Java' },
                                { value: 'C++', label: 'C++' },
                                { value: 'Go', label: 'Go' },
                                { value: 'JavaScript', label: 'JavaScript'},
                            ]}
                        />
                    </Form.Item>
                    <Form.Item
                        label="代码"
                    >
                        <TextArea
                            showCount
                            maxLength={1000}
                            style={{ height: 120, resize: 'none' }}
                            onChange={(v) => {
                                setCode(v.currentTarget.value)
                            }}
                            placeholder="code...."
                        />
                    </Form.Item>
                </Form>
            </Modal>
            <FloatButton.Group>
                <FloatButton
                    icon={<SettingOutlined />}
                    onClick={() => {
                        setModalShow(true);
                    }}/>

                <FloatButton
                    icon={<UploadOutlined />}
                    onClick={() => {
                        setUpdateCode(true);
                    }}/>
            </FloatButton.Group>
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
                                <Message key={item.id} id={idx} showId={showID} setShowId={setShowID} data={item} ws={ws}/>
                            </List.Item>
                        )}
                    />
                </Layout.Sider>

                <Layout.Content style={{height: "100vh", width: "100vw"}}>
                    <Show list={list} id={showID}/>
                </Layout.Content>
            </Layout>

        </div>
    )
}

export default App
