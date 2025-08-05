/* eslint-disable react-refresh/only-export-components */
import { CSSObject, styled } from '@mui/material/styles';
import { FC, MouseEvent, useRef } from 'react'

const MAX_MOVEMENT_LIMIT = 10;
const LOADER_SIZE = 35;

export interface HoldLoaderProps {
    onActionAfterHold: () => void;
    holdTime?: number;
    size?: number;
}

export const useActionAfterHold: UseHoldLoader = ({
    holdTime = 2500,
    onActionAfterHold,
    size = LOADER_SIZE
}) => {
    const options = useRef<Options>({});
    const holdDiv = useRef<HTMLDivElement>(null);

    const onMouseDown = (e: MouseEvent<HTMLDivElement>) => {
        if (options.current.pos || !holdDiv.current || e.button !== 0) return;
        
        options.current.pos = { x: (e.nativeEvent as any).layerX, y: (e.nativeEvent as any).layerY };
        options.current.timeoutId = setTimeout(() => {
            onActionAfterHold()
            options.current.timeoutId = undefined
        }, holdTime)

        holdDiv.current.toggleAttribute('data-helddown', true)
        holdDiv.current.style['left'] = options.current.pos.x - size / 2 +'px';
        holdDiv.current.style['top'] = options.current.pos.y - size / 2 +'px';
    }

    const onMouseLeave = () => {
        if (!options.current.pos || !holdDiv.current) return;

        options.current.pos = undefined;
        if (options.current.timeoutId) {
            clearTimeout(options.current.timeoutId)
        }
        options.current.timeoutId = undefined;

        holdDiv.current.toggleAttribute('data-helddown', false)
        holdDiv.current.style['left'] = '';
        holdDiv.current.style['top'] = '';
    }

    const onMouseMove = (e: MouseEvent<HTMLDivElement>) => {
        if (!options.current.pos) return;
        
        const movementSum = Math.abs(
            options.current.pos.x - (e.nativeEvent as any).layerX
        ) + Math.abs(
            options.current.pos.y - (e.nativeEvent as any).layerY
        )
        
        if (movementSum > MAX_MOVEMENT_LIMIT) {
            return onMouseLeave();
        }
    }

    return {
        onMouseDown,
        onMouseLeave,
        onMouseMove,
        CursorElement: () => (
            <UserPressPoint className='MusicItem_user-holding' ref={holdDiv} duration={`${holdTime}ms`} size={size}>
                <div className='loader'/>
            </UserPressPoint>
        ),
    }
}

const shouldForwardProp = (prop: string) => {
    switch (prop) {
        case 'duration':
        case 'size':
            return false
        default:
            return true
    }
}
const UserPressPoint = styled('div', { shouldForwardProp })<{ duration: string, size: number }>(({ duration, size }) => ({
    position: 'absolute',
    display: 'none',
    width: size,
    height: size,
    zIndex: 99,
    userSelect: 'none',
    pointerEvents: 'none',
    '&[data-helddown]': {
        display: 'initial',
        '@keyframes l18': {
            '0%': {
                clipPath: 'polygon(50% 50%, 50% 0, 50% 0, 50% 0, 50% 0, 50% 0%, 50% 0%)',
                borderColor: 'red',
            },
            '12.5%': {
                clipPath: 'polygon(50% 50%, 50% 0, 100% 0%, 100% 0%, 100% 0%, 100% 0%, 100% 0%)',
            },
            '25%': {
                clipPath: 'polygon(50% 50%, 50% 0, 100% 0%, 100% 50%, 100% 50%, 100% 50%, 100% 50%)',
            },
            '37.5%': {
                clipPath: 'polygon(50% 50%, 50% 0, 100% 0%, 100% 100%, 100% 100%, 100% 100%, 100% 100%)',
            },
            '50%': {
                clipPath: 'polygon(50% 50%, 50% 0, 100% 0%, 100% 100%, 50% 100%, 50% 100%, 50% 100%)',
            },
            '62.5%': {
                clipPath: 'polygon(50% 50%, 50% 0, 100% 0%, 100% 100%, 0% 100%, 0% 100%, 0% 100%)',
                borderColor: 'orange',
            },
            '75%': {
                clipPath: 'polygon(50% 50%, 50% 0, 100% 0%, 100% 100%, 0% 100%, 0% 50%, 0% 50%)',
            },
            '87.5%': {
                clipPath: 'polygon(50% 50%, 50% 0, 100% 0%, 100% 100%, 0% 100%, 0% 0%, 0% 0%)',
                borderColor: 'yellowgreen',
            },
            '100%': {
                clipPath: 'polygon(50% 50%, 50% 0, 100% 0%, 100% 100%, 0% 100%, 0% 0%, 50% 0%)',
                borderColor: 'green',
            },
        },
        '.loader': {
            width: size,
            aspectRatio: 1,
            border: `${size / 5}px solid #ddd5`,
            // borderRadius: '50%',
            position: 'relative',
            ':before': {
                content: '""',
                position: 'absolute',
                inset: `-${size / 5}px`,
                // borderRadius: '50%',
                border: `${size / 5}px solid #514b82`,
                animationName: 'l18',
                animationDuration: duration,
                animationTimingFunction: 'linear',
                animationTimeline: 'initial',

            } as CSSObject,
        },
    } as CSSObject,
}))

type UseHoldLoader = (props: HoldLoaderProps) => {
    CursorElement: FC;
    onMouseDown: (e: MouseEvent<HTMLDivElement>) => void,
    onMouseLeave: () => void,
    onMouseMove: (e: MouseEvent<HTMLDivElement>) => void,
}

export interface Options {
    pos?: {
        x: number;
        y: number;
    },
    timeoutId?: number;
}