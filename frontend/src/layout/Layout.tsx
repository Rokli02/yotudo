import { lazy, Suspense, useEffect, useState } from 'react'
import { RouteObject } from 'react-router-dom';
import NavbarLayout from './Navbar/NavbarLayout';
import { LoadingPage, UnknownPage } from '@src/pages/Common';
import { Status, StatusService } from '@src/api';
const MusicPageLazy = lazy(() => import('@src/pages/Music/MusicPage'))
const MiscPageLazy = lazy(() => import('@src/pages/Misc/MiscPage'))
const AuthorPageLazy = lazy(() => import('@src/pages/Author/AuthorPage'))

export const routes: Array<RouteObject> = [
    {
        path: "/",
        Component: NavbarLayout,
        children: [
            {
              index: true,
              element: (<Suspense fallback={<LoadingPage />}>
                  <MusicPageLazy />
              </Suspense>)
            },
            {
              path: "misc",
              element: (<Suspense fallback={<LoadingPage />}>
                  <MiscPageLazy />
              </Suspense>)
            },
            {
              path: "author",
              element: (<Suspense fallback={<LoadingPage />}>
                  <AuthorPageLazy />
              </Suspense>)
            }
        ],
    },
    {
        path: "*",
        Component: UnknownPage,
    }
]

function IndexPage() {
    const [status, setStatus] = useState<Status[]>([])

    useEffect(() => {
      StatusService.GetAllStatus().then((res) => {
        setStatus(res)
      })
    }, [])
    

    return <div>
        <h1>More, Test App</h1>
        <li>
            {status.map(s => (<ul key={s.id}>{s.id} - {s.name}</ul>))}
        </li>
    </div>
}