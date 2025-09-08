import { FC } from 'react';
import { Container, Content } from './styled.components';
import { useGetData } from '@src/hooks/useGetData';
import { StatusService, Status } from '@src/api';
import { LoadingPage } from '@src/pages/Common';
import { Box, SxProps, Theme } from '@mui/material';
import { StatusIcon, Title } from '@src/components/common';

export const StatusesComponent: FC = () => {
    const [loading, statuses] = useGetData(() => StatusService.GetAllStatus());

    return loading ? <Container><LoadingPage size='medium'/></Container> : (
        <Container>
            <Title>Státuszok</Title>
            <Content data-dir="col">
                {statuses?.map((status, index) =>
                    <StatusItem key={`${index}_${status.id}`} status={status} />
                )}
            </Content>
        </Container>
    )
}

const StatusItem: FC<{ status: Status}> = ({ status }) => {
    const _StatusIcon = StatusIcon[status.id];

    return (
        <Box sx={ItemContainerStyle}>
            <_StatusIcon />
            <Box sx={StatusNameStyle} >{status.name}</Box>
            <Box sx={StatusDescriptionStyle} >{status.description}</Box>
        </Box>
    )
}

export default StatusesComponent;

const StatusNameStyle: SxProps<Theme> = {
    width: '12ch',
    flexWrap: 'nowrap',
}

const StatusDescriptionStyle: SxProps<Theme> = {
    flexWrap: 'wrap',
}

const ItemContainerStyle: SxProps<Theme> = {
    width: '100%',
    backgroundColor: '#0003',
    height: 'fit-content',
    padding: '10px 16px',
    borderRadius: '8px',
    columnGap: '12px',
    overflowX: 'clip',
    cursor: 'default',
    display: 'grid',
    gridTemplateColumns: 'min-content max-content auto',
    alignItems: 'center',
}