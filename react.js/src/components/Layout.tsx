import { Outlet } from 'react-router-dom';
import Header from './Header';
import Navigation from './Navigation';
import Footer from './Footer';

const Layout = () => {
  return (
    <>
      <Header />
      <Navigation />
      <Outlet />
      <Footer />
    </>
  )
}

export default Layout
