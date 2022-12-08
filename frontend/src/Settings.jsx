
import { Component } from "react"
import FormRender, {connectForm} from "form-render";

import schema from "./SettingsSchema.json";

import {Button} from "antd";

import {SaveText, LoadText} from "../wailsjs/go/main/App";


class Settings extends Component {
    constructor(props) {
        super(props);
    }

    KFile = "etc/config.json"

    componentDidMount(){
        LoadText(this.KFile).then(this.onTextLoad).catch(this.onError);
    }

    onTextLoad = (data) => {
        const v = JSON.parse(data);
        this.props.form.setValues(v);
    }

    onError = (reason) => {
        console.log(reason)
    }

    onFinish = (formData, errors) => {
        const buff = JSON.stringify(formData, null, "  ");
        console.log(formData)
        SaveText(this.KFile, buff);
    }

    render (){
        const {form} = this.props;
        return (
            <div>
                <FormRender form={form} schema={schema} onFinish={this.onFinish} />
                <Button style={{float:"right"}} type="primary" onClick={form.submit}>Save</Button>
            </div>
        )
    }
}

export default connectForm(Settings);