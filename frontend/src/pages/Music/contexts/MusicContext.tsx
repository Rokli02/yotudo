/* eslint-disable no-case-declarations */
import { createContext, FC, ReactElement, useEffect, useState } from "react";
import { StatusService, Music, MusicService, NewMusic, Page, Status, Pagination, MusicUpdate } from "@src/api";
import { PageSetter, usePage } from "@src/hooks/usePage";

export interface IMusicContext {
    musics: Pagination<Music[]>;
    page: Page;
    setPage: PageSetter<[number]>,
    addMusic: (music: NewMusic) => Promise<boolean>;
    modifyMusic: (musicToUpdate: MusicUpdate, index?: number) => Promise<boolean>;
    performAction: (music: Music) => Promise<void>;
}

const PAGE_SIZE = 5;

export const MusicContext = createContext<IMusicContext>(null as unknown as IMusicContext);

export const MusicProvider: FC<{ children: ReactElement | ReactElement[] }> = ({ children }) => {
    const [musics, setMusics] = useState<Pagination<Music[]>>({ data: [], count: 0 })
    const [statuses, setStatuses] = useState<Status[]>([]);
    const [page, setPage, _setPage] = usePage<[number]>(PAGE_SIZE, (state, status) => MusicService.GetMusics(state, status).then(setMusics));

    async function addMusic(music: NewMusic): Promise<boolean> {
        return await MusicService.SaveMusic(music).then((value) => {
            if (!value) return false

            setMusics((pre) => {
                if (pre.data.unshift(value) > page.size) {
                    pre.count++;
                    pre.data.pop();
                }

                return {...pre};
            });

            return true;
        });
    }

    async function modifyMusic(musicToUpdate: MusicUpdate, index?: number) {
        return MusicService.UpdateMusic(musicToUpdate).then((updatedValue) => {
            if (!updatedValue) return false;

            setMusics((pre) => {
                if (modifyMusicAt(pre.data, updatedValue, index)) {
                    return { ...pre };
                }

                return pre;
            })

            return true;
        });
    }

    async function performAction(music: Music, index?: number) {
        switch (music.status.id) {
            case 0:
                setMusics((pre) => {
                    if (modifyMusicAt(pre.data, { ...music, status: statuses[1] }, index)) {
                        return { ...pre };
                    }

                    return pre;
                })

                const response = await MusicService.ProcessMusic(music.id);

                if (!response) {
                    return setMusics((pre) => {
                        if (modifyMusicAt(pre.data, { ...music, status: statuses[0] }, index)) {
                            return { ...pre };
                        }

                        return pre;
                    })
                }

                return setMusics((pre) => {
                    if (modifyMusicAt(pre.data, response, index)) {
                        return { ...pre };
                    }

                    return pre;
                })
            case 2:
                await MusicService.DownloadMusic(music.id);

                return;
        }
    }

    function modifyMusicAt(source: Music[], music: Music, index?: number): boolean {
        if (!index) {
            for (let i = 0; i < source.length; i++) {
                if (source[i].id === music.id) {
                    source[i] = music;

                    return true;
                }
            }
        } else if (source[index].id === music.id) {
            source[index] = music;

            return true;
        }

        return false;
    }

    useEffect(() => {
        MusicService.GetMusics(page).then((res) => {
            if (!res || (Array.isArray(res) && res.length == 0)) {
                return;
            }

            setMusics(res);
        })

        StatusService.GetAllStatus().then((s) => setStatuses(s))
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    return (
        <MusicContext.Provider value={{
            musics,
            page,
            setPage,
            addMusic,
            modifyMusic,
            performAction,
        }}>
            {children}
        </MusicContext.Provider>
    )
}