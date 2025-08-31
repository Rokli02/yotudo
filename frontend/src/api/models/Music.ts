import { AutocompleteOptions } from "@src/contexts/form";
import { Author } from "./Author";
import { Genre, ID, Status } from "./Misc";

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

export interface NewMusic {
    name: string;
    published: number;
    url: string;
    album: string;
    genre: AutocompleteOptions;
    author: AutocompleteOptions;
    contributor: AutocompleteOptions[];
    picName?: string;
    picName_chosenType?: string;
}

export interface MusicUpdate extends NewMusic, ID {
    status: Status;
}

export interface MusicFilter {
    search?: string;
    status?: number;
    size: number;
    skip: number;
}