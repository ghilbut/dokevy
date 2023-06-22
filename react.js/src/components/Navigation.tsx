import { Link } from 'react-router-dom';
import Box from '@mui/material/Box';

const Navigation = () => {
  return (
    <Box component="nav">
      <h1>This is navigation</h1>
      <ul>
        <li><Link to={'/'}>Home</Link></li>
        <li><Link to={'/about'}>About</Link></li>
      </ul>
    </Box>
  )
}

export default Navigation;
