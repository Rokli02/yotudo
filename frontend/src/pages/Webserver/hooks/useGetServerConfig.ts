import { ServerConfig, ServerConfigParser } from "@src/api";
import { useEffect, useState } from "react";
import { GetServerConfig } from '@wailsjs/go/src/App'

export function useGetServerConfig() {
    const [refetchTrigger, setRefetchTrigger] = useState(false);
    const [serverConfig, setServerConfig] = useState<ServerConfig>({ port: 0, isHosted: false, listeningOn: "" });
    const [state, setState] = useState<'loading' | 'complete'>('loading');

    function refetch() {
        setRefetchTrigger((pre) => !pre)
    }

    useEffect(() => {
        setState('loading')

        GetServerConfig().then((config) => {
            setServerConfig(ServerConfigParser.fromGo(config))

            setState('complete')
        })

    }, [refetchTrigger])

    return { serverConfig, state, refetch }
}