import { useRef } from "react"

export interface useLoadingReturn {
    value: boolean,
    startLoading(): void;
    stopLoading(): void;
}

export const useLoading = (initialValue: boolean = false): useLoadingReturn => {
    const loadingRef = useRef<useLoadingReturn>({
        value: initialValue,
        startLoading() { this.value = true },
        stopLoading() { this.value = false },
    })

    return loadingRef.current
}