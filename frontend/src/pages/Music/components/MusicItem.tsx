import { memo, MouseEvent, useRef } from 'react';
import { Music } from '@src/api';
import { SxProps, Theme } from '@mui/material/styles';
import { Divider, StatusActionIcon } from '@src/components/common';
import { useActionAfterHold } from '@src/hooks/useHoldLoader';
import { Box, Typography } from '@mui/material';

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

export const MusicItem = memo(({ music, onAction, onActionAfterHold = () => { console.log(`After ${HOLD_TIME_IN_MS}ms of holding, there is an action`) } }: MusicItemProps) => {
    const loadingState = useRef<boolean>(false)
    const Status = StatusActionIcon[music.status.id];
    const {
        CursorElement,
        onChildMouseDown,
        onMouseDown,
        onMouseLeave,
        onMouseMove,
        onMouseUp,
    } = useActionAfterHold({ onActionAfterHold: onActionAfterHold, holdTime: HOLD_TIME_IN_MS, size: 30 })

    const loading = (target: EventTarget & HTMLSpanElement, isLoading: boolean) => {
        loadingState.current = isLoading;
        target.toggleAttribute('data-loading', isLoading);
    }

    return (
        <Box
            sx={ContainerStyle}
            component='article'
            onMouseDown={onMouseDown}
            onMouseMove={onMouseMove}
            onMouseUp={onMouseUp}
            onMouseLeave={onMouseLeave}
        >
            { music.picName
                ? <Box component={'img'} sx={ThumbnailBackgroundStyle} src={`image/${music.picName}`}/>
                : undefined
            }
            <CursorElement />
            <Box component={'header'} sx={ItemHeader.ContainerStyle}>
                <Box sx={ItemHeader.TopRowStyle}>
                    <Typography component={'h2'} sx={ItemHeader.NameStyle}>{music.name}</Typography>
                    { !music.published ? undefined : <Typography component={'span'} sx={ItemHeader.PublishedStyle}>{ music.published }</Typography>}
                </Box>
                { !music.album ? undefined : <Typography component={'h3'} sx={ItemHeader.AlbumStyle}>{music.album}</Typography>}
            </Box>
            <Divider dir='horizontal'/>
            <Box sx={ItemContent.ContainerStyle}>
                <Box sx={ItemContent.AuthorStyle}>
                    { music.author.name }
                </Box>
                {
                    !music.contributor || music.contributor.length === 0 ? 
                        undefined : (
                        <Box sx={ItemContent.ContributorsStyle}>
                            {
                                music.contributor.map((c, i) => <Box key={`${i}_${c.id}`} sx={ItemContent.ContributorItemStyle}> { c.name } </Box>)
                            }
                        </Box>
                    )
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
                    <Status/>
                </Box>
            </Box>
        </Box>
    )
})

//#region Components
const ContainerStyle: SxProps<Theme> = {
    position: 'relative',
    cursor: 'default',
    userSelect: 'none',
    display: 'flex',
    flexDirection: 'column',
    width: '380px',
    minHeight: 'max-content',
    height: '300px',
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
        marginBottom: '8px',
        pointerEvents: 'none',
    },
    TopRowStyle:{
        display: 'flex',
        flexDirection: 'row',
        flexWrap: 'nowrap',
        justifyContent: 'space-between',
    },
    NameStyle: {
        fontSize: '1.3rem',
        fontWeight: 500,
        color: ItemColors.primaryFontColor,
        paddingBottom: 'auto',
        marginBottom: '2px',
    },
    PublishedStyle: {
        fontSize: '1.25rem',
        fontWeight: 500,
        color: ItemColors.secondaryFontColor,
        paddingBottom: 'auto',
        marginBottom: '2px',
    },
    AlbumStyle: {
        position: 'relative',
        fontSize: '1.25rem',
        width: 'fit-content',
        fontWeight: 500,
        paddingBottom: 'auto',
        marginLeft: '20px',
        marginBottom: '2px',
        color: ItemColors.secondaryFontColor,
    },
} satisfies Record<string, SxProps<Theme>>;
const ItemContent = {
    ContainerStyle: {
        flexGrow: 1,
        fontSize: '1.2rem',
        fontWeight: 400,
        color: ItemColors.primaryFontColor,
        paddingLeft: '20px',
    },
    AuthorStyle: {
        width: 'fit-content',
        marginBlock: '6px 6px',
        padding: '2px 6px',
        backgroundColor: '#0002',
        borderRadius: '6px',
    },
    ContributorsStyle: {
        display: 'flex',
        flexDirection: 'row',
        flexWrap: 'wrap',
        marginLeft: '8px',
        gap: '4px 6px',
        color: ItemColors.accentFontColor,
    },
    ContributorItemStyle: {
        width: 'fit-content',
        padding: '2px 6px',
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
    'object-fit': 'cover',
    'mask-image': `radial-gradient(
        circle at center,
        rgba(255,255,255,1) 0%,
        rgba(255,255,255,0) 67%
    )`,
};
//#endregion Components