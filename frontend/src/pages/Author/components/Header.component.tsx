import { FC, useMemo } from 'react'
import { Searchbar } from '@src/components/form';
import { SearchbarProps } from '@src/components/form/searcbar';
import { useAuthorContext } from '../contexts';
import { Pagination } from '@src/components/form';
import { Box } from '@mui/material';

export const HeaderComponent: FC = () => {
    const { authors: { count }, page, setPage } = useAuthorContext();

    const numOfPages = useMemo(() => {
        if (!count || !page.size) {
            return 1
        }

        return Math.ceil(count / page.size) || 1
    }, [page, count]);

    const onDebounce: SearchbarProps['onDebounce'] = (search) => {
        if (page.filter !== search) {
            setPage({ filter: search });
        }
    }

    return (
        <Box sx={HeaderStyle}>
            <Searchbar
                onDebounce={onDebounce}
                debounceTime={300}
                style={{ minWidth: '400px', maxWidth: '700px' }}
            />
            <Box sx={PaginationContainerStyle}>
                {
                    numOfPages === 1 ?
                        undefined :
                        <Pagination
                            count={numOfPages}
                            page={page.page + 1}
                            onChange={(_, currentPage) => {
                                setPage({ page: currentPage - 1 })
                            }}
                        />
                }
            </Box>

        </Box>
    )
}

const HeaderStyle = {
    display: 'grid',
    justifyContent: 'center',
    rowGap: '1rem',
    paddingTop: '1rem',
};

const PaginationContainerStyle = {
    display: 'grid',
    alignItems: 'center',
    justifyContent: 'center',
    height: 36,
};
