import { Outlet } from 'react-router-dom';
import Box from '@mui/material/Box';
import Header from './Header';
import Navigation from './Navigation';
import Footer from './Footer';
import Container from "@mui/material/Container";

const Layout = () => {
  return (
    <>
      <Header />
      <Navigation />
      <Box component="main">
        <Outlet />
      </Box>
      <Footer />
    </>
  )
}

export default Layout
