import { useRef } from "react";

export function useTextFetchGuard() {
    const fetchGuard = useRef(new FetchGuard());

    return fetchGuard.current;
}

class FetchGuard {
    private shouldFetch!: boolean;
    private previousText?: string;

    constructor() {
        this.makeItWorthFetching()
    }

    setShouldFetch(shouldFetch: boolean) {
        return this.shouldFetch = shouldFetch;
    }

    private setProperties(text: string, shouldFetch: boolean): boolean {
        this.previousText = text;
        return this.shouldFetch = shouldFetch
    }

    makeItWorthFetching() {
        this.shouldFetch = true;
        this.previousText = undefined;
    }

    worthFetching(text: string): boolean {
        if (!this.previousText) return this.setProperties(text, true);
        if (!this.shouldFetch && text.substring(0, this.previousText.length) === this.previousText) return false

        this.previousText = text;

        return true;
    }
}