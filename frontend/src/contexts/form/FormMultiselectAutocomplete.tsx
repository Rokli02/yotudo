import { ChangeEvent, ComponentProps, FC, ReactNode, useEffect, useState } from 'react'
import { Autocomplete, TextField } from '../../components/form'
import { useForm } from '.'
import { AutocompleteOptions } from './interface';
import { styled } from '@mui/material/styles';
import { Chip } from '../../components/common';
import { useTextFetchGuard } from '@src/hooks/useTextFetchGuard';

type MuiAutocompleteProps = Parameters<typeof Autocomplete>[0]
export interface MultiselectAutocompleteProps extends Omit<MuiAutocompleteProps, 'defaultValue' | 'renderInput' | 'options' | 'onChange' | 'value'> {
    readonly name: string;
    readonly label?: string;
    readonly debounceTime?: number;
    readonly selectedOptions?: AutocompleteOptions[];
    readonly options?: AutocompleteOptions[];
    readonly onChange?: (event: React.SyntheticEvent, value: AutocompleteOptions, values: AutocompleteOptions[]) => void
    readonly getOptions: (search: string, selectedOptionIds: number[], abortController?: AbortController) => Promise<AutocompleteOptions[]>;
    readonly renderChipContent?: (value: AutocompleteOptions, index: number, array: AutocompleteOptions[]) => ReactNode;
    readonly fetchOnce?: boolean;
    readonly onlyPreloadedOptions?: boolean;
}

export const FormMultiselectAutocomplete: FC<MultiselectAutocompleteProps> = ({
    name,
    label,
    onChange,
    debounceTime = 300,
    freeSolo = true,
    options = [],
    selectedOptions = [],
    fetchOnce = false,
    onlyPreloadedOptions = false,
    getOptions,
    renderChipContent = (v) => `${v.id} - ${v.label}`,
    ...props
}) => {
    const [inputValue, setInputValue] = useState<string>('');
    const [_selectedOptions, setSelectedOptions] = useState<AutocompleteOptions[]>(selectedOptions);
    const [_options, setOptions] = useState<AutocompleteOptions[]>([...options]);
    const { registerInput, unregisterInput, onValueChange, getErrors } = useForm();
    const fetchGuard = useTextFetchGuard()

    const onTyping = (event: ChangeEvent<HTMLInputElement>) => {
        setInputValue(event.target.value);
    }

    const _onChange: MuiAutocompleteProps['onChange'] = (event, value, reason) => {
        event.preventDefault();

        switch (reason) {
            case 'clear':
                setInputValue('')
            break;
            case 'selectOption':
                setInputValue('');
                setOptions((pre) => pre.filter((option) => option.id != (value as AutocompleteOptions).id))
                setSelectedOptions((pre) => {
                    const newSelectedOptions = [...pre, (value as AutocompleteOptions)];
                    
                    onChange?.(event, (value as AutocompleteOptions), newSelectedOptions);

                    onValueChange(name, newSelectedOptions);
                    return newSelectedOptions;
                })
            break;
            case 'createOption':
                const trimedValue = (value as string).trim();
                const foundOption = _options.find((o) => o.label.toLowerCase().search(trimedValue.toLowerCase()) !== -1);

                if (foundOption) {
                    setInputValue('');
                    setOptions((pre) => pre.filter((option) => option.id != foundOption.id))
                    setSelectedOptions((pre) => {
                        const newSelectedOptions = [...pre, foundOption];
                        
                        onChange?.(event, foundOption, newSelectedOptions);

                        onValueChange(name, newSelectedOptions);
                        return newSelectedOptions;
                    })
                } else if (!onlyPreloadedOptions) {
                    if (_selectedOptions.find((p) => p.label.toLowerCase() === trimedValue.toLowerCase())) break;

                    setInputValue('');
                    setSelectedOptions((pre) => {
                        const newOption = { name: trimedValue, label: trimedValue } satisfies AutocompleteOptions;
                        const newSelectedOptions = [...pre, newOption];
                        
                        onChange?.(event, newOption, newSelectedOptions);

                        onValueChange(name, newSelectedOptions);
                        return newSelectedOptions;
                    })
                }
        }
    }

    const onClear = () => {
        setInputValue('');
        setSelectedOptions(selectedOptions)
        onValueChange(name, selectedOptions);
        fetchGuard.makeItWorthFetching();
    }

    const onChipClose = (value: AutocompleteOptions, index: number) => () => {
        if (value.id !== undefined) {
            setOptions((pre) => [value, ...pre])
        }

        setSelectedOptions((pre) => {
            if (pre[index] && pre[index].id === value.id) {
                pre.splice(index, 1);

                const newSelectedOptions = [...pre];
                onValueChange(name, newSelectedOptions);
                return newSelectedOptions
            }

            const newSelectedOptions = pre.filter((so) => so.id !== value.id);
            onValueChange(name, newSelectedOptions);
            return newSelectedOptions
        })
    }

    useEffect(() => {
      registerInput(name, setSelectedOptions, onClear, selectedOptions);

      if (fetchOnce) {
        const filteredOptions: number[] = selectedOptions.map((so) => so.id as number).filter((id) => id !== undefined);

        getOptions(
            inputValue,
            filteredOptions,
        ).then((values) => setOptions(values))
      }

      return () => {
        unregisterInput(name);
      }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])
    
    useEffect(() => {
        if ((!options || options.length === 0) && !fetchOnce) {
            const abortController = new AbortController();

            let timeoutId: NodeJS.Timeout
            if (fetchGuard.worthFetching(inputValue)) {
                timeoutId = setTimeout(async () => {
                    const filteredOptions: number[] = _selectedOptions.map((so) => so.id as number).filter((id) => id !== undefined)
    
                    const fetchedOptions = await getOptions(
                        inputValue,
                        filteredOptions,
                        abortController,
                    );

                    if (!fetchedOptions.length) fetchGuard.worthFetching(inputValue, false)

                    setOptions(fetchedOptions);
                }, debounceTime)
            }
            
            return () => {
                clearTimeout(timeoutId)
                abortController.abort();
            }
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [inputValue, debounceTime])

    return (
        <>
            <Autocomplete
                {...props}
                renderInput={(params) => (
                    <TextField
                        variant='standard'
                        name={name}
                        {...params}
                        label={label}
                        onChange={onTyping}
                    />
                )}
                freeSolo={freeSolo}
                inputValue={inputValue}
                options={_options}
                slotProps={slotProps}
                onChange={_onChange}
                onKeyDown={(e) => { e.key === 'Enter' && e.preventDefault() }}
            />
            <ChipContainer>
                { _selectedOptions.map((so, index, arr) => (
                    <Chip key={`${index}_${so.id}`} onClose={onChipClose(so, index)}>
                        { renderChipContent(so, index, arr) }
                    </Chip>
                )) }
            </ChipContainer>
            { getErrors(name) }
        </>
    )
}

const slotProps: ComponentProps<typeof Autocomplete>['slotProps'] = { listbox: { sx: { maxHeight: '250px' } }}

const ChipContainer = styled('div')({
    position: 'relative',
    display: 'flex',
    flexWrap: 'wrap',
    flexDirection: 'row',
    gap: '.5rem .75rem',
    marginTop: '4px',
    marginInline: '8px',
});