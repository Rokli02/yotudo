import { styled } from '@mui/material/styles';
import { FC } from 'react';
import StatusesComponent from './components/Statuses.component';
import { GenresComponent } from './components/Genres.component';
import { AddGenreComponent } from './components/AddGenre.component';
import { GenreProvider } from './contexts';
import { RenameGenreComponent } from './components/RenameGenre.component';

export const MiscPage: FC = () => {
  
  return (
    <Container>
      <h1>Státuszok és műfajok</h1>
      <StatusesComponent />
      <GenreProvider>
        <GenresComponent />
        <AddGenreComponent />
        <RenameGenreComponent />
      </GenreProvider>
    </Container>
  )
};

export default MiscPage;

const Container = styled('div')({
  position: 'relative',
  flexWrap: 'wrap',
  display: 'flex',
  padding: '1rem 1rem',
  gap: '1rem 2%',
  justifyContent: 'center',
  'h1': {
    width: '100%',
    textAlign: 'center',
  },
});
