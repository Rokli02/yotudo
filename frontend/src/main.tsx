import { StrictMode } from 'react'
import { RouterProvider , createBrowserRouter } from 'react-router-dom'
import {createRoot} from 'react-dom/client'
import { routes } from '@src/layout/Layout'

const router = createBrowserRouter(routes)

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <RouterProvider router={router} />
    </StrictMode>
)
