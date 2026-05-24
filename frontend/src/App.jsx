import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Login from './pages/Login';
import Register from './pages/Register';
import Dashboard from './pages/Dashboard';
import Profile from './pages/Profile';
import AdminPanel from './pages/AdminPanel';
import useAuthStore from './store/authStore';

function App() {
  const { isAuthed } = useAuthStore();

  // VULNERABILITY: No CSP (Content Security Policy) headers
  // VULNERABILITY: No protection against clickjacking

  return (
    <Router>
      <div className="min-h-screen bg-gray-50">
        <Routes>
          <Route
            path="/login"
            element={isAuthed ? <Navigate to="/dashboard" /> : <Login />}
          />
          <Route
            path="/register"
            element={isAuthed ? <Navigate to="/dashboard" /> : <Register />}
          />
          <Route
            path="/dashboard"
            element={isAuthed ? <Dashboard /> : <Navigate to="/login" />}
          />
          <Route
            path="/profile"
            element={isAuthed ? <Profile /> : <Navigate to="/login" />}
          />
          <Route
            path="/admin"
            element={isAuthed ? <AdminPanel /> : <Navigate to="/login" />}
          />
          <Route path="/" element={<Navigate to="/dashboard" />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
