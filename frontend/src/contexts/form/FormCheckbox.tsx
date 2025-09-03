import { Dispatch, FC, useEffect, useState } from 'react'
import { Checkbox, FormControlLabel } from '../../components/form'
import { useForm } from '.';
import { FormControlLabelProps } from '@mui/material';

type MuiCheckboxProps = Parameters<typeof Checkbox>[0];
interface FormCheckboxProps extends Omit<MuiCheckboxProps, ''> {
    label: string;
    name: string;
    value?: boolean;
    onChange?: FormControlLabelProps['onChange']
}

export const FormCheckbox: FC<FormCheckboxProps> = ({ label, name, onChange, value = false, ...props }) => {
    const [_value, setValue] = useState<boolean>(value);
    const { registerInput, unregisterInput, onValueChange, getErrors } = useForm();

    const _onChange: FormControlLabelProps['onChange'] = (e, checked) => {
        setValue(checked);
        onValueChange(name, checked);
        onChange?.(e, checked)
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
            <FormControlLabel
                label={label}
                name={name}
                checked={_value}
                onChange={_onChange}
                control={<Checkbox {...props} />}
            />
            { getErrors(name) }
        </>
    )
}

export default FormCheckbox