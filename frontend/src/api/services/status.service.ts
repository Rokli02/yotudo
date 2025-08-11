import { Status } from "../models/Misc";
import { GetAll } from "@controller/StatusController"

let statuses: Array<Status> = []

export async function GetAllStatus<T extends boolean = false>(asMap: T = false as T): Promise<T extends true ? Record<number, Status> : Array<Status>> {
    if (!statuses.length) {
        statuses = (await GetAll()).map((status) => ({
            id: status.Id,
            name: status.Name,
            description: status.Description,
        }))
    }

    return asMap
        ? statuses.reduce((obj, curr) => {
            obj[curr.id] = curr

            return obj;
        }, {} as T extends true ? Record<number, Status> : Array<Status>)
        : statuses;
}