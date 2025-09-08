import { ChangeEvent, ComponentProps, FC, useEffect, useState } from 'react'
import { Autocomplete, TextField } from '../../components/form'
import { useForm } from '.'
import { AutocompleteOptions } from './interface';
import { useTextFetchGuard } from '@src/hooks/useTextFetchGuard';
import { useLoading } from '@src/hooks/useLoading';

type MuiAutocompleteProps = Parameters<typeof Autocomplete>[0]
export interface AutocompleteProps extends Omit<MuiAutocompleteProps, 'defaultValue' | 'renderInput' | 'options'> {
    readonly name: string;
    readonly value?: AutocompleteOptions;
    readonly label?: string;
    readonly debounceTime?: number;
    readonly options?: AutocompleteOptions[];
    readonly getOptions: (search: string, abortController?: AbortController) => Promise<AutocompleteOptions[]>;
    readonly fetchOnce?: boolean;
    readonly onlyPreloadedOptions?: boolean;
}

export const FormAutocomplete: FC<AutocompleteProps> = ({
    name,
    value = null,
    label,
    onChange: onChangeArg,
    debounceTime = 300,
    freeSolo = true,
    options = [],
    fetchOnce = false,
    onlyPreloadedOptions = false,
    getOptions,
    ...props
}) => {
    const [textFieldValue, setTextFieldValue] = useState<string>(value?.label ?? '')
    const [selected, setSelected] = useState<AutocompleteOptions | null>(value)
    const [_options, setOptions] = useState<AutocompleteOptions[]>([...options])
    const { registerInput, unregisterInput, onValueChange, getErrors } = useForm();
    const fetchGuard = useTextFetchGuard()
    const loadingState = useLoading()

    const onTyping = (event: ChangeEvent<HTMLInputElement>) => {
        setTextFieldValue(event.target.value);
    }

    const onChange: AutocompleteProps['onChange'] = (event, value, reason, details) => {
        event.preventDefault();

        switch (reason) {
            case 'createOption':
                const trimedValue = (value as string).trim();
                const foundOption = _options.find((o) => o.label.toLowerCase().search(trimedValue.toLowerCase()) !== -1);
                value = foundOption
                    ? foundOption
                    : onlyPreloadedOptions
                    ? null
                    : { name: trimedValue, label: trimedValue } satisfies AutocompleteOptions;
            case 'clear':
            case 'selectOption':
                if (!loadingState.value) {
                    setSelected(value as AutocompleteOptions | null);
                    onValueChange(name, value);
                    onChangeArg?.(event, value, reason, details);
                }
            break;
        }
    }

    const onClear = () => {
        setTextFieldValue(value?.label ?? '');
        setSelected(value);
        fetchGuard.makeItWorthFetching();
    }

    useEffect(() => {
      registerInput(name, setSelected, onClear, value);

      if (fetchOnce) {
        getOptions(textFieldValue).then((values) => setOptions(values))
      }

      return () => {
        unregisterInput(name);
      }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])
    
    useEffect(() => {
        if (!fetchOnce && !options?.length) {
            const abortController = new AbortController();
            let timeoutId: NodeJS.Timeout;

            if (fetchGuard.worthFetching(textFieldValue)) {
                loadingState.startLoading();

                timeoutId = setTimeout(async () => {
                    const fetchedOptions = await getOptions(textFieldValue, abortController);
    
                    fetchGuard.setShouldFetch(!!fetchedOptions.length);
    
                    setOptions(fetchedOptions)
                    loadingState.stopLoading();
                }, debounceTime);
            }

            return () => {
                clearTimeout(timeoutId)
                abortController.abort();
                loadingState.stopLoading();
            }
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [textFieldValue, debounceTime])

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
                        value={textFieldValue}
                    />)
                }
                freeSolo={freeSolo}
                value={selected}
                options={_options}
                slotProps={slotProps}
                onChange={onChange}
            />
            { getErrors(name) }
        </>
    )
}

const slotProps: ComponentProps<typeof Autocomplete>['slotProps'] = { listbox: { sx: { maxHeight: '250px' } } }