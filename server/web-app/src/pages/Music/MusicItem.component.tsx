import { FC, memo, MouseEvent, useRef } from 'react';
import { Music } from '@src/api/index.js';
import { Download, CachedRounded } from "@mui/icons-material";
import { SxProps, Theme } from '@mui/material/styles';
import { Divider } from '@src/components/common/index.js';
import { Box, Typography } from '@mui/material';
import { HoldDownAction, onChildMouseDown } from '@src/components/common/HoldDownAction.js';

const ItemColors = {
    container: 'var(--primary-color)',
    containerHeldDown: '#611818',
    containerShadow: '#0005',
    dividerColor: '#333333',
    primaryFontColor: 'var(--font-color)',
    secondaryFontColor: '#858585',
    accentFontColor: '#a6a6a6',
}

const HOLD_TIME_IN_MS = 850;

interface MusicItemProps {
    music: Music;
    onAction?: (event: MouseEvent<HTMLSpanElement>) => Promise<void>;
    onActionAfterHold?: () => void;
}

function DefaultOnActionAfterHold() { console.log(`After ${HOLD_TIME_IN_MS}ms of holding, there is an action`) }

export const MusicItem = memo(({
    music,
    onAction,
    onActionAfterHold = DefaultOnActionAfterHold,
}: MusicItemProps) => {
    const loadingState = useRef<boolean>(false)

    const loading = (target: EventTarget & HTMLSpanElement, isLoading: boolean) => {
        loadingState.current = isLoading;
        target.toggleAttribute('data-loading', isLoading);
    }

    const StatusIcon = StatusActionIcon[music.status.id]

    return (
        <HoldDownAction
            sx={ContainerStyle}
            component={'article'}
            onActionAfterHold={onActionAfterHold}
            holdTime={HOLD_TIME_IN_MS}
            size={60}
        >
            { music.picName
                ? <Box component={'img'} sx={ThumbnailBackgroundStyle} src={`image/${music.picName}`}/>
                : undefined
            }
            <Box component={'header'} sx={ItemHeader.ContainerStyle}>
                <Box sx={ItemHeader.TopRowStyle}>
                    <Typography component={'h2'} sx={ItemHeader.NameStyle}>{music.name}</Typography>
                    { !music.published ? undefined : <Typography component={'span'} sx={ItemHeader.PublishedStyle}>{ music.published }</Typography>}
                </Box>
                { !music.album ? undefined : <Typography component={'h3'} sx={ItemHeader.AlbumStyle}>{music.album}</Typography>}
            </Box>
            <Divider dir='horizontal'/>
            <Box sx={ItemContent.ContainerStyle}>
                <Box sx={ItemContent.ItemStyle}>{ music.author.name }</Box>
                {
                !music.contributor || music.contributor.length === 0 ? 
                    undefined : 
                    music.contributor.map((c, i) => <Box key={`${i}_${c.id}`} sx={ItemContent.ItemStyle}> { c.name } </Box>)
                }
            </Box>
            <Divider dir='horizontal'/>
            <Box sx={ItemFooter.ContainerStyle}>
                <Box component={'span'} sx={ItemFooter.GenreStyle}> { music.genre.name } </Box>
                <Box component={'span'} sx={ItemFooter.StatusStyle} title={music.status.name} onMouseDown={onChildMouseDown} onClick={(e) => {
                    if (loadingState.current || !onAction) return;

                    const target = e.currentTarget

                    loading(target, true)

                    onAction(e).finally(() => {
                        loading(target, false);
                    })
                }} data-status={music.status.id}>
                    <StatusIcon />
                </Box>
            </Box>
        </HoldDownAction>
    )
})

