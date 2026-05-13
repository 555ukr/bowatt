import './App.css';
import Home from './pages/home';
import NewPhoto from './pages/new_photo';

import { BrowserRouter, Routes, Route, Link } from 'react-router-dom';

function App() {
  return (
    <BrowserRouter>
      <nav className="navbar">
        <Link to="/">Feed</Link>
        <Link to="/new-photo">+ New Photo</Link>
      </nav>

      <div className="page">
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/new-photo" element={<NewPhoto />} />
        </Routes>
      </div>
    </BrowserRouter>
  );
}

export default App;
