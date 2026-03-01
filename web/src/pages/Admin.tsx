import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { login } from '../api/client';
import StatusBar from '../components/StatusBar';
import './Admin.css';

export default function Admin() {
  const navigate = useNavigate();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const { token } = await login(username, password);
      localStorage.setItem('token', token);
      navigate('/write');
    } catch (err) {
      setError((err as Error).message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="admin-layout">
      <div className="admin-container">
        <div className="admin-title">-- INSERT --</div>
        <form onSubmit={handleSubmit} className="admin-form">
          <label className="admin-label">
            <span>username</span>
            <input
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              autoFocus
              required
            />
          </label>
          <label className="admin-label">
            <span>password</span>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </label>
          {error && <div className="admin-error">{error}</div>}
          <button type="submit" disabled={loading} className="admin-submit">
            {loading ? 'authenticating...' : '[ login ]'}
          </button>
        </form>
      </div>

      <StatusBar left=".env" right="INSERT" />
    </div>
  );
}
