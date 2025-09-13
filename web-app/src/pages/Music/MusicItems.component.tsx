import { Box, SxProps, Theme } from "@mui/material";
import { Music } from "@src/api/index.js";
import { FC } from "react";
import { LoadingPage } from "../Common/LoadingPage.js";
import { MusicItem } from "./MusicItem.component.js";
import { MusicService } from '@src/api/index.js'

export interface MusicItemsProps {
    musics: Music[] | null;
    onHoldItem: (id: number) => void;
}

export const MusicItems: FC<MusicItemsProps> = ({musics, onHoldItem}) => {
    if (musics === null) {
        return <Box sx={LoadingWrapperStyle}><LoadingPage size="large" /></Box>
    }

    if (musics.length === 0) {
        return <Box sx={NoContentStyle}>Nem találhatóak zenék a keresési feltételek alapján!</Box>
    }

    return <Box sx={ContentStyle}>
        { musics.map((music) => 
            <MusicItem
                key={music.id}
                music={music}
                onAction={() => MusicService.DownloadMusic(music.id)}
                onActionAfterHold={() =>  onHoldItem(music.id)}
            />)
        }
    </Box>
}

const LoadingWrapperStyle: SxProps<Theme> = { width: 'fit-content', marginTop: '48px', marginInline: 'auto' }

const ContentStyle: SxProps<Theme> = {
    display: 'flex',
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: '.5rem .75rem',
    justifyContent: 'center',
    alignItems: 'flex-start',
};

const NoContentStyle: SxProps<Theme> = {
    marginTop: '36px',
    textAlign: 'center',
    fontSize: '1.3rem',
    fontWeight: 500,
}