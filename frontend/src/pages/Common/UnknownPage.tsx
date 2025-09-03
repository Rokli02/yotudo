import { FC } from 'react'
import { useNavigate } from 'react-router-dom';
import { ArrowBack } from '@mui/icons-material';
import { Box, Button as MuiButton } from '@mui/material';
import { CSSObject, SxProps, Theme } from '@mui/material/styles';

export const UnknownPage: FC = () => {
    const navigate = useNavigate();

    return (
        <Box sx={ContainerStyle}>
            <div className='back-container'>
                <MuiButton sx={ButtonStyle} onClick={() => navigate(-1)}>
                    <ArrowBack />
                    <span>Vissza</span>
                </MuiButton>
            </div>
            <div className='content'>
                <h1>Ilyen oldal nincs, de azért</h1>
                <div className='img-wrapper'>
                    <img alt='Noice' src={`/imgs/noice.png`}/>
                </div>
            </div>
        </Box>
    )
}

export default UnknownPage;

const ContainerStyle: SxProps<Theme> = {
    width: '100vw',
    position: 'relative',
    display: 'flex',
    paddingTop: 8,
    paddingInline: 6,
    flexDirection: 'column',
    justifyContent: 'center',
    '& .back-container': {
        minWidth: 'fit-content',
        width: 'min(200px, 90%)',
    } as CSSObject,
    '& .content': {
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        flexDirection: 'column',
        '& h1': {
            textAlign: 'center',
            fontSize: '2.5rem',
            marginBottom: 8,
            letterSpacing: 2,
            '@media screen and (max-width: 1000px)': {
                fontSize: '1.8rem',
                letterSpacing: 1,
            } as CSSObject,
            '@media screen and (max-width: 500px)': {
                fontSize: '1.3rem',
                letterSpacing: 'initial',
            } as CSSObject,
        } as CSSObject,
        '& div.img-wrapper': {
            width: '90%',
            maxWidth: '500px',
            aspectRatio: 1,
            position: 'relative',
            display: 'inline-flex',
            backgroundImage: 'radial-gradient(white 0%, transparent 70%, #f5db6788 70%, transparent 71%)',
            justifyContent: 'center',
            alignItems: 'center',
            '& > img': {
                width: '50%',
                aspectRatio: .84,
                position: 'relative',
            } as CSSObject,
        } as CSSObject,
    } as CSSObject,
};

const ButtonStyle = {
    width: '100%',
    fontSize: 20,
    columnGap: 10,
    color: 'crimson',
    borderRadius: 16,
    '@media screen and (max-width: 1000px)': {
        fontSize: 16,
        columnGap: 6,
    } as CSSObject,
    '& svg': {
        height: 24,
        aspectRatio: 1,
        '@media screen and (max-width: 1000px)': {
            height: 20,
        } as CSSObject,
    } as CSSObject,
};