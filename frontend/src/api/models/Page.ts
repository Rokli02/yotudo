export interface Pagination<T = unknown> {
    data: T;
    count: number;
}

export interface Page {
    page: number;
    size: number;
}

export interface Sort {
    key: string;
    dir: number;
}