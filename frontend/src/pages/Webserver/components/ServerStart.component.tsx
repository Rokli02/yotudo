import { Box, SxProps, Theme } from '@mui/material';
import { ServerConfig } from '@src/api';
import { Button, FormControl, InputLabel } from '@src/components/form';
import { Form, FormCheckbox, FormConstraints, FormInput } from '@src/contexts/form';
import { FC, useState } from 'react'

export interface ServerStartProps {
    config?: ServerConfig;
    startServer: (value: ServerConfig) => Promise<boolean>;
}

export const ServerStart: FC<ServerStartProps> = ({ config, startServer }) => {
    const [isLoading, setLoading] = useState(false)

    async function onSubmit(v: any) {
        setLoading(true)

        return await startServer(v).finally(() => {
            setLoading(false)
        })
    }

    return (
        <Box sx={ContainerStyle}>
            <Form sx={FormContent.Container} onSubmit={onSubmit} constraints={formConstraints} transformFlatObjectTo={transformObject}>
                <FormControl sx={FormContent.Port}>
                    <InputLabel>Port</InputLabel>
                    <FormInput name='port' type='number' value={config?.port}/>
                </FormControl>
                <FormCheckbox sx={FormContent.Host} name='host' label='Megosztás a hálózaton' value={config?.isHosted} />
                <div/>
                <Button
                    sx={FormContent.LaunchButton}
                    variant='contained'
                    color='secondary'
                    type='submit'
                    disabled={isLoading}
                >Indítás</Button>
            </Form>
        </Box>
    )
}

export default ServerStart;

const formConstraints: FormConstraints = {
    port: (_value, errors) => {
        const value = Number(_value)
        if (isNaN(value)) errors.push("A mezőnek Kötelezően számnak kell lennie")
        if (value < 0 || value > 65535) errors.push("A számnak nagyobb vagy egyenlőnek kell lenni, mint 0 és kisebbnek mint 65535 vagy egyenlő")
    }
}

interface RawServerConfig {
    host: boolean;
    port: string;
}

function transformObject(value: RawServerConfig): ServerConfig {
    const conFing: ServerConfig = {
        listeningOn: '',
        isHosted: value.host,
        port: Number(value.port),
    }

    return conFing
}

const ContainerStyle: SxProps<Theme> = {
    height: 'calc(85vh - 100px)',
    marginInline: 'auto',
    paddingTop: '32px',
    maxWidth: '800px',
};
const FormContent = {
    Container: {
        height: '100%',
        display: 'grid',
        gridTemplateColumns: '1fr 1fr',
        gridTemplateRows: '50px 1fr 40px',
        gridTemplateAreas: `
            "port host"
            "spacer spacer"
            "launch launch"
        `,
        gap: '24px 16px',
    },
    Port: {
        gridArea: 'port',
    },
    Host: {
        gridArea: 'host',
        flexDirection: 'row-reverse',
        justifyContent: 'center',
    },
    LaunchButton: {
        gridArea: 'launch',
        marginInline: 'auto',
        paddingInline: '3.5ch',
        fontSize: '1.05rem',
    }
} satisfies Record<string, SxProps<Theme>>