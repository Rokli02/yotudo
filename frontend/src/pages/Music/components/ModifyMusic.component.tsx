import { FC, useMemo } from 'react'
import { useSearchParams } from 'react-router-dom'
import { useMusicContext } from '../contexts';
import { ModifyMusicModal } from './ModifyMusic.modal';

export const ModifyMusicComponent: FC = () => {
    const [searchParams, setSearchParams] = useSearchParams();
    const { musics, modifyMusic } = useMusicContext();

    const spData = useMemo(() => {
        const id = Number(searchParams.get('id'));
        let index = Number(searchParams.get('index'))

        if (isNaN(index)) index = 0;

        if (isNaN(id) || !musics || musics.data.length === 0) {
            return null
        }

        if (musics.data[index].id === id) {
            return musics.data[index]
        }

        return musics.data.find((music) => music.id === id) ?? null
    }, [searchParams, musics])

    return (
        <>
            { !spData ? undefined : <ModifyMusicModal open={!!spData} onClose={() => setSearchParams({})} music={spData} onSubmit={modifyMusic}/> }
        </>
    )
}
