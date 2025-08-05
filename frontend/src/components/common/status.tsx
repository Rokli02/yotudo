import { CachedRounded, Done, Inventory, SvgIconComponent, DoNotDisturb, Download, SaveAlt } from "@mui/icons-material";
import styled from "styled-components";

export const StatusIcon: Record<number, SvgIconComponent> = {
    '-1': DoNotDisturb,
    0: Inventory,
    1: CachedRounded,
    2: Done,
}

const AnimatedCachedRounded = styled(CachedRounded)({
    '@keyframes spin': {
        '0%': {
            transform: 'rotate(0deg)',
        },
        '100%': {
            transform: 'rotate(360deg)',
        },
    },
    animationName: 'spin',
    animationIterationCount: 'infinite',
    animationDirection: 'reverse',
    animationDuration: '1.5s'
})

export const StatusActionIcon: Record<number, SvgIconComponent> = {
    0: Download,
    1: AnimatedCachedRounded,
    2: SaveAlt,
}