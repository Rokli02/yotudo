import { ChangeEvent, Dispatch, FC, useEffect, useState } from 'react'
import { Autocomplete, TextField } from '../../components/form'
import { useForm } from '.'
import { AutocompleteOptions } from './interface';

type MuiAutocompleteProps = Parameters<typeof Autocomplete>[0]
export interface AutocompleteProps extends Omit<MuiAutocompleteProps, 'defaultValue' | 'renderInput' | 'options'> {
    name: string;
    value?: AutocompleteOptions;
    label?: string;
    debounceTime?: number;
    readonly options?: AutocompleteOptions[];
    getOptions: (search: string, abortController?: AbortController) => Promise<AutocompleteOptions[]>;
    fetchOnce?: boolean;
}
//FIXME: Az 'onClear' nem nullázza le a mezőt!
export const FormAutocomplete: FC<AutocompleteProps> = ({
    name,
    value = null,
    label,
    onChange,
    debounceTime = 750,
    freeSolo = true,
    options = [],
    fetchOnce = false,
    getOptions,
    ...props
}) => {
    const [_value, setValue] = useState<string>(value?.label ?? '')
    const [selected, setSelected] = useState<AutocompleteOptions | null>(value)
    const [_options, setOptions] = useState<AutocompleteOptions[]>([...options])
    const { registerInput, unregisterInput, onValueChange, getErrors } = useForm();

    const onTyping = (event: ChangeEvent<HTMLInputElement>) => {
        setValue(event.target.value);
    }

    const _onChange: AutocompleteProps['onChange'] = (event, value, reason, details) => {
        if (typeof value !== 'object' || !value || !(value as Record<string, unknown>)['id']) {
            return
        }

        setSelected(value as AutocompleteOptions);
        onValueChange(name, value);
        onChange?.(event, value, reason, details);
    }

    const onClear = () => {
        setValue('');
        setSelected(value);
    }

    useEffect(() => {
      registerInput(name, setSelected as Dispatch<React.SetStateAction<unknown>>, onClear, value);

      if (fetchOnce) {
        getOptions(_value).then((values) => setOptions(values))
      }

      return () => {
        unregisterInput(name);
      }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])
    
    useEffect(() => {
        if ((!options || options.length === 0) && !fetchOnce) {
            const abortController = new AbortController();
            const timeoutId = setTimeout(async () => {
                setOptions(await getOptions(_value, abortController))
            }, debounceTime)
            
            return () => {
                clearTimeout(timeoutId)
                abortController.abort();
            }
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [_value, debounceTime])

    return (
        <>
            <Autocomplete
                {...props}
                renderInput={(params) => (<TextField variant='standard' name={name} {...params} label={label} onChange={onTyping} value={_value}/>)}
                freeSolo={freeSolo}
                value={selected}
                options={_options}
                slotProps={{ listbox: { sx: { maxHeight: '250px' } } }}
                onChange={_onChange}
            />
            { getErrors(name) }
        </>
    )
}
