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
        _setPage((pre) => {
            const newState = {
                ...pre,
                ...pageUpdate,
            }
            
            onChange?.(newState, ...extraArgs)
            
            return newState;
        })

        // Talán 'await'-elni kéne rá, de most jó lesz így
    }

    return [page, setPage, _setPage];
}