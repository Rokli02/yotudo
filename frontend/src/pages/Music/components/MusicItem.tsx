import { memo, MouseEvent, useRef } from 'react';
import { Music } from '@src/api';
import { CSSObject, styled } from '@mui/material/styles';
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

const HOLD_TIME_IN_MS = 1500;

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
        onMouseDown,
        onMouseLeave,
        onMouseMove
    } = useActionAfterHold({ onActionAfterHold: onActionAfterHold, holdTime: HOLD_TIME_IN_MS, size: 30 })

    const loading = (target: EventTarget & HTMLSpanElement, isLoading: boolean) => {
        loadingState.current = isLoading;
        target.toggleAttribute('data-loading', isLoading);
    }

    return (
        <MusicItemContainer
            onMouseDown={onMouseDown}
            onMouseMove={onMouseMove}
            onMouseUp={onMouseLeave}
            onMouseLeave={onMouseLeave}
        >
            <CursorElement />
            <header className='MusicItem_header'>
                <div className='MusicItemHeader_np'>
                    <h1 className='MusicItemHeader_name'>{music.name}</h1>
                    { !music.published ? undefined : <span className='MusicItemHeader_published'>{ music.published }</span>}
                </div>
                { !music.album ? undefined : <h2 className='MusicItemHeader_album'>{music.album}</h2>}
            </header>
            <Divider dir='horizontal'/>
            <div className='MusicItem_content'>
                <div className='MusicItemContent_author'>
                    { music.author.name }
                </div>
                {
                    !music.contributor || music.contributor.length === 0 ? 
                        undefined : (
                        <div className='MusicItemContent_contributors'>
                            {
                                music.contributor.map((c, i) => <div key={`${i}_${c.id}`} className='MusicItemContent_contributor'> { c.name } </div>)
                            }
                        </div>
                    )
                }
            </div>
            <Divider dir='horizontal'/>
            <div className='MusicItem_footer'>
                <span className='MusicItemFooter_genre'> { music.genre.name } </span>
                <span className='MusicItemFooter_status' title={music.status.name} onMouseDown={(e) => { e.preventDefault(); e.stopPropagation() }} onClick={(e) => {
                    if (loadingState.current || !onAction) return;

                    const target = e.currentTarget

                    loading(target, true)

                    onAction(e).finally(() => {
                        loading(target, false);
                    })
                }} data-status={music.status.id}>
                    <Status/>
                </span>
            </div>
        </MusicItemContainer>
    )
})

const MusicItemContainer = styled('article')({
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
    '&.MusicItem_header': {
        height: 'min-content',
        marginBottom: '8px',
        pointerEvents: 'none',
        '&.MusicItemHeader_np': {
            display: 'flex',
            flexDirection: 'row',
            flexWrap: 'nowrap',
            justifyContent: 'space-between',
            '&.MusicItemHeader_name': {
                fontSize: '1.4rem',
                fontWeight: 500,
                color: ItemColors.primaryFontColor,
                paddingBottom: 'auto',
                marginBottom: '2px',
            },
            '&.MusicItemHeader_published': {
                fontSize: '1.25rem',
                fontWeight: 500,
                color: ItemColors.secondaryFontColor,
                paddingBottom: 'auto',
                marginBottom: '2px',
            },
        },
        '&.MusicItemHeader_album': {
            position: 'relative',
            fontSize: '1.25rem',
            width: 'fit-content',
            fontWeight: 500,
            paddingBottom: 'auto',
            marginLeft: '20px',
            marginBottom: '2px',
            color: ItemColors.secondaryFontColor,
        }
    },
    '&.MusicItem_content': {
        overflowY: 'hidden',
        flexGrow: 1,
        fontSize: '1.2rem',
        fontWeight: 400,
        color: ItemColors.primaryFontColor,
        paddingLeft: '20px',
        '&.MusicItemContent_author': {
            width: 'fit-content',
            marginBlock: '6px 6px',
            padding: '2px 6px',
            backgroundColor: '#0002',
            borderRadius: '6px',
        },
        '&.MusicItemContent_contributors': {
            display: 'flex',
            flexDirection: 'row',
            marginLeft: '8px',
            gap: '4px 6px',
            color: ItemColors.accentFontColor,
            '&.MusicItemContent_contributor': {
                width: 'fit-content',
                padding: '2px 6px',
                backgroundColor: '#0002',
                borderRadius: '6px',
            }
        },
    },
    '&.MusicItem_footer': {
        display: 'flex',
        justifyContent: 'space-between',
        flexWrap: 'nowrap',
        paddingBottom: '2px',
        paddingInline: '4px',
        '&.MusicItemFooter_genre': {
            fontSize: '1.10rem',
            fontWeight: 400,
            color: ItemColors.secondaryFontColor,
        },
        '&.MusicItemFooter_status': {
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
        }
    },
} as CustomCSS)