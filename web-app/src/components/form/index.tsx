import { SxProps, Theme } from "@mui/material/styles";
import {
    Button as MuiButton,
    Input as MuiInput,
    Pagination as MuiPagination,
} from "@mui/material";
import { ComponentProps, FC, memo } from "react";

const InputStyle: SxProps<Theme> = {
    '&.MuiInputBase-root.MuiInput-root': {
        width: '100%',
        '&::before': {
            borderBottomColor: '#fff9',
        },
        '&::after': {
            borderBottomColor: 'crimson',
        },
        '& input': {
            color: 'var(--font-color)'
        },
    },
}

export const Input: FC<ComponentProps<typeof MuiInput>> = memo(({ sx, ...props }) => {
    return <MuiInput sx={{
        ...InputStyle,
        ...sx,
    }} {...props} />
})

export const Button = MuiButton;

export const Pagination = MuiPagination;

export { Searchbar } from './searcbar.js';