import { Outlet } from 'react-router-dom';

function DexLayout() {
    return (
        <div>
            <h2>Dex Layout</h2>
            <Outlet />
        </div>
    );
}

export default DexLayout;
