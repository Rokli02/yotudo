import { CachedRounded, CloudDownload, SvgIconComponent, DoNotDisturb, Download } from "@mui/icons-material";
import styled from "styled-components";

export const StatusIcon: Record<number, SvgIconComponent> = {
    '-1': DoNotDisturb,
    0: CloudDownload,
    1: CachedRounded,
    2: Download,
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
    0: CloudDownload,
    1: AnimatedCachedRounded,
    2: Download,
}