import { Close } from "@mui/icons-material";
import { CSSObject, styled } from "@mui/material/styles";
import { CSSProperties, memo, ReactNode } from "react";

export interface ChipProps {
    children: ReactNode;
    style?: CSSProperties;
    className?: string;
    onClose?: () => void;
}
export const Chip = memo(({ children, onClose, ...props }: ChipProps) => {

    return (
        <ChipBody {...props}>
            <ChipContent>
                { children }
            </ChipContent>
            { !onClose ? undefined : <IconButton onClick={onClose}><Close style={{ width: '85%', height: '85%' }} /></IconButton> }
        </ChipBody>
    )
})

const ChipBody = styled('div')({
    position: 'relative',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    columnGap: '.5rem',
    backgroundColor: 'var(--primary-color)',
    paddingRight: '3px',
    borderRadius: '8px',
    width: 'max-content',
})

const ChipContent = styled('div')({
    padding: '4px 10px 4px 10px',
})

const IconButton = styled('div')({
    cursor: 'pointer',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    width: '24px',
    height: '24px',
    borderRadius: '50%',
    ':hover': {
        backgroundColor: '#3333'
    } as CSSObject,
    ':active': {
        backgroundColor: '#1113'
    } as CSSObject,
})
