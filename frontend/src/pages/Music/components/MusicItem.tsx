import { memo, MouseEvent, useRef } from 'react';
import { Music } from '@src/api';
import { styled } from '@mui/material/styles';
import { Divider, StatusActionIcon } from '@src/components/common';
import { useActionAfterHold } from '@src/hooks/useHoldLoader';
import { CustomCSS } from '@src/components/common/interface';

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
        <Container
            onMouseDown={onMouseDown}
            onMouseMove={onMouseMove}
            onMouseUp={onMouseUp}
            onMouseLeave={onMouseLeave}
        >
            { music.picName
                ? <ThumbnailBackground src={`image/${music.picName}`}/>
                : undefined
            }
            <CursorElement />
            <ItemHeader.Container>
                <ItemHeader.TopRow>
                    <ItemHeader.Name>{music.name}</ItemHeader.Name>
                    { !music.published ? undefined : <ItemHeader.Published>{ music.published }</ItemHeader.Published>}
                </ItemHeader.TopRow>
                { !music.album ? undefined : <ItemHeader.Album>{music.album}</ItemHeader.Album>}
            </ItemHeader.Container>
            <Divider dir='horizontal'/>
            <ItemContent.Container>
                <ItemContent.Author>
                    { music.author.name }
                </ItemContent.Author>
                {
                    !music.contributor || music.contributor.length === 0 ? 
                        undefined : (
                        <ItemContent.Contributors>
                            {
                                music.contributor.map((c, i) => <ItemContent.ContributorItem key={`${i}_${c.id}`}> { c.name } </ItemContent.ContributorItem>)
                            }
                        </ItemContent.Contributors>
                    )
                }
            </ItemContent.Container>
            <Divider dir='horizontal'/>
            <ItemFooter.Container>
                <ItemFooter.Genre> { music.genre.name } </ItemFooter.Genre>
                <ItemFooter.Status title={music.status.name} onMouseDown={onChildMouseDown} onClick={(e) => {
                    if (loadingState.current || !onAction) return;

                    const target = e.currentTarget

                    loading(target, true)

                    onAction(e).finally(() => {
                        loading(target, false);
                    })
                }} data-status={music.status.id}>
                    <Status/>
                </ItemFooter.Status>
            </ItemFooter.Container>
        </Container>
    )
})

//#region Components
const Container = styled('article')({
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
} as CustomCSS);
const ItemHeader = {
    Container: styled('header')({
        height: 'min-content',
        marginBottom: '8px',
        pointerEvents: 'none',
    }),
    TopRow: styled('div')({
        display: 'flex',
        flexDirection: 'row',
        flexWrap: 'nowrap',
        justifyContent: 'space-between',
    }),
    Name: styled('h2')({
        fontSize: '1.3rem',
        fontWeight: 500,
        color: ItemColors.primaryFontColor,
        paddingBottom: 'auto',
        marginBottom: '2px',
    }),
    Published: styled('span')({
        fontSize: '1.25rem',
        fontWeight: 500,
        color: ItemColors.secondaryFontColor,
        paddingBottom: 'auto',
        marginBottom: '2px',
    }),
    Album: styled('h3')({
        position: 'relative',
        fontSize: '1.25rem',
        width: 'fit-content',
        fontWeight: 500,
        paddingBottom: 'auto',
        marginLeft: '20px',
        marginBottom: '2px',
        color: ItemColors.secondaryFontColor,
    }),
};
const ItemContent = {
    Container: styled('div')({
        // overflowY: 'hidden',
        flexGrow: 1,
        fontSize: '1.2rem',
        fontWeight: 400,
        color: ItemColors.primaryFontColor,
        paddingLeft: '20px',
    }),
    Author: styled('div')({
        width: 'fit-content',
        marginBlock: '6px 6px',
        padding: '2px 6px',
        backgroundColor: '#0002',
        borderRadius: '6px',
    }),
    Contributors: styled('div')({
        display: 'flex',
        flexDirection: 'row',
        marginLeft: '8px',
        gap: '4px 6px',
        color: ItemColors.accentFontColor,
    }),
    ContributorItem: styled('div')({
        width: 'fit-content',
        padding: '2px 6px',
        backgroundColor: '#0002',
        borderRadius: '6px',
    }),
};
const ItemFooter = {
    Container: styled('div')({
        display: 'flex',
        justifyContent: 'space-between',
        flexWrap: 'nowrap',
        paddingBottom: '2px',
        paddingInline: '4px',
    }),
    Genre: styled('span')({
        fontSize: '1.10rem',
        fontWeight: 400,
        color: ItemColors.secondaryFontColor,
    }),
    Status: styled('span')({
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
    }),
};
const ThumbnailBackground = styled('img')({
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
})
//#endregion Components