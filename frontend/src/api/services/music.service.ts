import { model } from "@wailsjs/go/models";
import { Music, NewMusic, MusicUpdate } from "../models/Music";
import { Page, Pagination } from "../models/Page";
import { GetManyByPagination, Save, Update } from '@controller/MusicController';
import { DownloadByMusicId } from '@controller/YtController';
import { GetAllStatus } from "./status.service";
import { Status } from "../models/Misc";

export async function GetMusics(page: Page = { page: 0, size: 25 }, statusId: number = 0): Promise<Pagination<Music[]>> {
    const result = await GetManyByPagination(page.filter ?? '', statusId, { Page: page.page, Size: page.size }, [{ Key: 'updated_at', Dir: -1 }]);
    const statusMap = await GetAllStatus();

    return {
        data: result.Data.map((m) => convertGoMusicToTsMusic(m, statusMap)),
        count: result.Count,
    }
}

export async function SaveMusic(newMusic: NewMusic): Promise<Music | null> {
    const response = await Save(new model.NewMusic({
        Name: newMusic.name,
        Author: { Id: newMusic.author.id, Name: newMusic.author.label },
        Album: newMusic.album,
        GenreId: newMusic.genre.id,
        Url: newMusic.url,
        Published: newMusic.published,
        PicFilename: newMusic.picUri ? newMusic.picUri : newMusic.useThumbnail ? 'thumbnail': "",
        Contributors: newMusic.contributor.map((c) => ({ Id: c.id, Name: c.label })),
    } as model.NewMusic));

    const statusMap = await GetAllStatus(true);

    return convertGoMusicToTsMusic(response, statusMap)
}

export async function UpdateMusic(newMusic: MusicUpdate): Promise<Music | null> {
    const response = await Update(new model.UpdateMusic({
        Id: newMusic.id,
        Name: newMusic.name,
        Author: { Id: newMusic.author.id, Name: newMusic.author.label },
        Album: newMusic.album,
        GenreId: newMusic.genre.id,
        Url: newMusic.url,
        Published: newMusic.published,
        Contributors: newMusic.contributor.map((c) => ({ Id: c.id, Name: c.label })),
        Status: newMusic.status.id,
        PicFilename: newMusic.picUri ? newMusic.picUri : newMusic.useThumbnail ? 'thumbnail': undefined,
    }))

    const statusMap = await GetAllStatus(true);

    return convertGoMusicToTsMusic(response, statusMap)
}

export async function DeleteMusic(id: number): Promise<boolean> {
    return Promise.resolve(false);
}

export async function MoveMusicTo(id: number): Promise<Music | null> {
    return Promise.resolve(null);
}

export async function DownloadMusic(id: number, eventName: string): Promise<void> {
    return DownloadByMusicId(id, eventName);
}

function convertGoMusicToTsMusic(music: model.Music, statusMap: Record<number, Status>) {
    return {
        id: music.Id,
        name: music.Name,
        genre: {
            id: music.Genre.Id,
            name: music.Genre.Name,
        },
        album: music.Album,
        url: music.Url,
        status: statusMap[music.Status],
        published: music.Published,
        author: {
            id: music.Author.Id,
            name: music.Author.Name,
        },
        contributor: music.Contributors.map((c) => ({ id: c.Id, name: c.Name })),
        useThumbnail: undefined,
    }
}