import logo from './logo.svg';
import './App.css';
import RunForm from "./Form";
import { DockerMuiThemeProvider } from '@docker/docker-mui-theme';
import CssBaseline from '@mui/material/CssBaseline';

function App() {
    return (
        <DockerMuiThemeProvider>
            <CssBaseline />
            <div className="App">
                <p>Storj local regisry is a local docker registry which bridges which stores all the data in the decentralized cloud. Please start the registry with your access grant and access images with prefix localhost:9999/</p>
                <RunForm></RunForm>
            </div>
        </DockerMuiThemeProvider>
    );
}

export default App;
