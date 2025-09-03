import { FC } from 'react';
import AuthorItemsComponent from './components/AuthorItems.component';
import { AuthorProvider } from './contexts/AuthorContext';
import AddAuthorComponent from './components/AddAuthor.component';
import { HeaderComponent } from './components/Header.component';
import { Box } from '@mui/material';

export const AuthorPage: FC = () => {
    return (
        <Box sx={ContainerStyle}>
            <AuthorProvider>
            <h1>Szerzők</h1>
                <HeaderComponent />
                <AuthorItemsComponent />
                <AddAuthorComponent />
            </AuthorProvider>
        </Box>
    )
}

export default AuthorPage;

const ContainerStyle = {
    position: 'relative',
    textAlign: 'center',
    padding: '1rem 1rem',    
}

