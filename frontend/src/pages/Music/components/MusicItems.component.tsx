import { SxProps, Theme } from '@mui/material/styles'
import { FC } from 'react'
import { useMusicContext } from '../contexts'
import { MusicItem } from './MusicItem';
import {useSearchParams } from 'react-router-dom';
import { Box } from '@mui/material';

export const MusicItemsComponent: FC = () => {
    const { musics, performAction } = useMusicContext();
    const [searchParams, setSearchParams] = useSearchParams();

    if (musics.count === 0 || !musics.data?.length) {
        return <h3 style={{ fontSize: '1.75rem', textAlign: 'center' }}>Nincs találat!</h3>;
    }

    return (
        <Box sx={ContentStyle}>
            { musics.data.map((music, index) => 
                <MusicItem
                    key={`${index}_${music.id}`}
                    music={music}
                    onAction={music.status.id === 1 ? undefined : () => performAction(music)}
                    onActionAfterHold={() =>  {
                        const idSP = searchParams.get('id');
                        const indexSP = searchParams.get('index');
                        if ((idSP && idSP === music.id + '') && (indexSP && indexSP === index + '')) {
                            return
                        }

                        setSearchParams({ id: music.id + '', index: index + '' }, { replace: true })
                    }}
                />)
            }
        </Box>
    )
}

const ContentStyle: SxProps<Theme> = {
    display: 'flex',
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: '.5rem .75rem',
    justifyContent: 'center',
    alignItems: 'flex-start',
};