import { Box, SxProps, Theme } from "@mui/material";
import { FC } from "react";
import { MusicItems } from "./MusicItems.component.js";
import { Searchbar } from "@src/components/form/index.js";
import { PaginationComponent } from "./Pagination.component.js";
import { useMusicPageState } from "./hooks.js";

export const MusicsPage: FC = () => {
    const {
        musics,
        filter,
        page,
        onSearchDebounce,
        onPaginationChange,
        holdMusicItem,
        downloadMusic,
    } = useMusicPageState();
    const count = musics?.count ?? 0

    return (
        <Box sx={ContainerStyle}>
            <Searchbar value={filter} sx={SearchbarStyle} onDebounce={onSearchDebounce} debounceTime={400} />
            <PaginationComponent page={page} count={count} onChange={onPaginationChange} />
            <MusicItems musics={musics ? musics.data : null} holdItem={holdMusicItem} downloadMusic={downloadMusic}/>
        </Box>
    )
};

const ContainerStyle: SxProps<Theme> = {
    position: 'relative',
    padding: '32px 16px',
}

const SearchbarStyle: SxProps<Theme> = {
    maxWidth: 600,
    width: '90%',
    marginInline: 'auto',
}