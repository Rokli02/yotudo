import { FC, FormEvent, useEffect, useState } from 'react'
import { Container } from './styled.components'
import { Title, TopRightButton } from '@src/components/common'
import { Close } from '@mui/icons-material'
import { useGenreContext } from '../contexts'
import { Button, Form, FormControl, Input, InputLabel } from '@src/components/form'
import { Genre, GenreService } from '@src/api'
import { SxProps, Theme, Typography } from '@mui/material'

export const RenameGenreComponent: FC = () => {
    const { selectedGenre, setSelectedGenre, renameGenre, deleteGenre} = useGenreContext();

    return !selectedGenre
        ? <></>
        : <RenameGenre
            selectedGenre={selectedGenre}
            renameGenre={renameGenre}
            onClose={() => setSelectedGenre(null)}
            deleteGenre={deleteGenre}
        />
}

const RenameGenre: FC<{
    selectedGenre: Genre,
    onClose: () => void,
    renameGenre: ReturnType<typeof useGenreContext>['renameGenre'],
    deleteGenre: ReturnType<typeof useGenreContext>['deleteGenre'],
}> = ({ selectedGenre, onClose, renameGenre, deleteGenre }) => {
    const [name, setName] = useState<string>(!selectedGenre ? '' : selectedGenre.name)

    async function onSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault();

        const renamedGenre = await GenreService.RenameGenre(selectedGenre!.id, { name })

        if (renamedGenre) return renameGenre(renamedGenre.id, { name: renamedGenre.name });
    }

    function revertName() {
        setName(selectedGenre?.name ?? '')
    }

    function deleteGenreById() {
        deleteGenre(selectedGenre.id).then((canDelete) => {
            if (canDelete) onClose()
        })
    }

    useEffect(() => {
        if (!selectedGenre && !name) return

        setName(selectedGenre?.name ?? '');
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [selectedGenre])

    return (
        <Container>
            <Typography component={'label'} sx={TopLeftIdTextStyle}>id: {selectedGenre.id}</Typography>
            <TopRightButton Icon={Close} onClick={onClose}/>
            <Title>Műfaj átnevezés</Title>
            <Form onSubmit={onSubmit}>
                <FormControl>
                    <InputLabel htmlFor='name'>Műfaj neve</InputLabel>
                    <Input value={name} id='name' onChange={(e) => setName(e.target.value)}/>
                </FormControl>
                <div style={{ display: 'flex', flexDirection: 'row', columnGap: '1rem' }}>
                    <Button type='submit' color='success' variant='text'>Módosítás</Button>
                    <Button type='button' color='primary' variant='outlined' onClick={revertName}>Undo</Button>
                    <Button type='button' color='error' variant='outlined' onClick={deleteGenreById}>Törlés</Button>
                </div>
            </Form>
        </Container>
    )
}

const TopLeftIdTextStyle: SxProps<Theme> = { // label
    position: 'absolute',
    top: '.5rem',
    left: '1rem',
    color: '#fff8',
    fontSize: '.8rem',
}