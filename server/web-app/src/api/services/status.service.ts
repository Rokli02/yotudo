import { customFetch } from "../customFetch.js";
import { Status, StatusDTO } from "../models/Misc.js";

const statuses: {
    array: Array<Status>,
    object: Record<number, Status>,
} = {
    array: [],
    object: {},
}

export async function GetAllStatus<T extends boolean = false>(asMap: T = false as T): Promise<T extends true ? Record<number, Status> : Array<Status>> {
    if (!statuses.array.length) {

        statuses.array = (await customFetch<StatusDTO[]>("/status").catch(() => [])).map((status) => ({
            id: status.Id,
            name: status.Name,
            description: status.Description,
        }))

        statuses.object = statuses.array.reduce((obj, curr) => {
            obj[curr.id] = curr

            return obj;
        }, {} as Record<number, Status>)
    }

    return asMap
        ? statuses.object as (T extends true ? Record<number, Status> : Status[])
        : statuses.array;
}