import React, { Component, useState } from 'react';
import './App.css';
//import { Greet } from "../wailsjs/go/main/App";
import { EventsOn, BrowserOpenURL } from '../wailsjs/runtime/runtime';
import { HomeOutlined, SettingOutlined, UserOutlined } from "@ant-design/icons";
import { NavLink, HashRouter as Router, Route, Routes } from "react-router-dom";

import Home from './Home';
import Settings from './Settings';
import Account from './Account';
import { Layout, Menu, Col, Row, message, Space, Button } from 'antd';

const { Content, Sider, Footer } = Layout;

class App extends Component {
    state = {
        collapsed: true,
        log: ""
    }

    MenuItems = [
        {
            label: <NavLink to="/">Home</NavLink>,
            key: "home",
            icon: <HomeOutlined />,
            component: Home,
            path: "/",
            title: "Home"
        },
        {
            label: <NavLink to="/settings">Settings</NavLink>,
            key: "settings",
            icon: <SettingOutlined />,
            component: Settings,
            path: "/settings",
            title: "Settings"
        },
        {
            label: <NavLink to="/account">Account</NavLink>,
            key: "account",
            icon: <UserOutlined />,
            component: Account,
            path: "/account",
            title: "Account"
        }
    ]



    setCollapsed = () => {
        this.setState({ collapsed: !this.state.collapsed });
    }

    renderRoutes = (items) => {
        return items.map(
            (v, i) => {
                if (v.children) {
                    return this.renderRoutes(v.children)
                } else {
                    return <Route path={v.path} element={React.createElement(v.component, {
                        log: this.state.log
                    }) } key={v.key} />
                }
            }
        )
    }

    onReadMesasge = (msg) => {
        console.log(msg)
        this.setState(
            {
                log: this.state.log +  msg.topic + ":" + msg.message
            }
        )

    }

    componentDidMount() {
        EventsOn("message", this.onReadMesasge)
    }

    BrowserOpenLinks = () => {
        /*
        <a href='https://github.com/ntop/n2n'>N2N</a>{" | "}
        <a href='https://wails.io/'>Wails</a> {" | "}
        <a href='https://4x.ant.design/index-cn'>Antd</a>
        */
        var links = [
            {
                title:"N2N",
                url:"https://github.com/ntop/n2n"
            },
            {
                title:"Wails",
                url:"https://wails.io/"
            },
            {
                title:"Antd",
                url:"https://4x.ant.design/index-cn"
            }
        ]

        var last = links[links.length-1]

        return (
            <span>
                {links.slice(0, links.length-1).map((data, i)=> {
                    return <span><a href="#" onClick={()=>{BrowserOpenURL(data.url)}}> {data.title} </a> {" | "} </span>
                })}

                <a href="#" onClick={()=>{BrowserOpenURL(last.url)}}> {last.title} </a>

            </span>
        )
    }


    render() {
        return (
                <Router>
                <Layout hasSider style={{ minHeight: "100vh" }} >
                    <Sider theme='dark' collapsible collapsed={this.state.collapsed} onCollapse={this.setCollapsed} style={{
                        overflow: 'auto',
                        height: '100vh',
                        position: 'fixed',
                        left: 0,
                        top: 0,
                        bottom: "48px",
                        zIndex: 10,
                    }}>
                        <div className="logo" >
                            {
                                this.state.collapsed ? "N2N" : "N2N User Edge"
                            }
                        </div>
                        <Menu theme='dark' mode="inline" items={this.MenuItems} defaultSelectedKeys={"home"} />
                    </Sider>
                    <Layout className="site-layout" style={{ marginLeft: this.state.collapsed ? "80px" : "200px", backgroundColor:"white" }} >
                        <Content style={{ margin: '0' , backgroundColor:"white", padding:"1.6rem" }}>
                            <Row>
                                <Col span={20} offset={2}>
                                    <Space>
                                        <Routes>
                                            {
                                                this.renderRoutes(this.MenuItems)
                                            }
                                        </Routes>
                                    </Space>
                                </Col>
                            </Row>

                        </Content>
                        <Footer style={{backgroundColor:"white"}}>
                            <Row>
                                <Col span={24} style={{textAlign:"center"}}> Powered by  {this.BrowserOpenLinks()}</Col>
                            </Row>

                        </Footer>
                    </Layout>
                </Layout>
            </Router>
        )
    }
}

export default App
