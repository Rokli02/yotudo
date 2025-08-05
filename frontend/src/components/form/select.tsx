/* eslint-disable @typescript-eslint/no-explicit-any */
import { BaseSelectProps, MenuItem as MuiMenuItem, Select as MuiSelect } from "@mui/material";
import { styled, CSSObject } from "@mui/material/styles";
import { FC, memo, ReactNode, useMemo } from "react";

type MuiSelectProps = Parameters<typeof MuiSelect>[0];

export interface SelectProps extends Omit<MuiSelectProps, 'onChange' | 'children'> {
    options: Option[];
    onChange?: (value: any, props: Props) => void;
}

export interface Option<T = unknown> {
    value: T;
    label: ReactNode;
}

export interface Props {
    value: any;
    children: unknown;
}

export const Select: FC<SelectProps> = memo(({ options, onChange, ...props }) => {

    const _onChange: MuiSelectProps['onChange'] = !onChange ? undefined : (ev, child) => {
        onChange(ev.target.value, (child as any)?.props ?? { children: null, value: ev.target.value })
    }

    const _options = useMemo(() => {
        return options.map((option, i) => <StyledOption key={`${i}_${option.value}`} value={option.value as string}>{option.label}</StyledOption>)
    }, [options])

    return (
        <StyledSelect onChange={_onChange} MenuProps={menuProps} variant="outlined" {...props}>
            { _options }
        </StyledSelect>
    )
})

const menuProps: BaseSelectProps['MenuProps'] = { slotProps: { paper: { sx: { backgroundColor: 'var(--background-color)', color: 'var(--font-color)' } } } }

const StyledSelect = styled(MuiSelect)({
    '&.MuiInputBase-root': {
        color: 'var(--font-color)',
        '.MuiSvgIcon-root': {
            color: '#ffffff66',
            transition: 'transform 250ms',
        },
        '.MuiOutlinedInput-notchedOutline, :before': {
            borderColor: '#ffffff3a',
        },
        ':after': {
            borderColor: 'var(--primary-color)',
        },
    } as CSSObject,
});

const StyledOption = styled(MuiMenuItem)({
    '&.MuiButtonBase-root': {
        ':hover': {
            backgroundColor: '#c8c8e00c',
        }
    },
});