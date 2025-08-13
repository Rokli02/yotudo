import { FC } from 'react';
import { styled } from '@mui/material/styles';
import AuthorItemsComponent from './components/AuthorItems.component';
import { AuthorProvider } from './contexts/AuthorContext';
import AddAuthorComponent from './components/AddAuthor.component';
import { HeaderComponent } from './components/Header.component';

export const AuthorPage: FC = () => {
    return (
        <PageContainer>
            <AuthorProvider>
            <h1>Szerzők</h1>
                <HeaderComponent />
                <AuthorItemsComponent />
                <AddAuthorComponent />
            </AuthorProvider>
        </PageContainer>
    )
}

export default AuthorPage;

const PageContainer = styled('div')({
    position: 'relative',
    textAlign: 'center',
    padding: '1rem 1rem',    
})

