import { styled } from '@mui/material/styles';
import { FC } from 'react'
import { MusicItemsComponent } from './components/MusicItems.component';
import { MusicProvider } from './contexts';
import { AddMusicComponent } from './components/AddMusic.component';
import { HeaderComponent } from './components/Header.component';
import { ModifyMusicComponent } from './components/ModifyMusic.component';

export const MusicPage: FC = () => {
  return (
    <PageContainer>
      <MusicProvider>
        <h1>Zenék</h1>
        <HeaderComponent />
        <MusicItemsComponent />
        <ModifyMusicComponent />
        <AddMusicComponent />
      </MusicProvider>
    </PageContainer>
  )
}

export default MusicPage;

const PageContainer = styled('div')({
    position: 'relative',
    padding: '1rem 1rem',
    '& > h1': {
        textAlign: 'center',
    },
})