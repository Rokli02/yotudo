import { Status } from "../models/Misc";
import { GetAll } from "@service/StatusService"

let statuses: {
    array: Array<Status>,
    object: Record<number, Status>,
} = {
    array: [],
    object: {},
}

export async function GetAllStatus<T extends boolean = false>(asMap: T = false as T): Promise<T extends true ? Record<number, Status> : Array<Status>> {
    if (!statuses.array.length) {
        statuses.array = (await GetAll()).map((status) => ({
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