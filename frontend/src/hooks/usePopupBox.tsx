import { Box, SxProps, Theme } from "@mui/material";
import { MouseEvent, useCallback, useRef } from "react";

/**
 * 
 * @param popupText Szöveg amit meg szeretnénk jeleníteni a kis felugró dobozban.
 * @param action Művelet amit szeretnénk végrehajtani, amikor megtörténik a kattintás.
 * @returns `onClick` Annak az elemnek kell átadni amelyikre történő kattintáskor szeretnénk hogy mgejelenjen a felugró doboz.
 *          `PopupBox` A doboz ami kattintást követően megjelenik. Amellé az elem mellé kell helyezni, amelyiknek átadtuk az `onClick` függvényt.
 */
export function usePopupBox(popupText: string, action?: () => Promise<unknown>) {
    const linkCopyBoxRef = useRef<HTMLDivElement>(null);
    const timeoutIdRef = useRef<NodeJS.Timeout | null>(null);

    async function onClick(e: MouseEvent<HTMLElement>) {
        e.preventDefault();

        await action?.()

        const parent = (e.target as any).offsetParent as HTMLDivElement
        if (parent && linkCopyBoxRef.current) {
            const rect = parent.getBoundingClientRect();
            const offsetX = e.clientX - rect.left;
            const offsetY = e.clientY - rect.top;

            linkCopyBoxRef.current.style.left = `${offsetX}px`;
            linkCopyBoxRef.current.style.top = `${offsetY}px`;
            linkCopyBoxRef.current.toggleAttribute("data-copied", true);

            if (timeoutIdRef.current) clearTimeout(timeoutIdRef.current);

            timeoutIdRef.current = setTimeout(() => {
                linkCopyBoxRef.current?.toggleAttribute("data-copied", false);
            }, 1000);
        }
    }

    const PopupBox = useCallback(() => {
        return <Box sx={LinkCopyBox} ref={linkCopyBoxRef}>{popupText}</Box>
    }, [popupText])

    return {
        onClick,
        PopupBox,
    }
}

export default usePopupBox;

const LinkCopyBox: SxProps<Theme> = {
    position: 'absolute',
    left: '50%',
    top: '100%',
    width: 'fit-content',
    padding: '2px 12px 4px',
    borderRadius: '4px',
    backgroundColor: 'var(--primary-color)',
    userSelect: 'none',
    pointerEvents: 'none',
    transition: 'opacity ease-in-out 200ms',
    opacity: 0,
    '&[data-copied]': {
        opacity: 1,
    }
};
