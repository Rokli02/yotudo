import { Box } from "@mui/material";
import { ComponentProps, CSSProperties, FC, memo } from "react";

export const Divider: FC<ComponentProps<typeof Box> & { dir?: 'horizontal' | 'vertical', length?: string | number }> = memo(({ dir, length = '100%', sx, ...props }) => {
    const resultObj: CSSProperties = {
        display: 'inline-flex',
        backgroundColor: '#fff6',
        border: 'none',
        borderRadius: '3px',
        padding: '-6px -12px'
    }

    if (dir === 'horizontal') {
        resultObj.width = length;
        resultObj.height = '2px';
        resultObj.marginInline = 'auto';
    } else {
        resultObj.width = '2px';
        resultObj.height = length;
        resultObj.marginBlock = 'auto';
    }

    return <Box
        component={'hr'}
        sx={{
            ...resultObj,
            ...sx,
        }}
        {...props}
    />
})