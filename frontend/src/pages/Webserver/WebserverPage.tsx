import { Box, SxProps, Theme } from '@mui/material';
import { FC } from 'react';
import ServerStart from './components/ServerStart.component';
import { useGetServerConfig } from './hooks/useGetServerConfig';
import { LoadingPage } from '../Common';
import { StartServer, StopServer } from '@wailsjs/go/src/App';
import ServerRunning from './components/ServerRunning.component';
import { ServerConfig, ServerConfigParser } from '@src/api';

export const WebserverPage: FC = () => {
    const { serverConfig, state, refetch } = useGetServerConfig()

    async function startServer(v: ServerConfig) {
        await StartServer(ServerConfigParser.toGo(v));
        refetch();
        return true;
    }

    async function stopServer() {
        await StopServer();
        refetch()
    }

    return <Box sx={ContainerStyle}>
        <h1>Webszerver</h1>
        { state === 'loading' ?
            <LoadingPage /> :
            !serverConfig.listeningOn ?
            <ServerStart config={serverConfig} startServer={startServer} /> :
            <ServerRunning config={serverConfig} stopServer={stopServer} />
        }
    </Box>
}

export default WebserverPage;

const ContainerStyle: SxProps<Theme> = {
    position: 'relative',
    padding: '1rem 1rem',
    '& > h1': {
        textAlign: 'center',
    },
}