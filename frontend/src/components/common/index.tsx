import { SvgIconComponent } from "@mui/icons-material";
import { Box } from "@mui/material";
import { ComponentProps, CSSProperties, FC, memo } from "react";

export const Title: FC<ComponentProps<typeof Box>> = memo(({ sx, ...props }) => {
    return <Box
        sx={{
            ...{
                fontSize: '1.25rem',
                marginInline: 'auto',
                width: 'max-content',
            },
            ...sx,
        }}  
        {...props}
    />
})

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

const TRCBWrapperStyle = {
    cursor: 'pointer',
    position: 'absolute',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    top: '4px',
    right: '6px',
    width: '28px',
    height: '28px',
    borderRadius: '50%',
    ':hover': {
        backgroundColor: '#3333'
    },
    ':active': {
        backgroundColor: '#1113'
    },
};

export const TopRightButton = memo(({ style = {}, Icon, onClick }: { style?: CSSProperties, Icon: SvgIconComponent, onClick: () => void }) => {
    return (
        <Box sx={TRCBWrapperStyle} onClick={onClick} style={style}>
            <Icon style={{
                width: '85%',
                height: '85%',
            }}/>
        </Box>
    )
})

export { Chip } from './chip';
export { StatusIcon, StatusActionIcon } from './status';