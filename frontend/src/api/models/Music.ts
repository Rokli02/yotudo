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
    useThumbnail?: boolean;
    status: Status;
    pic_id?: number;
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
    useThumbnail: boolean;
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