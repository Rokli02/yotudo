export interface Author {
    id: number;
    name: string;
}

export interface NewAuthor extends Omit<Author, 'id'> {}

export interface AuthorFilter {
    search?: string;
    size: number;
    skip: number;
}