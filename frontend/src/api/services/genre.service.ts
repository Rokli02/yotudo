import { Genre, NewGenre, UpdateGenre } from "@src/api/models/Misc";
import { GetAll, Save, Rename } from '@controller/GenreController'
import { DataCache } from './serviceCacher';

const genresCache = new DataCache<Array<Genre>>(undefined, 60);

export async function GetAllGenre(): Promise<Array<Genre>> {
    if (genresCache.data === undefined) {
        genresCache.data = (await GetAll()).map((g) => ({ id: g.Id, name: g.Name }))
    }

    return genresCache.data
}

export async function SaveGenre(newGenre: NewGenre): Promise<Genre> {
    const sg = await Save(newGenre.name);
    const savedGenre: Genre = { id: sg.Id, name: sg.Name };
    const genres = genresCache.data;

    if (genres !== undefined) {
        genres.push(savedGenre)
        genresCache.data = genres
    }
    
    return savedGenre;
}

export async function RenameGenre(id: number, updateGenre: UpdateGenre, index?: number): Promise<Genre> {
    // Rename in DB
    const ug = await Rename(id, updateGenre.name);
    const updatedGenre: Genre = { id: ug.Id, name: ug.Name };
    const genres = genresCache.data;

    // If rename was succesful rename cache
    if (genres !== undefined) {
        if (index === undefined) {
            index = genres.findIndex((v) => v.id === id);
        }

        if (index > 0) {
            genres[index] = updatedGenre;
        }

        genresCache.data = genres
    }

    // return with value
    return updatedGenre;
}