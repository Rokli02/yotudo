import { FC, FormEventHandler, useEffect, useState } from 'react'
import { Container, TopLeftIdText } from './styled.components'
import { Title, TopRightButton } from '@src/components/common'
import { Close } from '@mui/icons-material'
import { useGenreContext } from '../contexts'
import { Button, Form, FormControl, Input, InputLabel } from '@src/components/form'
import { Genre, GenreService } from '@src/api'

export const RenameGenreComponent: FC = () => {
    const { selectedGenre, setSelectedGenre, renameGenre } = useGenreContext();

    return !selectedGenre
        ? <></>
        : <RenameGenre selectedGenre={selectedGenre} renameGenre={renameGenre} onClose={() => setSelectedGenre(null)} />
}

export const RenameGenre: FC<{
    selectedGenre: Genre,
    onClose: () => void,
    renameGenre: ReturnType<typeof useGenreContext>['renameGenre'],
}> = ({ selectedGenre, onClose, renameGenre }) => {
    const [name, setName] = useState<string>(!selectedGenre ? '' : selectedGenre.name)

    const onSubmit: FormEventHandler<HTMLFormElement> = async (event) => {
        event.preventDefault();

        const renamedGenre = await GenreService.RenameGenre(selectedGenre!.id, { name })

        if (renamedGenre) return renameGenre(renamedGenre.id, { name: renamedGenre.name });
    }

    const revertName = () => {
        setName(selectedGenre?.name ?? '')
    }

    useEffect(() => {
        if (!selectedGenre && !name) return

        setName(selectedGenre?.name ?? '');
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [selectedGenre])

    return (
        <Container>
            <TopLeftIdText>id: {selectedGenre.id}</TopLeftIdText>
            <TopRightButton Icon={Close} onClick={onClose}/>
            <Title>RenameGenreComponent</Title>
            <Form onSubmit={onSubmit}>
                <FormControl>
                    <InputLabel htmlFor='name'>Műfaj neve</InputLabel>
                    <Input value={name} id='name' onChange={(e) => setName(e.target.value)}/>
                </FormControl>
                <div style={{ display: 'flex', flexDirection: 'row', columnGap: '1rem' }}>
                    <Button type='submit' color='success' variant='text'>Módosítás</Button>
                    <Button type='button' color='primary' variant='outlined' onClick={revertName}>Undo</Button>
                </div>
            </Form>
        </Container>
    )
}