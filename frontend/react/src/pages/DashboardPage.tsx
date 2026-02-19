import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import api from '../services/api';
import { Organization } from '../types';
import { Plus, Users, ArrowRight, Building2, Loader2, X } from 'lucide-react';

const DashboardPage: React.FC = () => {
  const [organizations, setOrganizations] = useState<Organization[]>([]);
  const [loading, setLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [creating, setCreating] = useState(false);
  
  // Form state
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [defaultWage, setDefaultWage] = useState(50);

  useEffect(() => {
    fetchOrganizations();
  }, []);

  const fetchOrganizations = async () => {
    try {
      const response = await api.get('/organizations');
      setOrganizations(response.data);
    } catch (err: any) {
      console.error('Failed to fetch orgs', err);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateOrg = async (e: React.FormEvent) => {
    e.preventDefault();
    setCreating(true);
    try {
      await api.post('/organizations', {
        name,
        description,
        default_wage: Number(defaultWage)
      });
      setIsModalOpen(false);
      setName('');
      setDescription('');
      setDefaultWage(50);
      fetchOrganizations();
    } catch (err) {
      console.error('Failed to create organization', err);
    } finally {
      setCreating(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center py-20">
        <Loader2 className="animate-spin text-primary" size={48} />
      </div>
    );
  }

  return (
    <div className="space-y-8 animate-fade-in">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-4xl mb-1">Your Organizations</h1>
          <p className="text-text-muted">Manage your teams and track meeting costs</p>
        </div>
        <button 
          onClick={() => setIsModalOpen(true)}
          className="primary flex items-center gap-2"
        >
          <Plus size={20} />
          <span>Create Organization</span>
        </button>
      </div>

      {organizations.length === 0 ? (
        <div className="card text-center py-16 space-y-6 max-w-2xl mx-auto">
          <div className="w-16 h-16 bg-surface rounded-full flex items-center justify-center mx-auto">
            <Building2 className="text-text-muted" size={32} />
          </div>
          <div className="space-y-2">
            <h2 className="text-2xl">No organizations yet</h2>
            <p className="text-text-muted">Create your first organization to start tracking meeting costs with your team.</p>
          </div>
          <button 
            onClick={() => setIsModalOpen(true)}
            className="primary px-8"
          >
            Create Organization
          </button>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {organizations.map((org) => (
            <Link key={org.id} to={`/org/${org.id}`} className="block group">
              <div className="card h-full flex flex-col hover:border-primary/50 transition-all border-glass-border">
                <div className="flex items-start justify-between mb-4">
                  <div className="w-12 h-12 bg-primary/10 rounded-xl flex items-center justify-center group-hover:bg-primary/20 transition-colors">
                    <Building2 className="text-primary" size={24} />
                  </div>
                  <div className="flex items-center gap-1 text-sm font-medium text-success">
                    <Users size={16} />
                    {org.member_count || 0}
                  </div>
                </div>
                
                <h3 className="text-xl font-bold mb-2 group-hover:text-primary transition-colors">{org.name}</h3>
                <p className="text-text-muted text-sm flex-1 line-clamp-2">{org.description || 'No description provided.'}</p>
                
                <div className="mt-6 flex items-center justify-between text-sm">
                  <span className="text-text-muted">Default Wage: ${org.default_wage}/hr</span>
                  <div className="flex items-center gap-1 font-semibold text-primary">
                    View <ArrowRight size={16} />
                  </div>
                </div>
              </div>
            </Link>
          ))}
        </div>
      )}

      {/* Create Organization Modal */}
      {isModalOpen && (
        <div className="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-background/80 backdrop-blur-sm animate-fade-in">
          <div className="card max-w-md w-full shadow-2xl space-y-6 bg-background border-glass-border">
            <div className="flex justify-between items-center">
              <h3 className="text-2xl font-bold">New Organization</h3>
              <button 
                onClick={() => setIsModalOpen(false)}
                className="text-text-muted hover:text-text transition-colors"
              >
                <X size={24} />
              </button>
            </div>
            
            <form onSubmit={handleCreateOrg} className="space-y-6">
              <div className="space-y-2">
                <label className="text-sm font-medium">Organization Name</label>
                <input
                  type="text"
                  required
                  className="w-full"
                  placeholder="e.g. Acme Corp, Engineering Team"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                />
              </div>

              <div className="space-y-2">
                <label className="text-sm font-medium">Description (Optional)</label>
                <textarea
                  rows={3}
                  className="w-full"
                  placeholder="What does this organization do?"
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                />
              </div>

              <div className="space-y-2">
                <label className="text-sm font-medium">Default Hourly Wage ($)</label>
                <input
                  type="number"
                  required
                  min="0"
                  className="w-full"
                  value={defaultWage}
                  onChange={(e) => setDefaultWage(Number(e.target.value))}
                />
                <p className="text-xs text-text-muted">Used as the fallback wage for members without a specific rate.</p>
              </div>

              <div className="flex gap-4 pt-4">
                <button 
                  type="button"
                  onClick={() => setIsModalOpen(false)}
                  className="flex-1 bg-surface hover:bg-surface-hover"
                >
                  Cancel
                </button>
                <button 
                  type="submit"
                  disabled={creating}
                  className="flex-1 primary flex items-center justify-center gap-2"
                >
                  {creating ? <Loader2 className="animate-spin" size={18} /> : 'Create'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default DashboardPage;
