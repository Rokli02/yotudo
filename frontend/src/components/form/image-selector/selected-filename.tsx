import { Close, Image } from "@mui/icons-material";
import { IconButton, SxProps, Theme, Typography } from "@mui/material";
import { Box } from "@mui/system";
import { FC } from "react";

export interface SelectedFilenameProps {
    value: string;
    onRemove: () => void;
}

export const SelectedFilename: FC<SelectedFilenameProps> = ({value, onRemove}) => {
    return (
        <Box sx={SelectedFilenameStyle.Wrapper}>
            <Image sx={SelectedFilenameStyle.Image} />
            <Box sx={SelectedFilenameStyle.TextWrapper}>
                <Typography sx={SelectedFilenameStyle.Text}>{value}</Typography>
            </Box>
            <IconButton sx={SelectedFilenameStyle.UnselectButton} onClick={onRemove} title="Eltávolítás">
                <Close />
            </IconButton>
        </Box>
    )
}

const SelectedFilenameStyle = {
    Wrapper: {
        display: 'inline-grid',
        width: '100%',
        gridTemplateColumns: 'min-content auto min-content',
        alignItems: 'center',
        padding: '4px',
        columnGap: '8px',
    },
    Image: {
        width: 28,
        height: 28,
        color: 'var(--primary-color)',
    },
    TextWrapper: {
        overflow: 'hidden',
        whiteSpace: 'nowrap',
    },
    Text: {
        maxWidth: '100%',
        width: 'fit-content',
        lineHeight: '100%',
        fontSize: 16,
        backgroundColor: 'var(--primary-color)',
        borderRadius: '4px',
        padding: '5px 16px',
        color: 'var(--font-color)',
        textOverflow: 'ellipsis',
        overflow: 'hidden',
        whiteSpace: 'nowrap',
    },
    UnselectButton: {
        color: 'var(--font-color)',
        padding: '4px',
        width: 30,
        height: 30,
        '>svg': {
            width: '100%',
            height: '100%',
        }
    },
} as const satisfies Record<string, SxProps<Theme>>;