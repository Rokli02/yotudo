import { SvgIconComponent } from "@mui/icons-material";
import { styled, CSSObject } from "@mui/material/styles";
import { CSSProperties, memo } from "react";

export const Title = styled('h2')({
    fontSize: '1.25rem',
    marginInline: 'auto',
    width: 'max-content',
})

export const Divider = styled('hr')<{ dir?: 'horizontal' | 'vertical', length?: string | number }>(({ dir, length = '100%' }) => {
    const resultObj: CSSObject = {
        display: 'inline-flex',
        marginInline: '4px',
        background: '#fff6',
        border: 'none',
        borderRadius: '3px',
        padding: '-6px -12px'
    }

    if (dir === 'horizontal') {
        resultObj.width = length;
        resultObj.height = '2px';
    } else {
        resultObj.width = '2px';
        resultObj.height = length;
    }

    return resultObj
})

const TRCBWrapper = styled('div')({
    cursor: 'pointer',
    position: 'absolute',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    top: '4px',
    right: '6px',
    width: '28px',
    height: '28px',
    borderRadius: '50%',
    ':hover': {
        backgroundColor: '#3333'
    } as CSSObject,
    ':active': {
        backgroundColor: '#1113'
    } as CSSObject,
})

export const TopRightButton = memo(({ style = {}, Icon, onClick }: { style?: CSSProperties, Icon: SvgIconComponent, onClick: () => void }) => {
    return (
        <TRCBWrapper onClick={onClick} style={style}>
            <Icon style={{
                width: '85%',
                height: '85%',
            }}/>
        </TRCBWrapper>
    )
})

export { Chip } from './chip';
export { StatusIcon, StatusActionIcon } from './status';