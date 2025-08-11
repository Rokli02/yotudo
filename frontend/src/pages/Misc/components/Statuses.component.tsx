import { FC } from 'react';
import { Container, Content, ItemContainer } from './styled.components';
import { useGetData } from '@src/hooks/useGetData';
import { StatusService, Status } from '@src/api';
import { LoadingPage } from '@src/pages/Common';
import { Divider } from '@mui/material';
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
        <ItemContainer className='data_status'>
            <_StatusIcon />
            <Divider orientation='vertical' variant='fullWidth' />
            <div className='text-div no-wrap'>{status.name}</div>
            <Divider orientation='vertical' variant='fullWidth' />
            <div className='text-div'>{status.description}</div>
        </ItemContainer>
    )
}

export default StatusesComponent;
