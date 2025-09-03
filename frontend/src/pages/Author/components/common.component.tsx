import { SxProps, Theme } from "@mui/material/styles";

export const DialogStyle: SxProps<Theme> = {
    '& .MuiDialog-container': {
        '& .MuiPaper-root': {
            maxWidth: '650px',
            width: '100%',
            backgroundColor: 'var(--background-color)',
            color: 'var(--font-color)',
            '& .MuiDialogContent-root': {
                width: '100%',
                '& > .MuiFormLabel-root': {
                    marginLeft: '6ch',
                },
                '& > .MuiFormControl-root > .MuiSlider-root': {
                    marginLeft: '21px',
                },
            },
        },
    },
    '& .form_items': {
        display: 'flex',
        flexDirection: 'column',
        rowGap: '1rem',
    },
};

export const DialogTitleStyle: SxProps<Theme> = {
    fontSize: '1.25rem',
    marginInline: 'auto',
    width: 'max-content',
};
