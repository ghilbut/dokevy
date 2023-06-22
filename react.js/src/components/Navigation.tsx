import { Link } from 'react-router-dom';

const Navigation = () => {
  return (
    <nav>
      <h1>This is navigation</h1>
      <ul>
        <li><Link to={'/'}>Home</Link></li>
        <li><Link to={'/about'}>About</Link></li>
      </ul>
    </nav>
  )
}

export default Navigation;
