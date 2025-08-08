import { createContext, FC, ReactElement, useEffect, useState } from 'react'
import { Author, AuthorService, NewAuthor, Page, Pagination } from '@src/api';

export interface IAuthorContext {
    authors: Pagination<Author[]>;
    page: Page;
    setPage: (page: Partial<Page>) => void,
    addAuthor: (author: NewAuthor) => Promise<boolean>;
    deleteAuthor: (id: number) => Promise<boolean>;
}

const PAGE_SIZE = 5;

export const AuthorContext = createContext<IAuthorContext>(null as unknown as IAuthorContext);

export const AuthorProvider: FC<{ children: ReactElement | ReactElement[] }> = ({ children }) => {
    const [authors, setAuthors] = useState<Pagination<Author[]>>({ data: [], count: 0 });
    const [page, _setPage] = useState<Page>({ page: 0, size: PAGE_SIZE });

    async function addAuthor(author: NewAuthor) {
        if (!author) return false;

        const newAuthor = await AuthorService.SaveAuthor(author);
        if (!newAuthor) return false;

        setAuthors((pre) => {
            if (pre.data.unshift(newAuthor) > page.size) {
                pre.data.pop();
            }

            return {...pre};
        })

        return true;
    }

    async function deleteAuthor(id: number) {
        const response = await AuthorService.DeleteAuthor(id)

        if (response) {
            if (authors.data.length == 0) {
                _setPage((pre) => ({ ...pre, page: Math.max(pre.page - 1, 0) }))
            }

            AuthorService.GetAuthors(page).then(setAuthors)
        }

        return response;
    }

    function setPage(pageUpdate: Partial<Page>) {
        const modifiedKeys = Object.entries(pageUpdate) as Array<[keyof Page, unknown]>
        if (
            modifiedKeys.length === 0 ||
            modifiedKeys.every(([key, value]) => page[key] === value)
        ) {
            console.log("Unnecessary state update was blocked")

            return;
        }

        _setPage((pre) => {
            const newState = {
                ...pre,
                ...pageUpdate,
            }
            
            return newState;
        })

        const mdfky = {...page, ...pageUpdate};

        AuthorService.GetAuthors(mdfky).then(setAuthors)
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