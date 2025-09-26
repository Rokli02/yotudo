import { RouteObject } from 'react-router-dom';
import NavbarLayout from './Navbar/NavbarLayout';
import { UnknownPage } from '@src/pages/Common';
import MusicPage from '@src/pages/Music/MusicPage';
import MiscPage from'@src/pages/Misc/MiscPage';
import AuthorPage from'@src/pages/Author/AuthorPage';
import WebserverPage from '@src/pages/Webserver/WebserverPage';

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
            },
            {
              path: "server",
              element: <WebserverPage />
            }
        ],
    },
    {
        path: "*",
        Component: UnknownPage,
    }
];
