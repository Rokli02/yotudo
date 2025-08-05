import { FC, useEffect, useState } from 'react'
import { Slider } from '../../components/form'
import { useForm } from '.'
import { SliderOwnProps } from '@mui/material'

type MuiSliderProps =  Parameters<typeof Slider>[0]
interface SliderProps extends MuiSliderProps {
    name: string,
    value?: number,
    onChange?: <T = string>(value: T) => void,
}

export const FormSlider: FC<SliderProps> = ({ onChange, name, value = 0, valueLabelDisplay = 'auto', ...props }) => {
    const [_value, setValue] = useState<unknown>(value)
    const { registerInput, unregisterInput, onValueChange, getErrors } = useForm();
    const _onChange: SliderOwnProps<number | number[]>['onChange'] = (_, value) => {
        setValue(value);
        onValueChange(name, value);
        onChange?.(value)
    }

    const onClear = () => {
        setValue(value);
    }
    
    useEffect(() => {
        registerInput(name, setValue, onClear, value);

        return () => {
            unregisterInput(name);
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    return (
        <>
            <Slider name={name} onChange={_onChange} value={_value as number} valueLabelDisplay={valueLabelDisplay} {...props} />
            { getErrors(name) }
        </>
    )
}
