import { lazy, Suspense, useEffect, useState } from 'react'
import { RouteObject } from 'react-router-dom';
import NavbarLayout from './Navbar/NavbarLayout';
import { LoadingPage, UnknownPage } from '@src/pages/Common';
import { Status, StatusService } from '@src/api';
// const MusicPageLazy = lazy(() => import('@src/pages/Music/MusicPage'))
const MiscPageLazy = lazy(() => wait(2000).then(() => import('@src/pages/Misc/MiscPage')))
// const AuthorPageLazy = lazy(() => import('@src/pages/Author/AuthorPage'))
const wait = (time: number) => new Promise((resolve) => setTimeout(resolve, time))
// export const Layout: FC = () => {
//   return (
//     <>
//       <Routes>
//         <Route Component={NavbarLayout}>
//           {/* <Route path="" element={<Suspense fallback={<LoadingPage />}><MusicPageLazy /></Suspense>}/>
//           <Route path="misc" element={<Suspense fallback={<LoadingPage />}><MiscPageLazy /></Suspense>} />
//           <Route path="author" element={<Suspense fallback={<LoadingPage />}><AuthorPageLazy /></Suspense>} /> */}
//         </Route>
//         <Route path="*" element={<UnknownPage />}/>
//         {/* <Route path='*' element={<Navigate to="/" replace={true}/>}/> */}
//       </Routes>
//     </>
//   )
// }

export const routes: Array<RouteObject> = [
    {
        path: "/",
        Component: NavbarLayout,
        children: [
            {
                index: true,
                Component: IndexPage,
            },
            {
              path: "misc",
              element: (<Suspense fallback={<LoadingPage />}>
                  <MiscPageLazy />
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