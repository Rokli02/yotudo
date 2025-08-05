import styles from './NavbarStyle.module.css'
import { FC, useLayoutEffect, useRef } from 'react';
import { Link, useLocation } from 'react-router-dom'
import { AudioFile, Dashboard, InterpreterMode } from '@mui/icons-material';
import { NavbarProps } from './Navbar.interface';

const Navbar: FC<NavbarProps> = () => {
  const navItemsRef = useRef<HTMLDivElement>(null)
  const location = useLocation()

  useLayoutEffect(() => {
    if (navItemsRef.current?.children) {
      for(const child of navItemsRef.current.children) {
        if (child.getAttribute('href') === location.pathname) {
          child.setAttribute('data-toggled', '')

          continue
        }
        
        child.removeAttribute('data-toggled')
      }
    }
  }, [location.pathname])

  return (
    <div id={styles['navbar']}>
      <div ref={navItemsRef} className={`${styles['nav-items']} horizontal-scrollbar`}>
        <Link to="" className={styles['nav-item']}>
          <AudioFile />
          Zenék
        </Link>
        <Link to="misc" className={styles['nav-item']}>
          <Dashboard />
          Misc
        </Link>
        <Link to="author" className={styles['nav-item']}>
          <InterpreterMode />
          Szerzők
        </Link>
      </div>
    </div>
  )
}

export default Navbar;