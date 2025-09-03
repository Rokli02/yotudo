import { FC, useEffect, useMemo, useState } from 'react'
import { useMusicContext } from '../contexts';
import { SxProps, Theme } from '@mui/material/styles';
import { FormControl, InputLabel, Pagination, Searchbar } from '@src/components/form';
import { SearchbarProps } from '@src/components/form/searcbar';
import { Select } from '@src/components/form';
import { Option } from '@src/components/form/select';
import { StatusIcon } from '@src/components/common';
import { StatusService, Status } from '@src/api';
import { Box } from '@mui/material';

export const HeaderComponent: FC = () => {
    const { musics: { count }, page, setPage } = useMusicContext();
    const [statuses, setStatuses] = useState<Status[]>([])
    const [currentStatus, setCurrentStatus] = useState<Option<number>>({value: -1, label: 'Nincs'})

    const numOfPages = useMemo(() => {
        if (!count || !page.size) {
            return 1
        }

        return Math.ceil(count / page.size) || 1
    }, [page, count]);

    const statusOptions: Option<number>[] = useMemo(() => {
        return statuses.map(convertStatusToOption)
    }, [statuses])

    const onDebounce: SearchbarProps['onDebounce'] = (search) => {
        setPage({ filter: search }, currentStatus!.value);
    }

    function onStatusSelect<T extends number>(value: T) {
        setCurrentStatus(convertStatusToOption(statuses.find((s) => s.id === value)!));
        setPage({}, value);
    }

    useEffect(() => {
        StatusService.GetAllStatus().then((statuses) => {
            setStatuses([{ id: -1, name: 'Nincs', description: '' }, ...statuses]);
        })
    }, [])

    return (
        <Box sx={HeaderStyle}>
            <Box sx={RowStyle}>
                <div />
                <Searchbar
                    onDebounce={onDebounce}
                    debounceTime={300}
                    className='searchbar'
                />
                <FormControl className='select_status'>
                    <InputLabel>Filter</InputLabel>
                    <Select label="Filter" fullWidth options={statusOptions} value={currentStatus.value} onChange={onStatusSelect}/>
                </FormControl>
            </Box>
            <Box sx={PaginationContainerStyle}>
                {
                    numOfPages === 1 ?
                        undefined :
                        <Pagination
                            count={numOfPages}
                            page={page.page + 1}
                            onChange={(_, currentPage) => {
                                setPage({ page: currentPage - 1 }, currentStatus.value);
                            }}
                        />
                }
            </Box>
        </Box>
    )
}

function convertStatusToOption(status: Status): Option<number> {
    const Icon = StatusIcon[status.id];

    return {
        value: status.id,
        label: (
            <div style={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'start',
                columnGap: 10,
            }}>
                <Icon />
                <span>{ status.name }</span>
            </div>
        )
    }
}

const HeaderStyle: SxProps<Theme> = {
    display: 'flex',
    flexDirection: 'column',
    flexWrap: 'wrap',
    justifyContent: 'center',
    rowGap: '1rem',
    paddingTop: '1rem',
    '& .status-option': {
        backgroundColor: 'purple',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'start',
        columnGap: 10,
    },
};

const PaginationContainerStyle: SxProps<Theme> = {
    height: 36,
    display: 'grid',
    alignItems: 'center',
    justifyContent: 'center',
};

const RowStyle: SxProps<Theme> = {
    width: '100%',
    display: 'grid',
    justifyContent: 'space-between',
    alignItems: 'center',
    gridTemplateColumns: '33% 33% 33%',
    '.searchbar': {
        marginInline: 'auto',
        minWidth: 300,
        width: '100%',
        maxWidth: 400,
    },
    '.select_status': {
        marginInline: 'auto',
        minWidth: 200,
        maxWidth: 250,
    },
    '@media screen and (max-width: 925px)': {
        gridTemplateColumns: 'none',
        justifyContent: 'center',
        rowGap: 24,
        '& > *:nth-of-type(1)': {
            display: 'none',
        }
    },
};