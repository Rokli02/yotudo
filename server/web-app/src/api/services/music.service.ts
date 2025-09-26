import { Music, MusicDTO } from "../models/Music.js";
import { Page, Pagination, PaginationDTO } from "../models/Page.js";
import { customFetch } from "../customFetch.js";
import { GetAllStatus } from "./status.service.js";
import { Status } from "../models/Misc.js";

export async function GetMusics(page: Page = { page: 0, size: 10 }): Promise<Pagination<Music[]>> {
    const result = await customFetch<PaginationDTO<MusicDTO[]>>("/musics", {
        method: "GET",
        query: {
            "f": page.filter ?? '',
            "p": page.page+"",
            "ps": page.size+"",
        },
    }).catch(() => ({ Data: [], Count: 0 }))

    const statusMap = await GetAllStatus();

    return {
        data: result.Data.map((m) => convertGoMusicToTsMusic(m, statusMap)),
        count: result.Count,
    }
}

export async function GetMusicById(id: number): Promise<Music> {
    const result = await customFetch<MusicDTO>(`/musics/${id}`)
    const statusMap = await GetAllStatus();

    return convertGoMusicToTsMusic(result, statusMap)
}

/**
 * Letölti a 'Music' modellhez kötődő zenét
 * @param id Music ID
 * @param eventName Event name, amin keresztül érkezik az aszinkron esemény információ
 * @throws MusicNotFoundError
 */
export async function DownloadMusic(id: number): Promise<void> {
    await customFetch(`/musics/${id}/download`, { method: 'GET', rawResponse: true })
        .then(async (response) => {
            let filename: string;
            const b64filename = response.headers.get('x-file-name')

            if (b64filename) {
                const b64bytes = new Uint8Array([...atob(b64filename)].map(char => char.charCodeAt(0)));
                filename = new TextDecoder('utf-8').decode(b64bytes)
            } else {
                filename = `Downloaded_${Date.now().toString(16)}.mp3`
            }

            const blob = await response.blob();
            const url = window.URL.createObjectURL(blob);

            const tempElement = document.createElement('a');

            tempElement.style.display = 'none';
            tempElement.href = url;
            tempElement.download = filename;

            document.body.appendChild(tempElement);
            tempElement.click();
            document.body.removeChild(tempElement);

            window.URL.revokeObjectURL(url);
        })
        .catch((err) => console.error(err));
}

function convertGoMusicToTsMusic(music: MusicDTO, status: Status[]): Music {
    return {
        id: music.Id,
        name: music.Name,
        genre: {
            id: music.Genre.Id,
            name: music.Genre.Name,
        },
        album: music.Album ?? undefined,
        url: music.Url,
        status: status[music.Status],
        published: music.Published ?? undefined,
        author: {
            id: music.Author.Id,
            name: music.Author.Name,
        },
        contributor: music.Contributors.map((c) => ({ id: c.Id, name: c.Name })),
        picName: music.Image?.Name,
    }
}