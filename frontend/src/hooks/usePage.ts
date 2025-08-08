import { Page } from "@src/api";
import { Dispatch, SetStateAction, useState } from "react";

export type PageSetter<ExtraArgs extends Array<unknown> = []> = (pageUpdate: Partial<Page>, ...extraArgs: ExtraArgs) => void;

export function usePage<ExtraArgs extends Array<unknown> = []>(
    initPageSize: number,
    onChange?: (state: Page, ...extraArgs: ExtraArgs) => unknown,
): [
    Page,
    PageSetter<ExtraArgs>,
    Dispatch<SetStateAction<Page>>,
] {
    const [page, _setPage] = useState<Page>({ page: 0, size: initPageSize });

    function setPage(pageUpdate: Partial<Page>, ...extraArgs: ExtraArgs) {
        // const modifiedKeys = Object.entries(pageUpdate) as Array<[keyof Page, unknown]>
        // if (
        //     modifiedKeys.length === 0 ||
        //     modifiedKeys.every(([key, value]) => page[key] === value)
        // ) {
        //     console.warn("Unnecessary state update was blocked, go find out what caused it")

        //     return;
        // }

        _setPage((pre) => {
            const newState = {
                ...pre,
                ...pageUpdate,
            }
            
            onChange?.(newState, ...extraArgs)
            
            return newState;
        })
        
        //TODO: Delete ha működik minden továbbra
        // const mdfky = {...page, ...pageUpdate};
        // onChange?.(mdfky, ...extraArgs)

        // Talán 'await'-elni kéne rá, de most jó lesz így
    }

    return [page, setPage, _setPage];
}