import { Dialog, DialogActions, DialogContent, DialogTitle } from '@mui/material'
import { styled, CSSObject } from '@mui/material/styles'
import { FC } from 'react'
import { TopRightButton } from '@src/components/common'
import { Button, FormControl, InputLabel } from '@src/components/form'
import { Form, FormConstraints, FormInput } from '@src/contexts/form'
import { Close } from '@mui/icons-material'
import { useAuthorContext } from '../contexts'
import { NewAuthor } from '@src/api'
import { FormSlider } from '@src/contexts/form/FormSlider'

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
        <CustomDialag open={open} fullWidth onClose={onClose}>
            <TopRightButton Icon={Close} onClick={onClose} />
            <Title>
                Új zenész felvétel
            </Title>
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
        </CustomDialag>
    )
}

export default AddAuthorModal

const CustomDialag = styled(Dialog)({
    '& .MuiDialog-container': {
        '& .MuiPaper-root': {
            maxWidth: '650px',
            width: '100%',
            backgroundColor: 'var(--background-color)',
            color: 'var(--font-color)',
            '& .MuiDialogContent-root': {
                width: '100%',
                '& > .MuiFormLabel-root': {
                    marginLeft: '6ch',
                } as CSSObject,
                '& > .MuiFormControl-root > .MuiSlider-root': {
                    marginLeft: '21px',
                } as CSSObject,
            } as CSSObject,
        } as CSSObject,
    } as CSSObject,
    '& .form_items': {
        display: 'flex',
        flexDirection: 'column',
        rowGap: '1rem',
    } as CSSObject,
})

const Title = styled(DialogTitle)({
    fontSize: '1.25rem',
    marginInline: 'auto',
    width: 'max-content',
})