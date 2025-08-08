import { Music, NewMusic } from "../models/Music";
import { Page, Pagination } from "../models/Page";

export async function GetMusics(page: Page = { page: 0, size: 25 }, statusId: number = 0): Promise<Pagination<Music[]>> {
    return Promise.resolve({ data: [], count: 0 });
}

export async function SaveMusic(newMusic: NewMusic): Promise<Music | null> {
    return Promise.resolve(null);
}

export async function UpdateMusic(newMusic: NewMusic): Promise<Music | null> {
    return Promise.resolve(null);
}

export async function DeleteMusic(id: number): Promise<boolean> {
    return Promise.resolve(false);
}

export async function ProcessMusic(id: number): Promise<Music | null> {
    return Promise.resolve(null);
}

export async function DownloadMusic(id: number): Promise<boolean> {
    return Promise.resolve(false);
}