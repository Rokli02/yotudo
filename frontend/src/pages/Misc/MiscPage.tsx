import { SxProps, Theme } from '@mui/material/styles';
import { FC } from 'react';
import StatusesComponent from './components/Statuses.component';
import { GenresComponent } from './components/Genres.component';
import { AddGenreComponent } from './components/AddGenre.component';
import { GenreProvider } from './contexts';
import { RenameGenreComponent } from './components/RenameGenre.component';
import { Box } from '@mui/material';

export const MiscPage: FC = () => {
  
  return (
    <Box sx={ContainerStyle}>
      <h1>Státuszok és műfajok</h1>
      <StatusesComponent />
      <GenreProvider>
        <GenresComponent />
        <AddGenreComponent />
        <RenameGenreComponent />
      </GenreProvider>
    </Box>
  )
};

export default MiscPage;

const ContainerStyle: SxProps<Theme> = {
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
};
