let ServerBase = import.meta.env.DEV ? 'localhost:3000' : window.location.host;
if (!ServerBase.startsWith("http://")) {
    ServerBase = 'http://'+ServerBase;
}

export async function customFetch(path: string, opts?: FetchOptions & { rawResponse: true }): Promise<Response>
export async function customFetch<T>(path: string, opts?: FetchOptions): Promise<T>
export async function customFetch<T>(path: string, opts: FetchOptions = { method: "GET", rawResponse: false }): Promise<T> {
        if (!opts.header) {
            opts.header = {}
        }

        if (opts.useToken) {
            const token = sessionStorage.getItem('token');

            if (token) {
                opts.header["Authorization"] = `Bearer ${token}`;
            }
        }

        const url = new URL(`/api${path}`, ServerBase)

        if (opts.query) {
            Object.entries(opts.query).forEach(([key, value]) => {
                if (value === undefined) return

                if (Array.isArray(value)) {
                    value.forEach(v => url.searchParams.append(key, v));
                } else {
                    url.searchParams.set(key, value)
                }
            })
        }

        try {
            const response = await fetch(
                url.toString(),
                {
                    method: opts.method,
                    headers: opts.header,
                    body: (opts as FetchOptionsWithoutBody).body ? JSON.stringify((opts as FetchOptionsWithoutBody).body) : null,
                },
            )

            if (!(response.status < 300 && response.status >= 200)) {
                console.error('Error reason:', await response.text())

                throw new HttpError(`Failed to fetch data (${response.status} - ${response.statusText})`, response.status)
            }

            if (opts.rawResponse) {
                return response as T;
            }

            return await response.json()
        } catch (err) {
            if (err instanceof HttpError) throw err

            console.error("Something terrible happened:", err)
            throw new HttpError("Unknown API error", 500)
        }
    }

export type AbstractFetchOptions = {
    header?: Record<string, string>;
    query?: Record<string, undefined | string | string[]>;
    useToken?: boolean;
    rawResponse?: boolean;
}

export type FetchOptions = (FetchOptionsWithBody | FetchOptionsWithoutBody) & AbstractFetchOptions

export type FetchOptionsWithBody = {
    method: "GET" | "DELETE";
}

export type FetchOptionsWithoutBody = {
    method: "POST" | "PUT";
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    body: Record<string, any> | Array<any>;
}

export class HttpError implements Error {
    name: string;
    message: string;
    status: number;
    stack?: string | undefined;

    constructor(message: string, status: number) {
        this.name = "HttpError";
        this.message = message;
        this.status = status;
    }
}