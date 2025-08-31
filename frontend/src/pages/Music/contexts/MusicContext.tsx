/* eslint-disable no-case-declarations */
import { createContext, FC, ReactElement, useEffect, useState } from "react";
import { Music, MusicService, NewMusic, Page, Pagination, MusicUpdate, StatusService, Status } from "@src/api";
import { PageSetter, usePage } from "@src/hooks/usePage";
import { EventsOn } from "@wailsjs/runtime/runtime";

export interface IMusicContext {
    musics: Pagination<Music[]>;
    page: Page;
    setPage: PageSetter<[number]>,
    addMusic: (music: NewMusic) => Promise<boolean>;
    modifyMusic: (musicToUpdate: MusicUpdate, index?: number) => Promise<boolean>;
    performAction: (music: Music) => Promise<void>;
}

const PAGE_SIZE = 24;
const MUSIC_STATUS_EVENT_NAME = 'download-progress'

export const MusicContext = createContext<IMusicContext>(null as unknown as IMusicContext);

export const MusicProvider: FC<{ children: ReactElement | ReactElement[] }> = ({ children }) => {
    const [musics, setMusics] = useState<Pagination<Music[]>>({ data: [], count: 0 });
    const [page, setPage, _setPage] = usePage<[number]>(PAGE_SIZE, (state, status) => MusicService.GetMusics(state, status).then(setMusics));
    const [status, setStatus] = useState<Status[]>([])

    async function addMusic(music: NewMusic): Promise<boolean> {
        return await MusicService.SaveMusic(music).then((value) => {
            if (!value) return false

            setMusics((pre) => {
                if (pre.data.unshift(value) > page.size) {
                    pre.data.pop();
                }

                pre.count++;
                
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
                MusicService.DownloadMusic(music.id, MUSIC_STATUS_EVENT_NAME);

                setMusics((pre) => ({
                    count: pre.count,
                    data: pre.data.map((_music) => {
                        if (_music.id !== music.id || _music.status.id == status[1].id) return _music;

                        return { ..._music, status: status[1] };
                    }),   
                }));

                return;
            case 2:
                setMusics((pre) => ({
                    count: pre.count,
                    data: pre.data.map((_music) => {
                        if (_music.id !== music.id || _music.status.id == status[1].id) return _music;

                        return { ..._music, status: status[1] };
                    }),   
                }));
                MusicService.MoveMusicTo(music.id).finally(() => {
                    setMusics((pre) => ({
                        count: pre.count,
                        data: pre.data.map((_music) => {
                            if (_music.id !== music.id || _music.status.id == status[2].id) return _music;

                            return { ..._music, status: status[2] };
                        }),   
                    }));
                });

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
        MusicService.GetMusics(page, -1).then(setMusics);
        // The chance of this line causing any error is extremely low, so I just don't care about it
        const status: Status[] = []
        StatusService.GetAllStatus().then((s) => {
            status.push(...s);
            setStatus(s)
        });

        const cancelEvent = EventsOn(MUSIC_STATUS_EVENT_NAME, function([musicId, progress, stat, err]: [number, number, string, string?]) {
            switch (stat) {
                case 'start':
                    setMusics((pre) => ({
                        count: pre.count,
                        data: pre.data.map((music) => {
                            if (music.id !== musicId || music.status.id == status[1].id) return music;

                            return { ...music, status: status[1] };
                        }),   
                    }));

                    break;
                case 'downloading':
                    console.log(`Download progress for id=${musicId} is ${progress}%`)
                    break;
                case 'completed':
                    MusicService.GetMusicById(musicId).then((m) => {
                        setMusics((pre) => ({
                            count: pre.count,
                            data: pre.data.map((music) => {
                                if (music.id !== musicId || music.status.id == status[2].id) return music;
    
                                return m;
                            }),   
                        }));
                    })

                    // TODO: Ez a kod visszaállítása, amint a thumbnail letöltés meg lett mókolva
                    // setMusics((pre) => ({
                    //     count: pre.count,
                    //     data: pre.data.map((music) => {
                    //         if (music.id !== musicId || music.status.id == status[2].id) return music;

                    //         return { ...music, status: status[2] };
                    //     }),   
                    // }));

                    break;
                case 'failed':
                    setMusics((pre) => ({
                        count: pre.count,
                        data: pre.data.map((music) => {
                            if (music.id !== musicId || music.status.id == status[0].id) return music;

                            return { ...music, status: status[0] };
                        }),   
                    }));

                    console.error(`Failed to download music with id=${musicId}. Reason: "${err ?? 'Unknown'}"`);

                    break;
            }
        })

        return () => {
            cancelEvent()
        }
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