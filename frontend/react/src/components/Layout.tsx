import React from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuthStore } from '../context/authStore';
import { LogOut, LayoutDashboard, Clock } from 'lucide-react';
import CookieConsent from './CookieConsent';

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const { user, logout, isAuthenticated } = useAuthStore();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <div className="min-h-screen flex flex-col">
      <header className="sticky top-0 z-50 bg-background/80 backdrop-blur-md border-b border-glass-border">
        <div className="max-w-7xl mx-auto px-4 h-16 flex items-center justify-between">
          <Link to="/" className="flex items-center gap-2">
            <div className="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
              <Clock className="text-white w-5 h-5" />
            </div>
            <span className="font-bold text-xl tracking-tight">MeetingCost</span>
          </Link>

          {isAuthenticated ? (
            <div className="flex items-center gap-6">
              <nav className="hidden md:flex items-center gap-4">
                <Link to="/dashboard" className="text-sm font-medium hover:text-primary transition-colors flex items-center gap-2">
                  <LayoutDashboard size={16} />
                  Dashboard
                </Link>
              </nav>

              <div className="h-8 w-[1px] bg-glass-border"></div>

              <div className="flex items-center gap-4">
                <div className="text-right hidden sm:block">
                  <p className="text-sm font-semibold">{user?.firstName} {user?.lastName}</p>
                  <p className="text-xs text-text-muted">{user?.email}</p>
                </div>
                <button
                  onClick={handleLogout}
                  className="p-2 text-text-muted hover:text-danger transition-colors bg-transparent border-none"
                  title="Logout"
                >
                  <LogOut size={20} />
                </button>
              </div>
            </div>
          ) : (
            <div className="flex items-center gap-4">
              <Link to="/login" className="text-sm font-medium px-4 py-2 rounded-md hover:bg-surface transition-colors">Login</Link>
              <Link to="/register" className="text-sm font-medium px-4 py-2 rounded-md bg-primary hover:bg-primary-hover text-white transition-colors">Register</Link>
            </div>
          )}
        </div>
      </header>

      <main className="flex-1 max-w-7xl mx-auto w-full px-4 py-8">
        {children}
      </main>

      <footer className="border-t border-glass-border py-8 text-center text-text-muted text-sm">
        <p>&copy; {new Date().getFullYear()} Meeting Cost Calculator. All rights reserved.</p>
      </footer>

      <CookieConsent />
    </div>
  );
};

export default Layout;
