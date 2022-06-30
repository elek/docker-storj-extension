import React from "react";
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import {Stack} from "@mui/material";
import {Service} from "./Desktop";

class RunForm extends React.Component {
    constructor(props) {
        super(props);
        this.state = {bucket: '',accessgrant: ''};

        this.handleChangeBucket = this.handleChangeBucket.bind(this);
        this.handleChangeAccessGrant = this.handleChangeAccessGrant.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);

        this.service = new Service();
    }

    handleChangeBucket(event) {
        this.setState({bucket: event.target.value});
    }

    handleChangeAccessGrant(event) {
        this.setState({accessgrant: event.target.value});
    }

    handleSubmit(event) {
        event.preventDefault();
        this.service.StartContainer(this.state.bucket, this.state.accessgrant);
    }

    render() {
        return (
            <Stack spacing={{ xs: 1, sm: 2, md: 4}} xl={8}>
                <TextField label="bucket" value={this.state.bucket} variant="filled" onChange={this.handleChangeBucket}/>
                <TextField label="access grant" value={this.state.accessgrant} variant="filled" onChange={this.handleChangeAccessGrant} />
                <Button variant="outlined" onClick={this.handleSubmit}>Start registry</Button>
            </Stack>
        );
    }
}

export default RunForm;
