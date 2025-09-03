import { FC } from 'react'
import { TopRightButton, Divider } from '@src/components/common'
import { Close } from '@mui/icons-material'
import {
    Form,
    FormInput,
    FormAutocomplete,
    FormMultiselectAutocomplete,
} from '@src/contexts/form'
import { IForm } from '@src/contexts/form/interface'
import { Button, FormControl, InputLabel } from '@src/components/form'
import { DialogActions, DialogContent } from '@mui/material'
import { Music, MusicUpdate, NewMusic } from '@src/api'
import {
    CustomDialag,
    getAuthorOptions,
    getContributorOptions,
    getGenreOptions,
    musicConstraints,
    Title,
    transformFormObjectToNewMusic,
} from './musicForm.utils'
import { FormImageSelector } from '@src/contexts/form/FormImageSelector'

interface ModifyMusicModalProps {
    open: boolean,
    onClose: () => void,
    music: Music,
    onSubmit: (musicToUpdate: MusicUpdate, index?: number) => Promise<boolean>,
}

export const ModifyMusicModal: FC<ModifyMusicModalProps> = ({ open, onClose, music, onSubmit }) => {
    const _onSubmit: IForm['onSubmit'] = async (value: NewMusic) => {
        const response = await onSubmit({ ...music, ...value});

        if (response) {
            onClose()
        }

        return response;
    }

    return ( <CustomDialag open={open} onClose={onClose}>
        <TopRightButton Icon={Close} onClick={onClose} />
        <Title>Zene módosítás</Title>
        <Form
            onSubmit={_onSubmit}
            transformFlatObjectTo={transformFormObjectToNewMusic}
            constraints={musicConstraints}
        >
            <DialogContent className='form_items'>
                <FormControl>
                    <InputLabel>Cím</InputLabel>
                    <FormInput name='name' value={music['name']} />
                </FormControl>
                <FormControl>
                    <FormAutocomplete
                        debounceTime={450}
                        name='author'
                        label='Szerző'
                        getOptions={getAuthorOptions}
                        value={{ label: music['author'].name, ...music['author'] }}
                    />
                </FormControl>
                <FormControl>
                    <InputLabel>URL</InputLabel>
                    <FormInput name='url' type='url' value={music['url']}/>
                </FormControl>
                <FormControl>
                    <FormMultiselectAutocomplete
                        debounceTime={450}
                        name='contributor'
                        label='Közreműködők'
                        getOptions={getContributorOptions}
                        renderChipContent={(v) => v.label}
                        selectedOptions={music['contributor']?.map((v) => ({ label: v.name, ...v }))}
                    />
                </FormControl>
                <Divider dir='horizontal' length='570px' sx={{ backgroundColor: 'var(--primary-color)' }}/>
                <FormControl>
                    <InputLabel>Album</InputLabel>
                    <FormInput name='album' type='text' value={music['album']}/>
                </FormControl>
                <FormControl>
                    <InputLabel>Kiadás dátuma</InputLabel>
                    <FormInput name='published' type='number' value={music['published']}/>
                </FormControl>
                <FormControl>
                    <FormAutocomplete
                        debounceTime={450}
                        fetchOnce
                        name='genre'
                        label='Műfaj'
                        getOptions={getGenreOptions}
                        value={{ label: music['genre'].name, ...music['genre'] }}
                    />
                </FormControl>
                <FormControl>
                    <FormImageSelector
                        name='picName'
                        defaultValue={music['picName']}
                    />
                </FormControl>
            </DialogContent>
            <DialogActions>
                <Button type='submit' color='success'>Módosítás</Button>
            </DialogActions>
        </Form>
    </CustomDialag>
    )
}