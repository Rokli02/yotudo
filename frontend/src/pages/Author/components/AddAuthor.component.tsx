import { PersonAdd } from '@mui/icons-material'
import { Button as MuiButton } from '@mui/material'
import { styled } from '@mui/material/styles'
import { FC, useState } from 'react'
import AddAuthorModal from './AddAuthorModal'

export const AddAuthorComponent: FC = () => {
    const [isModalOpen, setIsModalOpen] = useState(false)

    return (
        <>
            <OpenModel variant='contained' type='button' color='error' onClick={() => setIsModalOpen(true)}><PersonAdd /></OpenModel>
            { !isModalOpen ? undefined : <AddAuthorModal open={isModalOpen} onClose={() => setIsModalOpen(false)}/>}
        </>
    )
}

export default AddAuthorComponent

const OpenModel = styled(MuiButton)({
    position: 'fixed',
    right: 'min(4%, 40px)',
    bottom: '32px',
    minWidth: 'unset',
    padding: '12px',
    borderRadius: '50%',
    boxShadow: '3px 3px 9px #0004'
})