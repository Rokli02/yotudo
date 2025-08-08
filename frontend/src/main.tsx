import { StrictMode } from 'react'
import { RouterProvider , createBrowserRouter } from 'react-router-dom'
import {createRoot} from 'react-dom/client'
import { routes } from '@src/layout/Layout'

const router = createBrowserRouter(routes)

// Github ajánlás, ha netalán nem töltené be a 'go' runtime-ot
if (!('go' in window)) location.replace('/')

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <RouterProvider router={router} />
    </StrictMode>
)
