import React, { Component, useState } from 'react';
import './App.css';
import { Greet, ReadMessage } from "../wailsjs/go/main/App";
import { HomeOutlined, SettingOutlined, UserOutlined } from "@ant-design/icons";
import { NavLink, HashRouter as Router, Route, Routes } from "react-router-dom";

import Home from './Home';
import Settings from './Settings';
import Account from './Account';
import { Layout, Menu, Col, Row, message, Space } from 'antd';

const { Content, Sider } = Layout;

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
        //message.info( msg.topic + ":" + msg.message )

        this.SetupMessageChannel();
    }

    SetupMessageChannel() {
        ReadMessage().then(this.onReadMesasge)
    }

    componentDidMount() {
        this.SetupMessageChannel();
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
                    <Layout className="site-layout" style={{ marginLeft: this.state.collapsed ? "80px" : "200px" }} >
                        <Content style={{ margin: '16px' , backgroundColor:"white", padding:"16px" }}>
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
                    </Layout>
                </Layout>
            </Router>
        )
    }
}

export default App
