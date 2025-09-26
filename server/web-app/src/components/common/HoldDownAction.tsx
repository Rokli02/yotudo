import { Box } from "@mui/material";
import { useActionAfterHold, onContextMenu } from "@src/hooks/useHoldLoader.js";
import { ComponentProps, ElementType, FC, memo } from "react";

// eslint-disable-next-line react-refresh/only-export-components
export { onChildMouseDown } from "@src/hooks/useHoldLoader.js";

export interface HoldDownActionProps extends ComponentProps<typeof Box> {
    onActionAfterHold: () => void;
    holdTime?: number;
    size?: number;
    component?: ElementType,
} 

export const HoldDownAction: FC<HoldDownActionProps> = memo(({ onActionAfterHold, holdTime, size, children, component='div', ...props }) => {
    const {
        CursorElement,
        onMouseDown,
        onMouseMove,
        onMouseUp,
        onMouseLeave,
        onTouchStart,
        onTouchMove,
        onTouchEnd,
        onTouchCancel,
    } = useActionAfterHold({ onActionAfterHold, holdTime, size })

    return <Box
        component={component}
        onMouseDown={onMouseDown}
        onMouseMove={onMouseMove}
        onMouseUp={onMouseUp}
        onMouseLeave={onMouseLeave}
        onTouchStart={onTouchStart}
        onTouchMove={onTouchMove}
        onTouchEnd={onTouchEnd}
        onTouchCancel={onTouchCancel}
        onContextMenu={onContextMenu}
        {...props}
    >
        <CursorElement />
        {children}
    </Box>
})