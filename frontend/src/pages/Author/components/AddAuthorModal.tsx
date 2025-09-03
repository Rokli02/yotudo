import { Dialog, DialogActions, DialogContent, DialogTitle } from '@mui/material';
import { FC } from 'react';
import { TopRightButton } from '@src/components/common';
import { Button, FormControl, InputLabel } from '@src/components/form';
import { Form, FormConstraints, FormInput } from '@src/contexts/form';
import { Close } from '@mui/icons-material';
import { useAuthorContext } from '../contexts';
import { NewAuthor } from '@src/api';
import { DialogStyle, DialogTitleStyle } from './common.component';

const addAuthorConstraints: FormConstraints = {
    name: (_value, errors) => {
        const value = _value as string;

        if (!value || !value.trim()) errors.push('* Kötelező mező')
    },
}

export const AddAuthorModal: FC<{ open: boolean, onClose: () => void }> = ({ open, onClose }) => {
    const { addAuthor } = useAuthorContext();

    const onSubmit = async (newAuthor: NewAuthor) => {
        return await addAuthor(newAuthor);
    }
    
    return (
        <Dialog sx={DialogStyle} open={open} fullWidth onClose={onClose}>
            <TopRightButton Icon={Close} onClick={onClose} />
            <DialogTitle sx={DialogTitleStyle}>
                Új zenész felvétel
            </DialogTitle>
            <Form onSubmit={onSubmit} constraints={addAuthorConstraints}>
                <DialogContent className='form_items'>
                    <FormControl>
                        <InputLabel>Név</InputLabel>
                        <FormInput name='name' type='text'/>
                    </FormControl>
                </DialogContent>
                <DialogActions>
                    <Button type='submit' color='success'>Mentés</Button>
                </DialogActions>
            </Form>
        </Dialog>
    )
}

export default AddAuthorModal