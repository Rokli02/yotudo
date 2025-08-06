export interface Genre {
    id: number;
    name: string;
}

export interface NewGenre extends Omit<Genre, 'id'> {}

export interface UpdateGenre extends Pick<Genre, 'name'> {}

export interface Status {
    id: number;
    name: string;
    description: string;
  }
  