export interface Pagination<T = unknown> {
    data: T;
    count: number;
}

export interface PaginationDTO<T = unknown> {
    Data: T;
    Count: number;
}

export interface Page {
    filter?: string;
    page: number;
    size: number;
}

export interface Sort {
    key: string;
    dir: number;
}