//#region Components
const AnimatedCachedRounded: FC = () => {
    return <CachedRounded sx={{
    '@keyframes spin': {
        '0%': {
            transform: 'rotate(0deg)',
        },
        '100%': {
            transform: 'rotate(360deg)',
        },
    },
    animationName: 'spin',
    animationIterationCount: 'infinite',
    animationDirection: 'reverse',
    animationDuration: '1.5s'
}} />
}
const StatusActionIcon: Record<number, FC> = {
    1: AnimatedCachedRounded,
    2: Download,
}
const ContainerStyle: SxProps<Theme> = {
    position: 'relative',
    cursor: 'default',
    userSelect: 'none',
    display: 'flex',
    flexDirection: 'column',
    width: '380px',
    minHeight: 'max-content',
    height: '280px',
    maxHeight: '450px',
    backgroundColor: ItemColors.container,
    boxShadow: `3px 3px 9px ${ItemColors.containerShadow}`,
    borderRadius: '16px',
    padding: '6px 10px',
    transition: 'background-color 250ms',
    '&:has([data-helddown])': {
        cursor: 'none',
        transition: `background-color ${HOLD_TIME_IN_MS}ms linear 500ms`,
        backgroundColor: ItemColors.containerHeldDown,
    },
};
const ItemHeader = {
    ContainerStyle: {
        height: 'min-content',
        marginBottom: '2px',
        pointerEvents: 'none',
    },
    TopRowStyle:{
        display: 'flex',
        flexDirection: 'row',
        flexWrap: 'nowrap',
        justifyContent: 'space-between',
    },
    NameStyle: {
        fontSize: '1.15rem',
        fontWeight: 400,
        color: ItemColors.primaryFontColor,
        paddingBottom: 'auto',
    },
    PublishedStyle: {
        fontSize: '1.10rem',
        fontWeight: 400,
        color: ItemColors.secondaryFontColor,
        paddingBottom: 'auto',
    },
    AlbumStyle: {
        position: 'relative',
        fontSize: '1.10rem',
        width: 'fit-content',
        fontWeight: 400,
        paddingBottom: 'auto',
        marginLeft: '16px',
        color: ItemColors.secondaryFontColor,
    },
} satisfies Record<string, SxProps<Theme>>;
const ItemContent = {
    ContainerStyle: {
        display: 'flex',
        flexWrap: 'wrap',
        flexDirection: 'row',
        flexGrow: 1,
        justifyContent: 'start',
        alignContent: 'start',
        fontSize: '1rem',
        fontWeight: 300,
        color: ItemColors.primaryFontColor,
        gap: '4px 6px',
        paddingLeft: '16px',
        paddingBlock: '2px',
    },
    ItemStyle: {
        width: 'fit-content',
        height: 'min-content',
        padding: '1px 8px',
        backgroundColor: '#0002',
        borderRadius: '6px',
        overflowX: 'hidden',
        textOverflow: 'ellipsis',
        textWrap: 'nowrap',
    },
} satisfies Record<string, SxProps<Theme>>;
const ItemFooter = {
    ContainerStyle: {
        display: 'flex',
        justifyContent: 'space-between',
        flexWrap: 'nowrap',
        paddingBottom: '2px',
        paddingInline: '4px',
    },
    GenreStyle: {
        fontSize: '1.10rem',
        fontWeight: 400,
        color: ItemColors.secondaryFontColor,
    },
    StatusStyle: {
        padding: '4px',
        borderRadius: '8px',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        '&[data-status="0"], &[data-status="2"]': {
            transition: 'box-shadow 200ms, background-color 200ms',
            cursor: 'pointer',
            boxShadow: '1px 1px 2px #0000005a',
            ':hover': {
                backgroundColor: '#0000000e',
                boxShadow: 'unset',
            },
            ':active': {
                backgroundColor: '#0000001a',
                boxShadow: '0px 0px 5px 1px #0000001a',
            },
        },
        '&[data-loading]': {
            cursor: 'default',
            '--disabled-color': '#8885',
            backgroundColor: 'var(--disabled-color)',
            boxShadow: 'unset',
            ':hover': {
                backgroundColor: 'var(--disabled-color)',
            },
            ':active': {
                backgroundColor: 'var(--disabled-color)',
            },
        },
        '&[data-status="1"]:hover': {
            ':hover > svg:active': {
                cursor: 'none',
                animationPlayState: 'running',
                animationTimingFunction: 'linear',
                animationDuration: '750ms',
            },
        },
    },
} satisfies Record<string, SxProps<Theme>>;
const ThumbnailBackgroundStyle: SxProps<Theme> = {
    position: 'absolute',
    left: 0,
    top: 0,
    height: '100%',
    width: '100%',
    pointerEvents: 'none',
    opacity: .23,
    backgroundBlendMode: 'darken',
    objectFit: 'cover',
    maskImage: `radial-gradient(
        circle at center,
        rgba(255,255,255,1) 0%,
        rgba(255,255,255,0) 67%
    )`,
};
//#endregion Components