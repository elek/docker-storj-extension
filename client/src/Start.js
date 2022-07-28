import React from "react";
import Button from '@mui/material/Button';
import {Stack} from "@mui/material";
import {Service} from "./Service";

const Start = (props) => {
    const cb = props.cb
    const start = (event) => {
        event.preventDefault();
        Service.StartContainer();
        setTimeout(cb, 1000)
    }

    return (
        <Stack spacing={{ xs: 1, sm: 2, md: 4}} xl={8}>
            <p>Storj based registry image is already created, but not running.</p>
            <Button variant="outlined" onClick={start}>Start registry</Button>
        </Stack>
    );
}


export default Start;
