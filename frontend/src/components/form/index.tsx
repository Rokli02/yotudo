import { styled, CSSObject } from "@mui/material/styles";
import {
    FormControl as MuiFormControl,
    Button as MuiButton,
    Input as MuiInput,
    InputLabel as MuiInputLabel,
    Slider as MuiSlider,
    Autocomplete as MuiAutocomplete,
    TextField as MuiTextField,
    Paper as MuiPaper,
    Pagination as MuiPagination,
    Checkbox as MuiCheckbox,
    FormControlLabel as MuiFormControlLabel,
} from "@mui/material";
import { FC } from "react";

export const Form = styled('form')({
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    rowGap: '1rem'
})

export const FormControl = styled(MuiFormControl)({
    maxWidth: '500px',
    width: '100%',
    marginInline: 'auto',
    alignItems: 'start',
    rowGap: '.5rem'
})

export const FormControlLabel = styled(MuiFormControlLabel)({
    maxWidth: '500px',
    width: '100%',
    marginInline: 'auto',
    alignItems: 'center',
})

export const Input = styled(MuiInput)({
    '&.MuiInputBase-root.MuiInput-root': {
        width: '100%',
        '&::before': {
            borderBottomColor: '#fff9',
        } as CSSObject,
        '&::after': {
            borderBottomColor: 'crimson',
        } as CSSObject,
        '& input': {
            color: 'var(--font-color)'
        } as CSSObject,
    } as CSSObject,
})

export const InputLabel = styled(MuiInputLabel)({
    '&.MuiFormLabel-root.MuiInputLabel-root': {
        color: '#fffa',
        '&.Mui-focused': {
            color: 'inherit',
        } as CSSObject,
    } as CSSObject,
})

export const Button = styled(MuiButton)({
    width: 'max-content',
})

export const Slider = styled(MuiSlider)({
    color: 'var(--primary-color)'
})

const LocalAutocomplete = styled(MuiAutocomplete)({
    '&.MuiAutocomplete-root': {
        color: 'var(--font-color)',
        width: '100%',
        '& .MuiAutocomplete-endAdornment': {
            '& .MuiButtonBase-root': {
                color: 'var(--font-color)',
            },
        },
    },
})

export const Autocomplete: FC<Omit<Parameters<typeof LocalAutocomplete>[0], 'PaperComponent' | 'noOptionsText'>> = ({ ...props }) => <LocalAutocomplete PaperComponent={Paper} noOptionsText={<div style={{ color: 'var(--font-color)', textAlign: 'center'}}>Nincs szöveg</div>} {...props}/>

export const Paper = styled(MuiPaper)({
    '&.MuiPaper-root': {
        color: 'var(--font-color)',
        backgroundColor: 'var(--background-color)',
        border: '2px solid #fff1',
        borderRadius: '4px',
    },
})

export const TextField = styled(MuiTextField)({
    '&.MuiFormControl-root.MuiTextField-root': {
        width: '100%',
        '& .MuiFormLabel-root.MuiInputLabel-root': {
            left: 16,
            color: 'var(--font-color)',
        },
        '& .MuiInputBase-root': {
            color: 'var(--font-color)',
            '&.MuiInput-root': {
                ':before': {
                    borderColor: '#fff5',
                },
                ':after': {
                    borderColor: 'var(--primary-color)'
                }
            },
            '&.MuiOutlinedInput-root': {
                '& .MuiOutlinedInput-notchedOutline': {
                    borderColor: '#fff5',
                },
            },
            '&.MuiFilledInput-root': {
                ':before': {
                    borderColor: '#fff5',
                },
                ':after': {
                    borderColor: 'var(--primary-color)'
                }
            },
            '&.Mui-focused': {
                '&.MuiOutlinedInput-root .MuiOutlinedInput-notchedOutline': {
                    borderColor: 'var(--primary-color)',
                },
            },
        },
    } as CSSObject,
})

export const Pagination = styled(MuiPagination)({
    '&.MuiPagination-root': {
        '& .MuiPagination-ul': {
            justifyContent: 'center',
            alignItems: 'center',
            '& .MuiButtonBase-root': {
                color: 'var(--font-color)',
                ':hover': {
                    backgroundColor: '#ffffff0a',
                },
                '&.Mui-selected': {
                    backgroundColor: '#ffffff10',
                    ':hover': {
                        backgroundColor: '#ffffff1f',
                    }
                }
            }
        }
    }
})

export const Checkbox = styled(MuiCheckbox)({
    '&.MuiButtonBase-root.MuiCheckbox-root': {
        color: 'var(--font-color)',
        ':hover': {
            backgroundColor: '#aaa1',
        },
        '&.Mui-checked': {
            color: 'var(--primary-color)',
        }
    },
})

export { Searchbar } from './searcbar';
export { Select } from './select';
export { ImageSelector } from './image-selector/image-selector';
export type { Option } from './select';
export type { ImageSelectorProps } from './image-selector/image-selector';