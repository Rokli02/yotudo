import { Author, NewAuthor } from "../models/Author";
import { Page, Pagination } from "../models/Page";
import { GetManyByPagination, Save, Delete } from '@controller/AuthorController'

export async function GetAuthors(page: Page = { page: 0, size: 25 }): Promise<Pagination<Array<Author>>> {
    const response = await GetManyByPagination(page.filter ?? '', { Page: page.page, Size: page.size }, [{ Key: "id", Dir: -1 }]);

    return {
        count: response.Count,
        data: response.Data.map((a) => ({ id: a.Id, name: a.Name })),
    }
}

export async function SaveAuthor(newAuthor: NewAuthor): Promise<Author | null> {
    const response = await Save(newAuthor.name);

    return response ? { id: response.Id, name: response.Name } : null;
}

export async function DeleteAuthor(authorId: number): Promise<boolean> {
    return await Delete(authorId);
}