import { ChangeEvent, FC, memo } from 'react'
import { SxProps, Theme } from '@mui/material/styles';
import { Pagination } from '@src/components/form/index.js';
import { Box } from '@mui/material';
import { DefaultPageSize } from './constants.js';

export interface PaginationComponentProps {
    page: number;
    onChange(event: ChangeEvent<unknown>, page: number): void;
    count: number;
}

export const PaginationComponent: FC<PaginationComponentProps> = memo(({ page, onChange, count }) => {
    const pageCount = (count / DefaultPageSize) >> 0

    return (
        <Box sx={PaginationContainerStyle}>
            {
                count <= DefaultPageSize ?
                    undefined :
                    <Pagination
                        sx={PaginationStyle}
                        count={pageCount}
                        page={page+1}
                        onChange={onChange}
                    />
            }
        </Box>
    )
})

const PaginationStyle: SxProps<Theme> = {
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
}

const PaginationContainerStyle: SxProps<Theme> = {
    maxWidth: 600,
    marginInline: 'auto',
    padding: '16px 12px 16px',
};
