
import { Component } from "react"
import FormRender, {connectForm } from "form-render";

import schema from "./HomeSchema.json";
import {Button} from "antd";

import {SetupN2N, SaveText, LoadText, ShutdownN2N} from "../wailsjs/go/main/App";


class Home extends Component {
    constructor(props) {
        super(props);
        this.state = {
            prevData: ""
        }
    }

    KFile = "etc/account.json"

    componentDidMount(){
        LoadText(this.KFile).then(this.onTextLoad).catch(this.onError);
    }

    onTextLoad = (data) => {
        const v = JSON.parse(data);
        this.setState({
            prevData: JSON.stringify(v, null, "  ")
        })
        this.props.form.setValues(v);
    }

    onError = (reason) => {
        console.log("Error:", reason)
    }

    onFinish = (formData, errors) => {
        const buff = JSON.stringify(formData, null, "  ");
        if (buff != this.state.prevData) {
            SaveText(this.KFile, buff).then(this.trySetupEdge);
            this.setState(
                {prevData:buff}
            )
        } else {
            this.trySetupEdge();
        }
    }

    trySetupEdge = (error) => {
        if (error != null) {
            console.log(error)
        }
        SetupN2N().then(this.onError).catch(this.onError);
    }


    shutdownEdge = () => {
        ShutdownN2N();
    }

    render (){
        const {form} = this.props;
        return (
            <div>
                <FormRender form={form} schema={schema} onFinish={this.onFinish} />
                <Button style={{"float":"right"}} type="primary" onClick={form.submit}>Connect</Button>
                <Button style={{"float": "right"}} type="default" onClick={this.shutdownEdge}>Disconnect</Button>
            </div>
        )
    }
}

export default connectForm(Home);