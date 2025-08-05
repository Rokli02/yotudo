import { useState } from "react"

export const useFilter = <T extends object>(onChange: (filter: T) => Promise<void>, init: T = {} as T) => {
    const [filter, _setFilter] = useState<T>(init)

    const setFilter = (name: keyof T, value: unknown) => {
        _setFilter((pre) => {
            if (value === null) {
                delete pre[name];
            } else {
                (pre as Record<keyof T, unknown>)[name] = value;
            }

            onChange(pre)

            return {
                ...pre,
            }
        })
    }

    return { filter, setFilter }
}