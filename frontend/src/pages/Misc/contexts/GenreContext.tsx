import { createContext, Dispatch, FC, ReactElement, SetStateAction, useEffect, useState } from 'react'
import { Genre, NewGenre } from '@src/api'

export interface IGenreContext {
    genres: Genre[];
    selectedGenre: Genre | null;
    setGenres: Dispatch<SetStateAction<Genre[]>>;
    setSelectedGenre: (id: number | null, index?: number) => void;
    addGenre: (genre: Genre) => void;
    renameGenre: (id: number, genre: NewGenre) => void;
}

export const GenreContext = createContext<IGenreContext>(null as unknown as IGenreContext)

export const GenreProvider: FC<{ children: ReactElement | ReactElement[] }> = ({ children }) => {
    const [genres, setGenres] = useState<Genre[]>([])
    const [selectedGenre, setSelectedGenreState] = useState<Genre | null>(null)

    const addGenre: IGenreContext['addGenre'] = (genre) => {
        setGenres(pre => [...pre, genre])
    }

    const renameGenre: IGenreContext['renameGenre'] = (id, genre) => {
        setGenres((pre) => pre.map((value) => {
            if (value.id === id) {
                value.name = genre.name
            }

            return value;
        }))
    }

    const setSelectedGenre: IGenreContext['setSelectedGenre'] = (id, index) => {
        if (id === null || id <= 0) {
            return setSelectedGenreState(null);
        }

        if (index !== undefined && genres[index].id === id) {
            return setSelectedGenreState(genres[index]);
        }

        for (let index = 0; index < genres.length; index++) {
            if (genres[index].id === id) {
                return setSelectedGenreState(genres[index]);
            }
        }

        setSelectedGenreState(null);
    }

    useEffect(() => {
        console.log('Selected genre:', selectedGenre)
    }, [selectedGenre])

    return (
        <GenreContext.Provider value={{
            genres,
            selectedGenre,
            setGenres,
            setSelectedGenre,
            addGenre,
            renameGenre,
        }}>
            {children}
        </GenreContext.Provider>
    )
}
