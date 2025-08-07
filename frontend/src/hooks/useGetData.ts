import { useLayoutEffect, useState } from "react"

export const useGetData = <T>(dataFetcher: () => Promise<T>) => {
  const [data, setData] = useState<T>()
  const [loading, setLoading] = useState<boolean>(true)

  useLayoutEffect(() => {
    const fetchPromise = dataFetcher()
    setLoading(true)

    fetchPromise.then((res) => {
      console.log(res)
      if (!res || (Array.isArray(res) && res.length == 0)) {
        return;
      }

      setData(res)
    }).finally(() => {
      setLoading(false)
    })
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return [loading, data] as const;
}