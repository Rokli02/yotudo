import { createContext, FC, ReactElement, useEffect, useState } from 'react'
import { Author, AuthorService, NewAuthor, Page, Pagination } from '@src/api';
import { PageSetter, usePage } from '@src/hooks/usePage';

export interface IAuthorContext {
    authors: Pagination<Author[]>;
    page: Page;
    setPage: PageSetter,
    addAuthor: (author: NewAuthor) => Promise<boolean>;
    deleteAuthor: (id: number) => Promise<boolean>;
}

const PAGE_SIZE = 12;

export const AuthorContext = createContext<IAuthorContext>(null as unknown as IAuthorContext);

export const AuthorProvider: FC<{ children: ReactElement | ReactElement[] }> = ({ children }) => {
    const [authors, setAuthors] = useState<Pagination<Author[]>>({ data: [], count: 0 });
    const [page, setPage, _setPage] = usePage(PAGE_SIZE, (state) => AuthorService.GetAuthors(state).then(setAuthors));

    async function addAuthor(author: NewAuthor) {
        if (!author) return false;

        const newAuthor = await AuthorService.SaveAuthor(author);
        if (!newAuthor) return false;

        setAuthors((pre) => {
            if (pre.data.unshift(newAuthor) > page.size) {
                pre.data.pop();
            }

            pre.count++;

            return {...pre};
        })

        return true;
    }

    async function deleteAuthor(id: number) {
        const response = await AuthorService.DeleteAuthor(id);

        if (response) {
            if (authors.data.length <= 1) {
                _setPage((pre) => {
                    return { ...pre, page: Math.max(pre.page - 1, 0) };
                });
                
            }

            AuthorService.GetAuthors(page).then(setAuthors);
        }

        return response;
    }

    useEffect(() => {
        AuthorService.GetAuthors(page).then(setAuthors)
    }, [])

    return (
        <AuthorContext.Provider value={{
            authors,
            page,
            setPage,
            addAuthor,
            deleteAuthor,
        }}>
            {children}
        </AuthorContext.Provider>
    )
}