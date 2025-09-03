import { Dialog, DialogActions, DialogContent, DialogTitle } from '@mui/material'
import { FC } from 'react'
import { Divider, TopRightButton } from '@src/components/common'
import { Button, FormControl, InputLabel } from '@src/components/form'
import { Close } from '@mui/icons-material'
import {
    Form,
    FormInput,
    FormAutocomplete,
    FormMultiselectAutocomplete,
} from '@src/contexts/form'
import { IForm } from '@src/contexts/form/interface'
import { NewMusic } from '@src/api'
import { useMusicContext } from '../contexts'
import {
    CustomDialagStyle,
    getAuthorOptions,
    getContributorOptions,
    getGenreOptions,
    musicConstraints,
    TitleStyle,
    transformFormObjectToNewMusic,
} from './musicForm.utils'
import { FormImageSelector } from '@src/contexts/form/FormImageSelector'

export const AddMusicModal: FC<{ open: boolean, onClose: () => void }> = ({ open, onClose }) => {
    const { addMusic } = useMusicContext();
    const onSubmit: IForm['onSubmit'] = async (value: NewMusic) => {
        return addMusic(value);
    }

    return (
        <Dialog sx={CustomDialagStyle} open={open} onClose={onClose}>
            <TopRightButton Icon={Close} onClick={onClose} />
            <DialogTitle sx={TitleStyle}>Új zene hozzáadás</DialogTitle>
            <Form
                onSubmit={onSubmit}
                transformFlatObjectTo={transformFormObjectToNewMusic}
                constraints={musicConstraints}
            >
                <DialogContent className='form_items'>
                    <FormControl>
                        <InputLabel>Cím</InputLabel>
                        <FormInput name='name' />
                    </FormControl>
                    <FormControl>
                        <FormAutocomplete
                            debounceTime={600}
                            name='author'
                            label='Szerző'
                            getOptions={getAuthorOptions}
                        />
                    </FormControl>
                    <FormControl>
                        <InputLabel>URL</InputLabel>
                        <FormInput name='url' type='url'/>
                    </FormControl>
                    <FormControl>
                        <FormMultiselectAutocomplete
                            debounceTime={600}
                            name='contributor'
                            label='Közreműködők'
                            getOptions={getContributorOptions}
                            renderChipContent={(v) => v.label}
                        />
                    </FormControl>
                    <Divider dir='horizontal' length='570px' sx={{ backgroundColor: 'var(--primary-color)' }}/>
                    <FormControl>
                        <InputLabel>Album</InputLabel>
                        <FormInput name='album' type='text'/>
                    </FormControl>
                    <FormControl>
                        <InputLabel>Kiadás dátuma</InputLabel>
                        <FormInput name='published' type='number'/>
                    </FormControl>
                    <FormControl>
                        <FormAutocomplete
                            debounceTime={600}
                            fetchOnce
                            name='genre'
                            label='Műfaj'
                            getOptions={getGenreOptions}
                        />
                    </FormControl>
                    <FormControl>
                        <FormImageSelector
                            name='picName'
                        />
                    </FormControl>
                </DialogContent>
                <DialogActions>
                    <Button type='submit' color='success'>Mentése</Button>
                </DialogActions>
            </Form>
        </Dialog>
    )
}
