import { CSSProperties, FC, useState, ChangeEventHandler, useRef } from 'react'
import { CSSObject, styled } from '@mui/material/styles';
import { Search } from '@mui/icons-material';
import { Input } from '.'

export interface SearchbarProps {
    className?: string;
    style?: CSSProperties;
    debounceTime?: number;
    onDebounce: (value: string) => void;
}

export const Searchbar: FC<SearchbarProps> = ({ debounceTime = 500, onDebounce, className, style, ...props }) => {
    const [search, setSearch] = useState<string>('')
    const timeoutIdRef = useRef<any>(undefined)

    const onChange: ChangeEventHandler<HTMLInputElement | HTMLTextAreaElement> = (event) => {
        setSearch(event.target.value);

        if (timeoutIdRef.current != null) {
            clearTimeout(timeoutIdRef.current);
        }

        const timeoutId = setTimeout(() => {
            onDebounce(search)
            timeoutIdRef.current = undefined;
        }, debounceTime)

        timeoutIdRef.current = timeoutId;
    }

    return (
        <SearchbarContainer className={className} style={style}>
            <Input
                name='search'
                className='SearchBar_Input'
                inputProps={{
                    color: 'red'
                }}
                placeholder='Keresés...'
                type='text'
                autoComplete='off'
                onChange={onChange}
                value={search}
                renderSuffix={
                    ({ focused }) => <SearchSuffixContainer className={ focused ? 'focused' : undefined }><Search /></SearchSuffixContainer>
                }
                { ...props }
            />
        </SearchbarContainer>
    )
}

const SearchSuffixContainer = styled('div')({
    width: '24px',
    height: '24px',
    position: 'relative',
    marginInline: '6px 12px',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    pointerEvents: 'none',
    color: 'var(--font-color)',
    right: 0,
    opacity: 1,
    transition: `
        right 150ms,
        opacity 125ms,
        marginInline 150ms,
        width 200ms
    `,
    '&.focused': {
        right: -12,
        opacity: 0,
        marginInline: 0,
        width: 0,
    } as CSSObject,
})

const SearchbarContainer = styled('div')({
    position: 'relative',
    '& .SearchBar_Input': {
        '&.MuiInputBase-root': {
            '& .MuiInputBase-input': {},
        },
        lineHeight: 48,
        height: 48,
    } as CSSObject,
});
