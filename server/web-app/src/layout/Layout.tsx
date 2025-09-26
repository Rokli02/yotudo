import { lazy, Suspense } from "react";
import { Navigate, RouteObject } from "react-router-dom";
const UnknownPageLazy = lazy(() => import("@src/pages/Common/UnknownPage.js"))
import { MusicsPage } from "@src/pages/Music/Musics.page.js";
import { LoadingPage } from "@src/pages/Common/LoadingPage.js";

export const routes: Array<RouteObject> = [
    {
        path: '/',
        element: <Navigate to="/musics" replace />
    },
    {
        path: '/musics',
        Component: MusicsPage,
    },
    {
        path: "*",
        element: <Suspense fallback={<LoadingPage />}>
            <UnknownPageLazy />
        </Suspense>,
    }
]