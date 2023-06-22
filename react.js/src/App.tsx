import { Navigate, Route, Routes } from 'react-router-dom';
import Layout from './components/Layout';
import Home from './pages/Home';
import About from './pages/About';
import DexLayout from './pages/Dex/Layout';
import DexClientList from './pages/Dex/ClientList';
import DexClient from './pages/Dex/Client';

function App() {
    return (
        <Routes>
            <Route element={<Layout />}>
                <Route path="/" element={<Home />} />
                <Route path="/dex" element={<Navigate to="clients" replace />} />
                <Route element={<DexLayout />}>
                    <Route path="/dex/clients" element={<DexClientList />} />
                    <Route path="/dex/clients/:id" element={<DexClient />} />
                </Route>
                <Route path="/about" element={<About />} />
            </Route>
        </Routes>
    );
}

export default App;
