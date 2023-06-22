import { createBrowserRouter } from 'react-router-dom';
import App from './App';
import About from './pages/About';
import Layout from './components/Layout' ;

const router = createBrowserRouter([
  {
    element: <Layout />,
    children: [
      {
        path: "/",
        element: <App />,
      },
      {
        path: '/about',
        element: <About />,
      },
    ],
  },
]);

export default router;
