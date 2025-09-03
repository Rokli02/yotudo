import { FC } from 'react';
import AuthorItem from './AuthorItem';
import { useAuthorContext } from '../contexts';
import { Box } from '@mui/material';

export const AuthorItemsComponent: FC = () => {
    const { authors, deleteAuthor } = useAuthorContext()

    if (authors.count === 0 || !authors.data?.length) {
        return <h3 style={{ textAlign: 'center', fontSize: '1.75rem' }}>Nincs találat!</h3>
    }

    return (
        <Box sx={ContainerStyle}>
            {authors.data.map((author, index) => <AuthorItem key={`${index}_${author.id}`} author={author} onDelete={deleteAuthor} />)}
        </Box>
    )
}

export default AuthorItemsComponent

const ContainerStyle = {
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
};