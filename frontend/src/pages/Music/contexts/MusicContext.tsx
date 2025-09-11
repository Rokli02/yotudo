/* eslint-disable no-case-declarations */
import { createContext, FC, ReactElement, useEffect, useState } from "react";
import { Music, MusicService, NewMusic, Page, Pagination, MusicUpdate, StatusService, Status, DialogService } from "@src/api";
import { PageSetter, usePage } from "@src/hooks/usePage";
import { EventsOn } from "@wailsjs/runtime/runtime";

export interface IMusicContext {
    musics: Pagination<Music[]>;
    page: Page;
    currentStatus: Status,
    setCurrentStatus: (status: Status) => void;
    setPage: PageSetter,
    addMusic: (music: NewMusic) => Promise<boolean>;
    modifyMusic: (musicToUpdate: MusicUpdate, index?: number) => Promise<boolean>;
    deleteMusic(music: Music): Promise<boolean>;
    performAction: (music: Music) => Promise<void>;
}

const PAGE_SIZE = 16;
const MUSIC_STATUS_EVENT_NAME = 'download-progress'

export const MusicContext = createContext<IMusicContext>(null as unknown as IMusicContext);

export const MusicProvider: FC<{ children: ReactElement | ReactElement[] }> = ({ children }) => {
    const [musics, setMusics] = useState<Pagination<Music[]>>({ data: [], count: 0 });
    const [currentStatus, _setCurrentStatus] = useState<Status>({id: -1, name: 'Nincs', description: ''})
    const [page, setPage, _setPage] = usePage(PAGE_SIZE, (state) => MusicService.GetMusics(state, currentStatus.id).then(setMusics));
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

    async function deleteMusic(music: Music): Promise<boolean> {
        return DialogService.OpenConfirmationDialog("Zene törlés", "Biztosan törlöd a megnyitott zenét?")
            .then(async (res) => {
                if (!res) return false;
               
                const deleteResult = await MusicService.DeleteMusic(music.id);
                if (!deleteResult) return deleteResult;

                return MusicService.GetMusics(page, currentStatus.id).then(setMusics).then(() => true).catch(() => false);
            })
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

    function setCurrentStatus(status: Status) {
        _setCurrentStatus(status)

        MusicService.GetMusics(page, status.id).then(setMusics)
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
                    setMusics((pre) => ({
                        count: pre.count,
                        data: pre.data.map((music) => {
                            if (music.id !== musicId || music.status.id == status[2].id) return music;

                            return { ...music, status: status[2] };
                        }),   
                    }));

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
            currentStatus,
            setCurrentStatus,
            setPage,
            addMusic,
            modifyMusic,
            deleteMusic,
            performAction,
        }}>
            {children}
        </MusicContext.Provider>
    )
}