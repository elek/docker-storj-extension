import './App.css';
import React, {useEffect, useState} from "react";
import RunForm from "./CreateNew";
import Start from "./Start";
import Stop from "./Stop";
import {DockerMuiThemeProvider} from '@docker/docker-mui-theme';
import CssBaseline from '@mui/material/CssBaseline';
import Grid from '@mui/material/Grid';
import {Service} from "./Service";

const App = () => {
    const [status, setStatus] = useState("")
    const refresh = function () {

        Service.CheckStatus().then((r) => {
            console.log("refreshing state " + r)
            setStatus(r)
        })
    }
    useEffect(() => {
        refresh()
    }, [""])


    return (
        <DockerMuiThemeProvider>
            <CssBaseline/>
            <Grid container>
                <Grid item lg={4}>
                    <div className="App">
                        <p>{status}Storj local registry is a local docker registry stores all the data in the
                            decentralized cloud.</p>

                        <p>After starting the local registry you can push/pull to/from Storj bucket with using localhost:9999 prefix (eg. docker run localhost:9999/name/image)</p>
                        {status === "missing" &&
                            <RunForm cb={refresh}></RunForm>
                        }
                        {status === "stopped" &&
                            <Start cb={refresh}></Start>
                        }
                        {status === "running" &&
                            <Stop cb={refresh}></Stop>
                        }
                    </div>
                </Grid>
            </Grid>
        </DockerMuiThemeProvider>
    );
}


export default App;
