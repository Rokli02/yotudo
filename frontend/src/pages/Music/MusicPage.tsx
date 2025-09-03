import { SxProps, Theme } from '@mui/material/styles';
import { FC } from 'react'
import { MusicItemsComponent } from './components/MusicItems.component';
import { MusicProvider } from './contexts';
import { AddMusicComponent } from './components/AddMusic.component';
import { HeaderComponent } from './components/Header.component';
import { ModifyMusicComponent } from './components/ModifyMusic.component';
import { Box } from '@mui/material';

export const MusicPage: FC = () => {
  return (
    <Box sx={ContainerStyle}>
      <MusicProvider>
        <h1>Zenék</h1>
        <HeaderComponent />
        <MusicItemsComponent />
        <ModifyMusicComponent />
        <AddMusicComponent />
      </MusicProvider>
    </Box>
  )
}

export default MusicPage;

const ContainerStyle: SxProps<Theme> = {
    position: 'relative',
    padding: '1rem 1rem',
    '& > h1': {
        textAlign: 'center',
    },
}