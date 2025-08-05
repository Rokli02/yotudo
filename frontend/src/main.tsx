import React from 'react'
import { RouterProvider , createBrowserRouter } from 'react-router-dom'
import {createRoot} from 'react-dom/client'
import { UnknownPage } from './pages/Common'
import NavbarLayout from './layout/Navbar/NavbarLayout'

const router = createBrowserRouter([
    {
        path: "/",
        Component: NavbarLayout,
        children: [
            {
                index: true,
                Component: () => { return <div><h1>Hallo, Test App</h1></div> }
            }
        ],
    },
    {
        path: "*",
        Component: UnknownPage,
    }
])

createRoot(document.getElementById('root')!).render(
    <React.StrictMode>
        <RouterProvider router={router} />
    </React.StrictMode>
)
