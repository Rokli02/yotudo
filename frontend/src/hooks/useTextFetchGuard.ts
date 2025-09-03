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

    private setProperties(text: string, shouldFetch: boolean): boolean {
        this.previousText = text;
        return this.shouldFetch = shouldFetch
    }

    makeItWorthFetching() {
        this.shouldFetch = true;
        this.previousText = undefined;
    }

    worthFetching(text: string, shouldFetchNextTime?: boolean): boolean {
        if (shouldFetchNextTime !== undefined) return this.setProperties(text, shouldFetchNextTime);
        if (this.previousText === undefined) return this.setProperties(text, true);
        if (this.previousText === text) return this.setProperties(this.previousText, false);
        if (!this.shouldFetch && text.substring(0, this.previousText.length) === this.previousText) return this.setProperties(this.previousText, false);

        return this.setProperties(text, true)
    }
}