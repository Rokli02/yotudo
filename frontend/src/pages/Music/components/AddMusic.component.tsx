import { styled } from '@mui/material/styles'
import { Button as MuiButton } from '@mui/material'
import { FC, useState } from 'react'
import { PlaylistAdd } from '@mui/icons-material'
import { AddMusicModal } from './AddMusic.modal'

export const AddMusicComponent: FC = () => {
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false)

    return (
        <>
            <OpenModel variant='contained' type='button' color='error' onClick={() => setIsModalOpen(true)}><PlaylistAdd /></OpenModel>
            { !isModalOpen ? undefined : <AddMusicModal open={isModalOpen} onClose={() => setIsModalOpen(false)}/> }
        </>
    )
}

const OpenModel = styled(MuiButton)({
    position: 'fixed',
    right: 'min(4%, 40px)',
    bottom: '32px',
    minWidth: 'unset',
    padding: '12px',
    borderRadius: '50%',
    boxShadow: '3px 3px 9px #0004'
})