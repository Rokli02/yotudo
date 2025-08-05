import { ChangeEventHandler, FC, useEffect, useState } from 'react'
import { Input } from '../../components/form'
import { useForm } from '.'

type MuiInputProps =  Parameters<typeof Input>[0]
interface InputProps extends MuiInputProps {
    name: string,
    onChange?: <T = string>(value: T) => void,
}

export const FormInput: FC<Omit<InputProps, 'defaultValue'>> = ({ onChange, name, value = '', ...props }) => {
    const [_value, setValue] = useState<unknown>(value)
    const { registerInput, unregisterInput, onValueChange, getErrors } = useForm();

    const _onChange: ChangeEventHandler<HTMLTextAreaElement | HTMLInputElement> = (event) => {
        setValue(event.target.value);
        onValueChange(name, event.target.value);
        onChange?.(event.target.value)
    }

    const onClear = () => {
        setValue(value);
    }

    useEffect(() => {
        registerInput(name, setValue, onClear, value);

        return () => {
            unregisterInput;
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    return (
        <>
            <Input onChange={_onChange} value={_value} name={name} {...props} />
            { getErrors(name) }
        </>
    )
}

export default FormInput;
