import { styled, CSSObject } from "@mui/material/styles";

export const Container = styled('div')({
    position: 'relative',
    backgroundColor: 'var(--primary-color)',
    minWidth: '700px',
    maxWidth: '1000px',
    padding: '8px 12px',
    borderRadius: '8px',
    boxShadow: '3px 3px 6px #0004',
    color: 'inherit',
    '@media screen and (max-width: 1000px)': {
        minWidth: '400px',
        maxWidth: '700px',
    } as CSSObject,
})

export const ItemContainer = styled('div')({
    width: '100%',
    display: 'flex',
    flexWrap: 'nowrap',
    backgroundColor: '#0003',
    height: 'fit-content',
    padding: '.5rem .8rem',
    borderRadius: '8px',
    columnGap: '1rem',
    overflowX: 'clip',
    '&.min_width': {
        width: 'max-content',
    },
    '&.data_status': {
        alignItems: 'center',
        '& div.text-div': {
            display: 'inline',
            fontSize: '1rem',
            color: 'inherit',
            cursor: 'default',
            '&.no-wrap': {
                textWrap: 'nowrap',
            } as CSSObject,
        } as CSSObject,
        '& hr.MuiDivider-root.MuiDivider-fullWidth': {
            backgroundColor: 'black',
            marginBlock: '6px',
            height: '100%',
            '&.MuiDivider-vertical': {
                height: 30,
            } as CSSObject,
        } as CSSObject,
    } as CSSObject,
})

export const Content = styled('div')({
    display: 'flex',
    rowGap: '.45rem',
    columnGap: '.3rem',
    alignItems: 'center',
    flexGrow: 1,
    '&.dir_row': {
        flexDirection: 'row'
    } as CSSObject,
    '&.dir_col': {
        flexDirection: 'column'
    } as CSSObject,
})

export const TopLeftIdText = styled('label')({
    position: 'absolute',
    top: '.5rem',
    left: '1rem',
    color: '#fff8',
    fontSize: '.8rem',
})
