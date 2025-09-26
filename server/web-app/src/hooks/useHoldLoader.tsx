import { Box } from '@mui/material';
import {
    ComponentProps,
    FC,
    forwardRef,
    MouseEventHandler,
    TouchEventHandler,
    useRef,
} from 'react'

const MAX_MOVEMENT_LIMIT = 17;
const LOADER_SIZE = 35;
const START_TOUCH_EVENT_AFTER = 600;

export interface HoldLoaderProps {
    onActionAfterHold: () => void;
    holdTime?: number;
    size?: number;
}

type UseHoldLoader = (props: HoldLoaderProps) => {
    CursorElement: FC;
    onMouseDown: MouseEventHandler<HTMLDivElement>,
    onMouseLeave: () => void,
    onMouseUp: () => void,
    onMouseMove: MouseEventHandler<HTMLDivElement>,
    onTouchStart: TouchEventHandler<HTMLDivElement>,
    onTouchMove: TouchEventHandler<HTMLDivElement>,
    onTouchEnd: () => void,
    onTouchCancel: () => void,
}

interface Cords {
    x: number;
    y: number;
}

export interface Options {
    pos?: Cords,
    screenPos?: Cords,
    afterHoldActionTimeoutId?: NodeJS.Timeout;
    onTouchStartTimeoutId?: NodeJS.Timeout;
}

export const useActionAfterHold: UseHoldLoader = ({
    holdTime = 2000,
    onActionAfterHold,
    size = LOADER_SIZE
}) => {
    const options = useRef<Options>({});
    const holdDiv = useRef<HTMLDivElement>(null);

    const onMouseDown: ReturnType<UseHoldLoader>['onMouseDown'] = (e) => {
        if (options.current.pos || !holdDiv.current || e.button !== 0) return;
        
        const rect = e.currentTarget.getBoundingClientRect();

        options.current.pos = {
            x: e.nativeEvent.clientX - rect.left,
            y: e.nativeEvent.clientY - rect.top,
        };

        options.current.afterHoldActionTimeoutId = setTimeout(() => {
            onActionAfterHold()
            options.current.afterHoldActionTimeoutId = undefined
        }, holdTime)

        holdDiv.current.toggleAttribute('data-helddown', true)
        holdDiv.current.style['left'] = (options.current.pos.x - size / 2) +'px';
        holdDiv.current.style['top'] = (options.current.pos.y - size / 2) +'px';
    }

    const onMouseLeave: ReturnType<UseHoldLoader>['onMouseLeave'] = () => {
        if (!options.current.pos || !holdDiv.current) return;

        options.current.pos = undefined;

        if (options.current.afterHoldActionTimeoutId) {
            clearTimeout(options.current.afterHoldActionTimeoutId)
        }

        options.current.afterHoldActionTimeoutId = undefined;

        holdDiv.current.toggleAttribute('data-helddown', false)
        holdDiv.current.style['left'] = '';
        holdDiv.current.style['top'] = '';
    }

    const onMouseMove: ReturnType<UseHoldLoader>['onMouseMove'] = (e) => {
        if (!options.current.pos) return;
        
        const rect = e.currentTarget.getBoundingClientRect();

        const movementSum = Math.abs(
            options.current.pos.x - e.nativeEvent.clientX - rect.left
        ) + Math.abs(
            options.current.pos.y - e.nativeEvent.clientY - rect.top
        )
        
        if (movementSum > MAX_MOVEMENT_LIMIT) {
            return onMouseLeave();
        }
    }

    const onTouchStart: ReturnType<UseHoldLoader>['onTouchStart'] = (e) => {
        if (options.current.screenPos || !holdDiv.current || e.nativeEvent.touches.length !== 1) return;

        const touch = e.touches[0];

        options.current.screenPos = {
            x: touch.clientX,
            y: touch.clientY,
        };

        const element = e.currentTarget;
        const rect = element.getBoundingClientRect();

        options.current.onTouchStartTimeoutId = setTimeout(() => {
            options.current.onTouchStartTimeoutId = undefined;

            options.current.afterHoldActionTimeoutId = setTimeout(() => {
                onActionAfterHold()
                options.current.afterHoldActionTimeoutId = undefined
            }, holdTime)
            
            const currentScreenPos = {
                x: touch.clientX - rect.left,
                y: touch.clientY - rect.top,
            };

            holdDiv.current!.toggleAttribute('data-helddown', true)
            holdDiv.current!.style['left'] = (currentScreenPos.x - size / 2) +'px';
            holdDiv.current!.style['top'] = (currentScreenPos.y - size / 2) +'px';
        }, START_TOUCH_EVENT_AFTER)
    }

    const onTouchMove: ReturnType<UseHoldLoader>['onTouchMove'] = (e) => {
        if (!options.current.screenPos) return;
        
        const touch = e.touches[0];

        const currentScreenPos = {
            x: touch.clientX,
            y: touch.clientY,
        };

        const movementSum = Math.abs(
            options.current.screenPos.x - currentScreenPos.x
        ) + Math.abs(
            options.current.screenPos.y - currentScreenPos.y
        )

        if (movementSum > MAX_MOVEMENT_LIMIT) {
            return onTouchEnd();
        }
    }

    const onTouchEnd: ReturnType<UseHoldLoader>['onTouchEnd'] = () => {
        if (!options.current.screenPos || !holdDiv.current) return;
        
        options.current.screenPos = undefined;

        if (options.current.onTouchStartTimeoutId) {
            clearTimeout(options.current.onTouchStartTimeoutId);
        }

        if (options.current.afterHoldActionTimeoutId) {
            clearTimeout(options.current.afterHoldActionTimeoutId)
        }

        options.current.onTouchStartTimeoutId = undefined;
        options.current.afterHoldActionTimeoutId = undefined;

        holdDiv.current.toggleAttribute('data-helddown', false)
        holdDiv.current.style['left'] = '';
        holdDiv.current.style['top'] = '';
    }

    return {
        onMouseDown,
        onMouseLeave,
        onMouseUp: onMouseLeave,
        onMouseMove,
        onTouchStart,
        onTouchMove,
        onTouchEnd,
        onTouchCancel: onTouchEnd,
        CursorElement: () => (
            <UserPressPoint ref={holdDiv} duration={`${holdTime}ms`} size={size}>
                <div className='loader'/>
            </UserPressPoint>
        ),
    }
}

export const onChildMouseDown: MouseEventHandler = (e) => {
    e.preventDefault();
    e.stopPropagation();
}

export const onContextMenu: MouseEventHandler = (e) => {
    e.preventDefault();
    e.stopPropagation();

    return false
}

// eslint-disable-next-line react-refresh/only-export-components
const UserPressPoint = forwardRef<unknown, ComponentProps<typeof Box> & { duration: string, size: number }>(({ sx, duration, size, ...props }, ref) => {
    return <Box
        ref={ref}
        sx={{
            ...{
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
                        },
                    },
                },
            },
            ...sx,
        }}
        {...props}
    />
});