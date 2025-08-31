import { ImageSelector } from "@src/components/form";
import { ComponentProps, FC, useEffect, useState } from "react";
import { useForm } from ".";
import { Box } from "@mui/material";
import { SelectorType } from "@src/components/form/image-selector/select-image";

interface FormImageSelector extends Omit<ComponentProps<typeof Box>, 'onChange'>{
    name: string,
    defaultValue?: string;
    onChange?: (v?: string | null, type?: SelectorType) => void;
}

export const FormImageSelector: FC<FormImageSelector> = ({ name, defaultValue, onChange, ...props }) => {
    const [_value, setValue] = useState<unknown>(defaultValue)
    const { registerInput, unregisterInput, onValueChange, getErrors } = useForm();

    const _onChange: ComponentProps<typeof ImageSelector>['onChange'] = (type, v) => {
        setValue(v);
        onValueChange(name, v);
        onValueChange(`${name}_chosenType`, type);
        onChange?.(v, type)
    }

    const onClear = () => {
        setValue(defaultValue);
        onValueChange(`${name}_chosenType`, undefined);

    }

    useEffect(() => {
        registerInput(name, setValue, onClear, defaultValue);
        registerInput(`${name}_chosenType`, () => {}, () => {}, undefined);

        return () => {
            unregisterInput(name);
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    return <>
        <ImageSelector value={_value as string | undefined} onChange={_onChange} restoreValue={onClear} {...props} />
        { getErrors(name) }
    </>
}
