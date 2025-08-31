import { model } from "@wailsjs/go/models";
import { Music, NewMusic, MusicUpdate } from "../models/Music";
import { Page, Pagination } from "../models/Page";
import { GetManyByPagination, GetById, Save, Update } from '@service/MusicService';
import { DownloadByMusicId, MoveToDownloadDir } from '@service/YoutubeService';
import { GetAllStatus } from "./status.service";
import { Status } from "../models/Misc";

export async function GetMusics(page: Page = { page: 0, size: 25 }, statusId: number = -1): Promise<Pagination<Music[]>> {
    const result = await GetManyByPagination(page.filter ?? '', statusId, { Page: page.page, Size: page.size }, [{ Key: 'updated_at', Dir: -1 }]);
    const statusMap = await GetAllStatus();

    return {
        data: result.Data.map((m) => convertGoMusicToTsMusic(m, statusMap)),
        count: result.Count,
    }
}

export async function GetMusicById(id: number): Promise<Music> {
    const result = await GetById(id);
    const statusMap = await GetAllStatus();

    return convertGoMusicToTsMusic(result, statusMap)
}

export async function SaveMusic(newMusic: NewMusic): Promise<Music | null> {
    const response = await Save(new model.NewMusic({
        Name: newMusic.name,
        Author: { Id: newMusic.author.id, Name: newMusic.author.label },
        Contributors: newMusic.contributor.map((c) => ({ Id: c.id, Name: c.label })),
        Url: newMusic.url,
        Album: newMusic.album,
        GenreId: newMusic.genre.id,
        Published: newMusic.published,
        PicFilename: newMusic.picName,
        PicType: newMusic.picName_chosenType,
    } as model.NewMusic));

    const statusMap = await GetAllStatus();

    return convertGoMusicToTsMusic(response, statusMap)
}

export async function UpdateMusic(newMusic: MusicUpdate): Promise<Music | null> {
    const response = await Update(new model.UpdateMusic({
        Id: newMusic.id,
        Name: newMusic.name,
        Author: { Id: newMusic.author.id, Name: newMusic.author.label },
        Contributors: newMusic.contributor.map((c) => ({ Id: c.id, Name: c.label })),
        Url: newMusic.url,
        Album: newMusic.album,
        GenreId: newMusic.genre.id,
        Published: newMusic.published,
        Status: newMusic.status.id,
        PicFilename: newMusic.picName,
        PicType: newMusic.picName_chosenType,
    }))

    const statusMap = await GetAllStatus();

    return convertGoMusicToTsMusic(response, statusMap)
}

export async function DeleteMusic(id: number): Promise<boolean> {
    return Promise.resolve(false);
}

export async function MoveMusicTo(id: number): Promise<void> {
    return MoveToDownloadDir(id);
}

/**
 * Letölti a 'Music' modellhez kötődő zenét
 * @param id Music ID
 * @param eventName Event name, amin keresztül érkezik az aszinkron esemény információ
 * @throws MusicNotFoundError
 */
export async function DownloadMusic(id: number, eventName: string): Promise<void> {
    return DownloadByMusicId(id, eventName);
}

function convertGoMusicToTsMusic(music: model.Music, status: Status[]): Music {
    return {
        id: music.Id,
        name: music.Name,
        genre: {
            id: music.Genre.Id,
            name: music.Genre.Name,
        },
        album: music.Album,
        url: music.Url,
        status: status[music.Status],
        published: music.Published,
        author: {
            id: music.Author.Id,
            name: music.Author.Name,
        },
        contributor: music.Contributors.map((c) => ({ id: c.Id, name: c.Name })),
        picName: music.PicFilename,
    }
}