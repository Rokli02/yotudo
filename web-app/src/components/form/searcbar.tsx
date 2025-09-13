import { CSSProperties, FC, useEffect, useState, memo } from 'react'
import { Search } from '@mui/icons-material';
import { Input } from './index.js'
import { Box, SxProps, Theme } from '@mui/material';

export interface SearchbarProps {
    value?: string;
    className?: string;
    style?: CSSProperties;
    sx?: SxProps<Theme>;
    debounceTime?: number;
    onDebounce: (value: string) => void;
}

export const Searchbar: FC<SearchbarProps> = memo(({
    value = '',
    debounceTime = 500,
    onDebounce,
    className,
    style,
    sx,
    ...props
}) => {
    const [search, setSearch] = useState<string>(value)
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

    const _sx = {
        ...SearchbarContainerStyle,
        ...sx,
    } as SxProps<Theme>

    return (
        <Box sx={_sx} className={className} style={style}>
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
})

const SearchSuffixContainerStyle: SxProps<Theme> = {
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

const SearchbarContainerStyle: SxProps<Theme> = {
    position: 'relative',
    '& .SearchBar_Input': {
        '&.MuiInputBase-root': {
            '& .MuiInputBase-input': {},
        },
        lineHeight: 48,
        height: 48,
    },
};