import { FC, useEffect, useState } from 'react';
import { Container, Content, ItemContainer } from './styled.components';
import { Genre, GenreService } from '@src/api';
import { LoadingPage } from '@src/pages/Common';
import { Title } from '@src/components/common';
import { useGenreContext } from '../contexts';

export const GenresComponent: FC = () => {
    const [loading, setLoading] = useState<boolean>(true);
    const { genres, setGenres, setSelectedGenre } = useGenreContext();

    useEffect(() => {
        setLoading(true);

        GenreService.GetAllGenre().then((res) => {
        if (!res || (Array.isArray(res) && res.length == 0)) {
            return;
        }

        setGenres(res);
        }).finally(() => {
        setLoading(false);
        })
    }, [setGenres])

    return loading ? <Container><LoadingPage size='medium'/></Container>
    : genres.length === 0 ? <></> : (
        <Container className='direction_row'>
        <Title>Műfajok</Title>
        <Content data-dir='row'>
            {genres?.map((genre, index) => 
            <GenreItem key={`${index}_${genre.id}`} genre={genre} onClick={() => setSelectedGenre(genre.id, index)}/>
            )}
        </Content>
        </Container>
    )
}

const GenreItem: FC<{ genre: Genre, onClick: () => void }> = ({ genre, onClick }) => {
  return (
    <ItemContainer className='min_width' style={{cursor: 'pointer'}} onClick={onClick}>
      <div className="text-div">{genre.name}</div>
    </ItemContainer>
  )
}