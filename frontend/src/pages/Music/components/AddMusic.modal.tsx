import { DialogActions, DialogContent } from '@mui/material'
import { FC } from 'react'
import { TopRightButton } from '@src/components/common'
import { Button, FormControl, InputLabel } from '@src/components/form'
import { Close } from '@mui/icons-material'
import {
    Form,
    FormInput,
    FormCheckbox,
    FormAutocomplete,
    FormMultiselectAutocomplete,
} from '@src/contexts/form'
import { IForm } from '@src/contexts/form/interface'
import { NewMusic } from '@src/api'
import { useMusicContext } from '../contexts'
import {
    CustomDialag,
    getAuthorOptions,
    getContributorOptions,
    getGenreOptions,
    musicConstraints,
    Title,
    transformFormObjectToNewMusic,
} from './musicForm.utils'

export const AddMusicModal: FC<{ open: boolean, onClose: () => void }> = ({ open, onClose }) => {
    const { addMusic } = useMusicContext();
    const onSubmit: IForm['onSubmit'] = async (value: NewMusic) => {
        console.log(value)
        return addMusic(value);
    }

    return (
        <CustomDialag open={open} onClose={onClose}>
            <TopRightButton Icon={Close} onClick={onClose} />
            <Title>Új zene hozzáadás</Title>
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
                        <InputLabel>URL</InputLabel>
                        <FormInput name='url' type='url'/>
                    </FormControl>
                    <FormControl>
                        <InputLabel>Album</InputLabel>
                        <FormInput name='album' type='text'/>
                    </FormControl>
                    <FormControl>
                        <InputLabel>Kiadás dátuma</InputLabel>
                        <FormInput name='published' type='number'/>
                    </FormControl>
                    <FormControl>
                        <FormAutocomplete debounceTime={600} name='author' label='Szerző' getOptions={getAuthorOptions}/>
                    </FormControl>
                    <FormControl>
                        <FormAutocomplete debounceTime={600} fetchOnce name='genre' label='Műfaj' getOptions={getGenreOptions}/>
                    </FormControl>
                    <FormCheckbox label='Videó indexkép borítóképnek' name='useThumbnail' />
                    <FormControl>
                        <FormMultiselectAutocomplete debounceTime={600} name='contributor' label='Közreműködők' getOptions={getContributorOptions} renderChipContent={(v) => v.label}/>
                    </FormControl>
                </DialogContent>
                <DialogActions>
                    <Button type='submit' color='success'>Mentése</Button>
                </DialogActions>
            </Form>
        </CustomDialag>
    )
}
