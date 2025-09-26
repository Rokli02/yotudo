import { Author, AuthorDTO } from "./Author.js";
import { ImageDTO } from "./Image.js";
import { Genre, GenreDTO, Status } from "./Misc.js";

export interface Music {
    id: number;
    name: string;
    published?: number;
    url: string;
    album?: string;
    updatedAt?: string;
    author: Author;
    genre: Genre;
    status: Status;
    picName?: string;
    contributor?: Author[];
}

export interface MusicFilter {
    search?: string;
    status?: number;
    size: number;
    skip: number;
}

export interface MusicDTO {
    Id: number;
    Name: string;
    Published: number | null;
    Album: string | null;
    Url: string;
    Filename: string;
    Image: ImageDTO;
    Status: number;
    Genre: GenreDTO;
    Author: AuthorDTO;
    Contributors: AuthorDTO[];
}