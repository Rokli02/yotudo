import { ChangeEvent, FC, ReactNode, useEffect, useState } from 'react'
import { Autocomplete, TextField } from '../../components/form'
import { useForm } from '.'
import { AutocompleteOptions } from './interface';
import { styled } from '@mui/material/styles';
import { Chip } from '../../components/common';

type MuiAutocompleteProps = Parameters<typeof Autocomplete>[0]
export interface MultiselectAutocompleteProps extends Omit<MuiAutocompleteProps, 'defaultValue' | 'renderInput' | 'options' | 'onChange'> {
    name: string;
    value?: string;
    label?: string;
    debounceTime?: number;
    selectedOptions?: AutocompleteOptions[];
    readonly options?: AutocompleteOptions[];
    onChange?: (event: React.SyntheticEvent, value: AutocompleteOptions, values: AutocompleteOptions[]) => void
    getOptions: (search: string, selectedOptionIds: number[], abortController?: AbortController) => Promise<AutocompleteOptions[]>;
    renderChipContent?: (value: AutocompleteOptions, index: number, array: AutocompleteOptions[]) => ReactNode;
    fetchOnce?: boolean;
}

export const FormMultiselectAutocomplete: FC<MultiselectAutocompleteProps> = ({
    name,
    value: valueProp = '',
    label,
    onChange,
    debounceTime = 750,
    freeSolo = true,
    options = [],
    selectedOptions = [],
    fetchOnce = false,
    getOptions,
    renderChipContent = (v) => `${v.id} - ${v.label}`,
    ...props
}) => {
    const [_value, setValue] = useState<string>(valueProp);
    const [_selectedOptions, setSelectedOptions] = useState<AutocompleteOptions[]>(selectedOptions);
    const [_options, setOptions] = useState<AutocompleteOptions[]>([...options]);
    const { registerInput, unregisterInput, onValueChange, getErrors } = useForm();

    const onTyping = (event: ChangeEvent<HTMLInputElement>) => {
        setValue(event.target.value);
    }

    const _onChange: MuiAutocompleteProps['onChange'] = (event, value) => {
        if (typeof value !== 'object' || !value || !(value as Record<string, unknown>)['id']) {
            console.error('_onChange invalid input value', value)
            return
        }

        setValue(valueProp);
        setOptions((pre) => pre.filter((option) => option.id != (value as Record<string, unknown>)['id']))
        setSelectedOptions((pre) => {
            const newSelectedOptions = [...pre, (value as AutocompleteOptions)];
            
            onChange?.(event, (value as AutocompleteOptions), newSelectedOptions);

            return newSelectedOptions;
        })
    }

    const onClear = () => {
        setValue(valueProp);
        setSelectedOptions(selectedOptions)
    }

    const onChipClose = (value: AutocompleteOptions, index: number) => () => {
        setOptions((pre) => [value, ...pre])
        setSelectedOptions((pre) => {
            // Remove one element from selected options
            if (pre[index] && pre[index]['id'] === value['id']) {
                pre.splice(index, 1);

                return [...pre]
            }

            return pre.filter((so) => so['id'] !== value['id'])
        })
    }

    useEffect(() => {
      registerInput(name, setSelectedOptions as React.Dispatch<React.SetStateAction<unknown>>, onClear, selectedOptions);

      if (fetchOnce) {
        getOptions(_value, selectedOptions.map((so) => so['id'] as number)).then((values) => setOptions(values))
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
                setOptions(await getOptions(_value, _selectedOptions.map((so) => (so as Record<string, number>)['id']), abortController))
            }, debounceTime)
            
            return () => {
                clearTimeout(timeoutId)
                abortController.abort();
            }
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [_value, debounceTime])

    useEffect(() => {
        onValueChange(name, _selectedOptions);
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [_selectedOptions])

    return (
        <>
            <Autocomplete
                {...props}
                renderInput={(params) => (<TextField variant='standard' name={name} {...params} label={label} onChange={onTyping}/>)}
                freeSolo={freeSolo}
                inputValue={_value}
                options={_options}
                ListboxProps={{ sx: { maxHeight: '250px' } }}
                onChange={_onChange}
            />
            <ChipContainer>
                { _selectedOptions.map((so, index, arr) => {
                    return <Chip key={`${index}_${so.id}`} onClose={onChipClose(so, index)}>{ renderChipContent(so, index, arr) }</Chip>
                }) }
            </ChipContainer>
            { getErrors(name) }
        </>
    )
}


const ChipContainer = styled('div')({
    position: 'relative',
    display: 'flex',
    flexWrap: 'wrap',
    flexDirection: 'row',
    gap: '.5rem .75rem',
    marginTop: '4px',
    marginInline: '8px',
});