import { CircularProgress as MuiCircularProgress, Container as MuiContainer } from "@mui/material";
import { SxProps, Theme } from "@mui/material/styles"
import { CSSProperties, FC, useEffect, useMemo, useState } from "react";

type Size = 'large' | 'medium' | 'small'

interface ILoadingPage {
    size?: Size,
}

const SizeTable: Record<Size, CSSProperties> = {
    'small': {width: '32px', height: '32px'},
    'medium': {width: '48px', height: '48px'},
    'large': {width: '75px', height: '75px'},
}

export const LoadingPage: FC<ILoadingPage> = ({ size = 'large' }) => {
    const [numOfDots, setNumOfDots] = useState(0)
    
    const dots: string = useMemo(() => {
        return new Array(numOfDots).fill(0).map(() => '.').join('')
    }, [numOfDots])

    useEffect(() => {
      const intervalId = setInterval(() => {
        setNumOfDots((pre) => {
            if (pre > 2) {
                return 0
            }

            return pre + 1
        })
      }, 750)
    
      return () => {
        clearInterval(intervalId)
      }
    }, [])
    

    return (
        <MuiContainer sx={ContainerStyle}>
            <div>
                <h1 className={`${size}_size`}>
                    Loading
                    <div className="dots">{dots}</div>
                </h1>
                <MuiCircularProgress sx={CircularProgressStyle} variant="indeterminate" style={SizeTable[size]} />
            </div>
        </MuiContainer>
    )
}


export const ContainerStyle: SxProps<Theme> = {
    display: 'grid',
    justifyItems: 'center',
    alignContent: 'center',
    height: 'calc(100% - var(--navbar-height))',
    width: 'max-content',
    textAlign: 'center',
    '& div h1': {
        '& .dots': {
            display: 'inline-flex',
            width: '30px',
            textAlign: 'start',
        },
        '&.large_size': {
            fontSize: '2rem',
        },
        '&.medium_size': {
            fontSize: '1.7rem',
        },
        '&.small_size': {
            fontSize: '1.3rem',
        },
    },
};

const CircularProgressStyle: SxProps<Theme> = {
    color: 'inherit',
};