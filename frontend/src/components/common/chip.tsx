import { Close } from "@mui/icons-material";
import { Box } from "@mui/material";
import { SxProps, Theme } from "@mui/material/styles";
import { CSSProperties, memo, ReactNode } from "react";

export interface ChipProps {
    children: ReactNode;
    style?: CSSProperties;
    className?: string;
    onClose?: () => void;
}
export const Chip = memo(({ children, onClose, ...props }: ChipProps) => {

    return (
        <Box sx={ChipBodyStyle} {...props}>
            <Box sx={ChipContentStyle}>
                { children }
            </Box>
            { !onClose ? undefined : <Box sx={IconButtonStyle} onClick={onClose}><Close style={{ width: '85%', height: '85%' }} /></Box> }
        </Box>
    )
})

const ChipBodyStyle: SxProps<Theme> = {
    position: 'relative',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    columnGap: '.5rem',
    backgroundColor: 'var(--primary-color)',
    paddingRight: '3px',
    borderRadius: '8px',
    width: 'max-content',
};

const ChipContentStyle: SxProps<Theme> = {
    padding: '4px 10px 4px 10px',
};

const IconButtonStyle: SxProps<Theme> = {
    cursor: 'pointer',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    width: '24px',
    height: '24px',
    borderRadius: '50%',
    ':hover': {
        backgroundColor: '#3333'
    },
    ':active': {
        backgroundColor: '#1113'
    },
};
