import { FC } from 'react'
import { Author } from '@src/api'
import { CSSObject, styled } from '@mui/material/styles'
import DeleteIcon from '@mui/icons-material/Delete'
import { IconButton } from '@mui/material'
import { CustomCSS } from '@src/components/common/interface'

export interface AuthorItemProps {
    author: Author
    onDelete: (authorId: number) => Promise<boolean>;
};

export const AuthorItem: FC<AuthorItemProps> = ({ author, onDelete }) => {
  return (
    <Item>
        <div className='main_content'>
            <span className='name_field'>{ author.name }</span>
        </div>
        <div className='delete_btn'>
            <IconButton onClick={() => onDelete(author.id)}>
                <DeleteIcon />
            </IconButton>
        </div>
    </Item>
  )
}

const Item = styled('div')({
    position: 'relative',
    minHeight: 'fit-content',
    width: '100%',
    maxWidth: '750px',
    padding: '4px 18px',
    columnGap: '1rem',
    rowGap: '.5rem',
    display: 'flex',
    flexGrow: 1,
    flexWrap: 'wrap',
    alignItems: 'center',
    justifyContent: 'center',
    borderRadius: '8px',
    boxShadow: '3px 3px 9px #0004',
    backgroundColor: 'var(--primary-color)',
    '& .MuiSlider-root': {
        marginLeft: '8px',
        minWidth: '200px',
        maxWidth: '200px',
        '& .MuiSlider-markLabel': {
            color: '#fff7',
        },
        '@media screen and (max-width: 650px)': {
            minWidth: '200px',
            maxWidth: '80%',
        }
    },
    '& .main_content': {
        display: 'inherit',
        columnGap: '.3rem',
        height: '100%',
        flexGrow: 1,
        alignItems: 'center',
        rowGap: '.5rem',
        '& .id_field': {
            minWidth: 'fit-content',
            width: '6ch',
        },
        '& .name_field': {
            display: 'inherit',
        },
        '@media screen and (max-width: 650px)': {
            flexWrap: 'wrap',
            width: '100%',
            '& .name_field, & .id_field': {
                width: '100%',
                justifyContent: 'center',
                textAlign: 'center',
            },
        },
    },
} as CustomCSS)

export default AuthorItem