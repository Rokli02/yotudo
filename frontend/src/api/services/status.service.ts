import { Status } from "../models/Misc";
import { GetAll } from "@controller/StatusController"

let statuses: Array<Status> = []

export async function GetAllStatus(): Promise<Array<Status>> {
    if (!statuses.length) {
        statuses = (await GetAll()).map((s) => {
            console.log("HEEE?", s)

            return s
        }).map((status) => ({
            id: status.Id,
            name: status.Name,
            description: status.Description,
        }))
    }

    return statuses
}