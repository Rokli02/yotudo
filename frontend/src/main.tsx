import React, { useEffect, useState } from 'react'
import { RouterProvider , createBrowserRouter } from 'react-router-dom'
import {createRoot} from 'react-dom/client'
import { UnknownPage } from '@src/pages/Common'
import NavbarLayout from '@src/layout/Navbar/NavbarLayout'
import { Status } from '@src/api'
import { GetAllStatus } from '@src/api/services/status.service'

const router = createBrowserRouter([
    {
        path: "/",
        Component: NavbarLayout,
        children: [
            {
                index: true,
                Component: IndexPage,
            }
        ],
    },
    {
        path: "*",
        Component: UnknownPage,
    }
])

function IndexPage() {
    const [status, setStatus] = useState<Status[]>([])

    useEffect(() => {
      GetAllStatus().then((res) => {
        console.log("MORE", res)

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

createRoot(document.getElementById('root')!).render(
    <React.StrictMode>
        <RouterProvider router={router} />
    </React.StrictMode>
)
