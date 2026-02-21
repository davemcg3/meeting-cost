import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuthStore } from '../context/authStore';
import api from '../services/api';
import { Mail, Lock, Loader2 } from 'lucide-react';

const LoginPage: React.FC = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const setAuth = useAuthStore((state) => state.setAuth);
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const response = await api.post('/auth/login', { email, password });

      const { user, access_token, refresh_token } = response.data;
      if (!user || !access_token || !refresh_token) {
        throw new Error('Invalid response format from server');
      }

      setAuth(user, access_token, refresh_token);

      // Sync cookie consent session with the newly logged in user
      import('../services/consentService').then(({ consentService }) => {
        consentService.syncConsent();
      });

      navigate('/dashboard');
    } catch (err: any) {
      console.error("Login caught error:", err);
      const errorMessage = err.response?.data?.error || err.message || 'Invalid email or password';
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-md mx-auto mt-16">
      <div className="card shadow-2xl animate-fade-in">
        <h1 className="text-3xl text-center mb-2">Welcome Back</h1>
        <p className="text-text-muted text-center mb-8">Sign in to manage your meetings</p>

        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="space-y-2">
            <label className="text-sm font-medium">Email Address</label>
            <div className="relative">
              <Mail className="absolute left-3 top-1/2 -translate-y-1/2 text-text-muted" size={18} />
              <input
                type="email"
                required
                className="w-full pl-10"
                placeholder="name@company.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
            </div>
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">Password</label>
            <div className="relative">
              <Lock className="absolute left-3 top-1/2 -translate-y-1/2 text-text-muted" size={18} />
              <input
                type="password"
                required
                className="w-full pl-10"
                placeholder="••••••••"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </div>
          </div>

          {error && (
            <div className="p-3 bg-danger/10 border border-danger/20 text-danger text-sm rounded-md">
              {error}
            </div>
          )}

          <button
            type="submit"
            disabled={loading}
            className="w-full primary flex items-center justify-center gap-2 py-3"
          >
            {loading ? <Loader2 className="animate-spin" size={20} /> : 'Sign In'}
          </button>
        </form>

        <p className="mt-8 text-center text-sm text-text-muted">
          Don't have an account?{' '}
          <Link to="/register" className="text-primary hover:underline font-semibold">
            Create one for free
          </Link>
        </p>
      </div>
    </div>
  );
};

export default LoginPage;
