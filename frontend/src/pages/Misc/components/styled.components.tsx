import { Box } from "@mui/material";
import { ComponentProps, FC, memo } from "react";

export const Container: FC<ComponentProps<typeof Box>> = memo(({ sx, ...props}) => {
    return <Box
        sx={{
            ...{
                position: 'relative',
                backgroundColor: 'var(--primary-color)',
                minWidth: '490px',
                width: '100%',
                maxWidth: '700px',
                padding: '8px 12px',
                borderRadius: '8px',
                boxShadow: '3px 3px 6px #0004',
                color: 'inherit',
            },
            ...sx
        }}
        {...props}
    />
})

export const ItemContainer: FC<ComponentProps<typeof Box>> = memo(({ sx, ...props}) => {
    return <Box
        sx={{
            ...{
                width: '100%',
                display: 'flex',
                flexWrap: 'nowrap',
                backgroundColor: '#0003',
                height: 'fit-content',
                padding: '10px 16px',
                borderRadius: '8px',
                columnGap: '12px',
                overflowX: 'clip',
                '&.min_width': {
                    width: 'max-content',
                },
            },
            ...sx
        }}
        {...props}
    />
})

export const Content: FC<ComponentProps<typeof Box>> = memo(({ sx, ...props}) => {
    return <Box
        sx={{
            ...{
                display: 'flex',
                rowGap: '10px',
                columnGap: '8px',
                alignItems: 'center',
                flexGrow: 1,
                '&[data-dir=row]': {
                    flexDirection: 'row',
                    flexWrap: 'wrap',
                },
                '&[data-dir=col]': {
                    flexDirection: 'column',
                },
            },
            ...sx
        }}
        {...props}
    />
})