import { FC } from 'react'
import AuthorItem from './AuthorItem'
import { useAuthorContext } from '../contexts'
import { styled } from '@mui/material/styles'

export const AuthorItemsComponent: FC = () => {
    const { authors, deleteAuthor } = useAuthorContext()

    if (authors.count === 0 || !authors.data?.length) {
        return <h3 style={{ textAlign: 'center', fontSize: '1.75rem' }}>Nincs találat!</h3>
    }

    return (
        <Container>
            {authors.data.map((author, index) => <AuthorItem key={`${index}_${author.id}`} author={author} onDelete={deleteAuthor} />)}
        </Container>
    )
}

export default AuthorItemsComponent

const Container = styled('div')({
    position: 'relative',
    maxWidth: '1200px',
    width: '75%',
    minWidth: 'fit-content',
    marginInline: 'auto',
    padding: '1rem 1rem',
    rowGap: '1rem',
    display: 'flex',
    flexWrap: 'wrap',
    flexDirection: 'column',
    justifyContent: 'center',
    alignItems: 'center',
    '@media screen and (max-width: 710px)': {
        width: '100%',
    },
})