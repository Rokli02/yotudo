import { FC, MouseEvent, useState } from 'react'
import { Box, Link, SxProps, Theme, Typography } from '@mui/material';
import { ContentCopy } from '@mui/icons-material';
import { ServerConfig } from '@src/api';
import { Button } from '@src/components/form';
import { ClipboardSetText } from '@wailsjs/runtime';
import usePopupBox from '@src/hooks/usePopupBox';

export interface ServerRunningProps {
    config: ServerConfig;
    stopServer: () => Promise<void>;
}

export const ServerRunning: FC<ServerRunningProps> = ({ config, stopServer }) => {
    const { PopupBox, onClick: onLinkClick } = usePopupBox("Másolva", () => ClipboardSetText(config.listeningOn))
    const [isLoading, setLoading] = useState(false)

    async function onStopBtnClick(e: MouseEvent<HTMLButtonElement>) {
        e.preventDefault();

        setLoading(true);

        await stopServer().finally(() => {
            setLoading(false);
        })
    }

    return (
        <Box sx={ContainerStyle}>
            <Typography variant='h3' sx={ContentStyle.Title1}>A szerver az alábbi linken érhető el:</Typography>
            <Box sx={ContentStyle.LinkContainer}>
                <Link
                    align='center'
                    sx={ContentStyle.Link}
                    underline='hover'
                    target="_blank"
                    onClick={onLinkClick}
                >
                    {config.listeningOn}
                    <ContentCopy />
                </Link>
                <PopupBox />
            </Box>
            <Box sx={ContentStyle.QRCode}>QR code placeholder</Box>
            <Box sx={ContentStyle.StopBtnContainer}>
                <Button
                    sx={ContentStyle.StopBtn}
                    variant='outlined'
                    type='button'
                    color='error'
                    disabled={isLoading}
                    onClick={onStopBtnClick}
                >Leállítás</Button>
            </Box>
        </Box>
    )
}

export default ServerRunning;

const ContainerStyle: SxProps<Theme> = {
    height: 'calc(85vh - 100px)',
    marginInline: 'auto',
    paddingTop: '32px',
    maxWidth: '800px',
};
const ContentStyle = {
    Title1: {
        fontSize: '1.5rem',
        textAlign: 'center',
    },
    LinkContainer: {
        position: 'relative',
        marginTop: '8px',
        textAlign: 'center',
    },
    Link: {
        position: 'relative',
        fontSize: '1.35rem',
        cursor: 'pointer',
        color: 'var(--primary-color)',
        fontWeight: 500,
        '& svg': {
            marginLeft: '4px',
            width: '20px',
            height: '20px',
            userSelect: 'none',
            pointerEvents: 'none',
        },
    },
    LinkCopyBox: {
        position: 'absolute',
        left: '50%',
        top: '100%',
        width: 'fit-content',
        padding: '2px 12px 4px',
        borderRadius: '4px',
        backgroundColor: 'var(--primary-color)',
        userSelect: 'none',
        pointerEvents: 'none',
        transition: 'opacity ease-in-out 200ms',
        opacity: 0,
        '&[data-copied]': {
            opacity: 1,
        }
    },
    QRCode: {
        fontSize: '1rem',
        textAlign: 'center',
        marginTop: '20px',
    },
    StopBtnContainer: {
        width: '100%',
        height: '80px',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'end',
        marginInline: 'auto',
    },
    StopBtn: {
        height: 'min-content',
        paddingInline: '3.5ch',
        paddingBlock: '6px 4px',
        justifySelf: 'center',
    }
} satisfies Record<string, SxProps<Theme>>