import { useState, useCallback } from 'react';
import './App.css';
import Home from './pages/home';
import NewPhoto from './pages/new_photo';
import { useWebSocket } from './hooks/useWebSocket';

import { BrowserRouter, Routes, Route, Link } from 'react-router-dom';

interface Photo {
  Path: string;
  Tags: string[];
  CreatedAt: string;
  Data?: string;
}

function App() {
  const [notificationCount, setNotificationCount] = useState(0);
  const [showNotification, setShowNotification] = useState(false);
  const [newPhotos, setNewPhotos] = useState<Photo[]>([]);

  const handleWsMessage = useCallback((data: unknown) => {
    const photo = data as Photo;
    setNewPhotos((prev) => [photo, ...prev]);
    setNotificationCount((prev) => prev + 1);
    setShowNotification(true);
  }, []);

  const { isConnected } = useWebSocket({
    url: 'ws://127.0.0.1:8000/ws',
    onMessage: handleWsMessage,
  });

  const clearNotifications = () => {
    setNotificationCount(0);
  };

  return (
    <BrowserRouter>
      <nav className="navbar">
        <Link to="/" onClick={clearNotifications}>Feed</Link>
        <Link to="/new-photo">+ New Photo</Link>
        <div className="notification-bell-wrapper">
          <span
            className={`notification-bell ${showNotification ? 'ring' : ''}`}
            title={isConnected ? 'Live updates active' : 'Disconnected'}
          >
            🔔
          </span>
          {notificationCount > 0 && (
            <span className="notification-badge">{notificationCount}</span>
          )}
        </div>
      </nav>

      <div className="page">
        <Routes>
          <Route path="/" element={<Home newPhotos={newPhotos} />} />
          <Route path="/new-photo" element={<NewPhoto />} />
        </Routes>
      </div>
    </BrowserRouter>
  );
}

export default App;
