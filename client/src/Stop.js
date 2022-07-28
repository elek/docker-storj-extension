import React from "react";
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import {Stack} from "@mui/material";
import {Service} from "./Service";

const Stop = (props) => {
    const cb = props.cb
    const stop = (event) => {
        event.preventDefault();
        Service.StopContainer();
        setTimeout(cb, 1000)
    }


    return (
        <Stack spacing={{xs: 1, sm: 2, md: 4}} xl={8}>
            <p>Storj based registry is up and running</p>
            <Button variant="outlined" onClick={stop}>Stop registry</Button>
        </Stack>
    );
}


export default Stop;
