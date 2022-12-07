
import { Component } from "react"
import FormRender, {connectForm} from "form-render";

import schema from "./AccountSchema.json";
import {Button} from "antd";
import {ChangePassword} from "../wailsjs/go/main/App";


class Account extends Component {
    constructor(props) {
        super(props);
    }

    componentDidMount(){
    }

    onFinish = (formData, errors) => {
        if (formData.oldPassword == formData.newPassword) {
            return;
        }
        ChangePassword(formData);
    }

    render (){
        const {form} = this.props;
        return (
            <div>
                <FormRender form={form} schema={schema} onFinish={this.onFinish} />
                <Button style={{float:"right"}} type="primary" onClick={form.submit}>Change</Button>
            </div>
        )
    }
}

export default connectForm(Account);