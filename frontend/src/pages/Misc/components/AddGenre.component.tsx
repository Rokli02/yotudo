import { FC, FormEventHandler, useRef, useState } from 'react'
import { Container } from './styled.components'
import { FormControl, Button, Input, InputLabel, Form } from '@src/components/form';
import { Title } from '@src/components/common';
import { GenreService } from '@src/api';
import { useGenreContext } from '../contexts';

export const AddGenreComponent: FC = () => {
    const [name, setName] = useState<string>('')
    const loading = useRef<boolean>(false)
    const { addGenre } = useGenreContext()

    const onSubmit: FormEventHandler<HTMLFormElement> = async (event) => {
        event.preventDefault();
        if (loading.current) {
            return;
        }

        const cleanedName = name.trim();
        if (!cleanedName) {
            return;
        }

        loading.current = true;
        const newGenre = await GenreService.SaveGenre({ name: cleanedName });
        loading.current = false;

        if (newGenre) {
            setName('');
            return addGenre(newGenre);
        }
    }

    return (
        <Container>
            <Title>Új műfaj hozzáadás</Title>
            <Form onSubmit={onSubmit}>
                <FormControl>
                    <InputLabel htmlFor='name'>Műfaj neve</InputLabel>
                    <Input value={name} id='name' onChange={(e) => setName(e.target.value)}/>
                </FormControl>
                <Button type='submit' color='success' variant='text'>Létrehozás</Button>
            </Form>
        </Container>
    )
}
