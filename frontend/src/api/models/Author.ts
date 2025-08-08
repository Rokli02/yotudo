export interface Author {
    id: number;
    name: string;
}

export interface NewAuthor extends Omit<Author, 'id'> {}