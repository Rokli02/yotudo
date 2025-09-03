import { CSSProperties, FC, useEffect, useState } from 'react'
import { Search } from '@mui/icons-material';
import { Input } from '.'
import { Box } from '@mui/material';

export interface SearchbarProps {
    className?: string;
    style?: CSSProperties;
    debounceTime?: number;
    onDebounce: (value: string) => void;
}

export const Searchbar: FC<SearchbarProps> = ({ debounceTime = 500, onDebounce, className, style, ...props }) => {
    const [search, setSearch] = useState<string>('')
    const [firstRender, setFirstRender] = useState<boolean>(true)

    useEffect(() => {
        if (firstRender) return setFirstRender(false);

        const timeoutId = setTimeout(() => {
            onDebounce(search)
        }, debounceTime)
        
        return () => {
            clearTimeout(timeoutId)
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [search, debounceTime])
    

    return (
        <Box sx={SearchbarContainerStyle} className={className} style={style}>
            <Input
                name='search'
                className='SearchBar_Input'
                inputProps={{
                    color: 'red'
                }}
                placeholder='Keresés...'
                type='text'
                autoComplete='off'
                onChange={(event) => setSearch(event.target.value)}
                value={search}
                renderSuffix={
                    ({ focused }) => <Box sx={SearchSuffixContainerStyle} className={ focused ? 'focused' : undefined }><Search /></Box>
                }
                { ...props }
            />
        </Box>
    )
}

const SearchSuffixContainerStyle = {
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
    },
};

const SearchbarContainerStyle = {
    position: 'relative',
    '& .SearchBar_Input': {
        '&.MuiInputBase-root': {
            '& .MuiInputBase-input': {},
        },
        lineHeight: 48,
        height: 48,
    },
};