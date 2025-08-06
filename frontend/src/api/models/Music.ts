import { Author } from "./Author";
import { Genre, Status } from "./Misc";

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

export interface NewMusic extends Omit<Music, 'id'> {}

export interface MusicFilter {
    search?: string;
    status?: number;
    size: number;
    skip: number;
}