import { ComponentProps, useCallback, useEffect, useMemo, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { PaginationComponent } from "./Pagination.component.js";
import { DefaultPageSize } from "./constants.js";
import { Music, MusicService, Pagination, StatusService } from "@src/api/index.js";

export function useMusicPageState() {
    const [musics, setMusics] = useState<Pagination<Music[]> | null>(null)
    const [searchParams, setSearchParams] = useSearchParams()
    const router = useNavigate()
    const filter = searchParams.get('f') ?? '';
    const p = searchParams.get('p');
    const count = musics?.count ?? 0
    
    const page = useMemo(() => {
        const _page = Number(p);
        return isNaN(_page) ? 0 : count < DefaultPageSize ? 0 : _page;
    }, [p, count]);

    const onSearchDebounce = useCallback((search: string) => {
        setSearchParams((pre) => {
            if (!search) {
                pre.delete('f')
            } else {
                pre.set('f', search)
            }

            return pre.entries().toArray()
        })
    }, [setSearchParams]);

    const onPaginationChange: ComponentProps<typeof PaginationComponent>['onChange'] = useCallback((_, page) => {
        setSearchParams((pre) => {
            if (page == 1) {
                pre.delete('p')
            } else {
                pre.set('p', page-1+'')
            }

            return pre.entries().toArray()
        })
    }, [setSearchParams]);

    const holdMusicItem = useCallback((id: number) => {
        router(`${id}`, { relative: 'path' });
    }, [router]);

    const downloadMusic = useCallback(async (id: number) => {
        const status = await StatusService.GetAllStatus()
        const updateList = (statusId: number) => setMusics((pre) => ({
            ...pre,
            data: pre?.data.map((music) => music.id !== id 
                ? music
                : { ...music, status: status[statusId] }
            ) ?? [],
        } as Pagination<Music[]>))

        updateList(1)

        MusicService.DownloadMusic(id).finally(() => updateList(2))
    }, [setMusics])

    useEffect(() => {
        MusicService.GetMusics({ filter, page, size: DefaultPageSize }).then((musics) => {
            setMusics(musics)
        })
    }, [filter, page]);

    return {
        musics,
        filter,
        page,
        onSearchDebounce,
        onPaginationChange,
        holdMusicItem,
        downloadMusic,
    } as const
}
