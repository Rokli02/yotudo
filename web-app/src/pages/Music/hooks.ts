import { ComponentProps, useCallback, useEffect, useMemo, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { PaginationComponent } from "./Pagination.component.js";
import { DefaultPageSize } from "./constants.js";
import { Music, MusicService, Pagination } from "@src/api/index.js";

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

    const onHoldMusicItem = useCallback((id: number) => {
        router(`${id}`, { relative: 'path' });
    }, [router]);

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
        onHoldMusicItem,
    } as const
}
