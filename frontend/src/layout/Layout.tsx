import { FC, lazy, Suspense } from 'react'
import { Route, Routes } from 'react-router-dom';
import NavbarLayout from './Navbar/NavbarLayout';
import { LoadingPage, UnknownPage } from '../pages/Common';
// const MusicPageLazy = lazy(() => import('../pages/Music/MusicPage'))
// const MiscPageLazy = lazy(() => import('../pages/Misc/MiscPage'))
// const AuthorPageLazy = lazy(() => import('../pages/Author/AuthorPage'))

export const Layout: FC = () => {
  return (
    <>
      <Routes>
        <Route Component={NavbarLayout}>
          {/* <Route path="" element={<Suspense fallback={<LoadingPage />}><MusicPageLazy /></Suspense>}/>
          <Route path="misc" element={<Suspense fallback={<LoadingPage />}><MiscPageLazy /></Suspense>} />
          <Route path="author" element={<Suspense fallback={<LoadingPage />}><AuthorPageLazy /></Suspense>} /> */}
        </Route>
        <Route path="*" element={<UnknownPage />}/>
        {/* <Route path='*' element={<Navigate to="/" replace={true}/>}/> */}
      </Routes>
    </>
  )
}

export default Layout;

// const MusicPageLazy = lazy(() => import('../pages/Music/MusicPage'))
// const MiscPageLazy = lazy(() => import('../pages/Misc/MiscPage'))
// const AuthorPageLazy = lazy(() => import('../pages/Author/AuthorPage'))

// <Routes>
//     <Route Component={NavbarLayout}>
//         <Route path="" element={<Suspense fallback={<LoadingPage />}><MusicPageLazy /></Suspense>}/>
//         <Route path="misc" element={<Suspense fallback={<LoadingPage />}><MiscPageLazy /></Suspense>} />
//         <Route path="author" element={<Suspense fallback={<LoadingPage />}><AuthorPageLazy /></Suspense>} />
//     </Route>
//     <Route path="*" element={<UnknownPage />}/>
//     {/* <Route path='*' element={<Navigate to="/" replace={true}/>}/> */}
// </Routes>