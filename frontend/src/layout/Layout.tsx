import { useEffect, useState } from 'react'
import { RouteObject } from 'react-router-dom';
import NavbarLayout from './Navbar/NavbarLayout';
import { UnknownPage } from '@src/pages/Common';
import { Status, StatusService } from '@src/api';
import MusicPage from '@src/pages/Music/MusicPage';
import MiscPage from'@src/pages/Misc/MiscPage';
import AuthorPage from'@src/pages/Author/AuthorPage';

export const routes: Array<RouteObject> = [
    {
        path: "/",
        Component: NavbarLayout,
        children: [
            {
              index: true,
              element: <MusicPage />
            },
            {
              path: "misc",
              element: <MiscPage />
            },
            {
              path: "author",
              element: <AuthorPage />
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