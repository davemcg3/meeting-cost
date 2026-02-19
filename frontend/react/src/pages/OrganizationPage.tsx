import React, { useEffect, useState } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import api from '../services/api';
import { organizationService } from '../services/organizationService';
import { meetingService } from '../services/meetingService';
import { Organization, Meeting, MemberDTO } from '../types';
import {
  Users,
  Settings,
  Plus,
  ChevronRight,
  Calendar,
  Clock,
  TrendingUp,
  Loader2,
  Building2,
  X,
  Edit2,
  Mail,
  DollarSign,
  Trash2,
  AlertTriangle
} from 'lucide-react';
import { format } from 'date-fns';

const OrganizationPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [org, setOrg] = useState<Organization | null>(null);
  const [meetings, setMeetings] = useState<Meeting[]>([]);
  const [members, setMembers] = useState<MemberDTO[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  // Modals state
  const [isMeetingModalOpen, setIsMeetingModalOpen] = useState(false);
  const [isSettingsModalOpen, setIsSettingsModalOpen] = useState(false);
  const [isInviteModalOpen, setIsInviteModalOpen] = useState(false);
  const [isEditMemberModalOpen, setIsEditMemberModalOpen] = useState(false);
  const [isDeleteConfirmModalOpen, setIsDeleteConfirmModalOpen] = useState(false);
  const [deleteConfirmName, setDeleteConfirmName] = useState('');

  // Form states
  const [newMeetingPurpose, setNewMeetingPurpose] = useState('');
  const [editOrgName, setEditOrgName] = useState('');
  const [editOrgDescription, setEditOrgDescription] = useState('');
  const [editOrgDefaultWage, setEditOrgDefaultWage] = useState(0);
  const [inviteEmail, setInviteEmail] = useState('');
  const [inviteWage, setInviteWage] = useState(0);
  const [editingMember, setEditingMember] = useState<MemberDTO | null>(null);
  const [editMemberWage, setEditMemberWage] = useState(0);

  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    if (id) {
      fetchData();
    }
  }, [id]);

  const fetchData = async () => {
    try {
      const [o, ms, mems] = await Promise.all([
        organizationService.getOrganization(id!),
        organizationService.getMeetings(id!),
        organizationService.getMembers(id!)
      ]);
      setOrg(o);
      setMeetings(ms);
      setMembers(mems);

      // Init form states
      setEditOrgName(o.name);
      setEditOrgDescription(o.description);
      setEditOrgDefaultWage(o.default_wage);
      setInviteWage(o.default_wage);
    } catch (err) {
      setError('Failed to load organization data');
    } finally {
      setLoading(false);
    }
  };

  const handleCreateMeeting = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newMeetingPurpose.trim()) return;
    setSubmitting(true);
    try {
      const newMeeting = await meetingService.createMeeting(id!, newMeetingPurpose);
      navigate(`/meeting/${newMeeting.id}`);
    } catch (err: any) {
      console.error('Failed to create meeting', err);
      alert(err.response?.data?.error || 'Failed to create meeting');
    } finally {
      setSubmitting(false);
    }
  };

  const handleUpdateSettings = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    try {
      await api.patch(`/organizations/${id}`, {
        name: editOrgName,
        description: editOrgDescription,
        default_wage: editOrgDefaultWage
      });
      setIsSettingsModalOpen(false);
      fetchData();
    } catch (err) {
      console.error('Failed to update settings', err);
    } finally {
      setSubmitting(false);
    }
  };

  const handleInviteMember = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    try {
      await organizationService.addMember(id!, inviteEmail, inviteWage);
      setIsInviteModalOpen(false);
      setInviteEmail('');
      fetchData();
    } catch (err: any) {
      console.error('Failed to invite member', err);
      alert(err.response?.data?.error || 'Failed to invite member');
    } finally {
      setSubmitting(false);
    }
  };

  const handleUpdateMemberWage = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingMember) return;
    setSubmitting(true);
    try {
      await organizationService.updateWage(id!, editingMember.person_id, editMemberWage);
      setIsEditMemberModalOpen(false);
      setEditingMember(null);
      fetchData();
    } catch (err) {
      console.error('Failed to update member wage', err);
    } finally {
      setSubmitting(false);
    }
  };

  const openEditMember = (member: MemberDTO) => {
    setEditingMember(member);
    setEditMemberWage(member.hourly_wage || org?.default_wage || 0);
    setIsEditMemberModalOpen(true);
  };

  const handleDeleteOrganization = async () => {
    if (!org || deleteConfirmName !== org.name) return;
    setSubmitting(true);
    try {
      await api.delete(`/organizations/${id}`);
      navigate('/dashboard');
    } catch (err: any) {
      console.error('Failed to delete organization', err);
      alert(err.response?.data?.error || 'Failed to delete organization');
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center py-20">
        <Loader2 className="animate-spin text-primary" size={48} />
      </div>
    );
  }

  if (error || !org) {
    return (
      <div className="card max-w-lg mx-auto mt-10 text-center">
        <h2 className="text-2xl font-bold text-danger mb-4">Error</h2>
        <p className="text-text-muted mb-6">{error || 'Organization not found'}</p>
        <button onClick={() => navigate('/dashboard')} className="primary">Back to Dashboard</button>
      </div>
    );
  }

  return (
    <div className="space-y-8 animate-fade-in">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-end justify-between gap-6">
        <div className="space-y-2">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-primary/10 rounded-xl flex items-center justify-center">
              <Building2 className="text-primary" size={24} />
            </div>
            <h1 className="text-4xl mb-0">{org.name}</h1>
          </div>
          <p className="text-text-muted max-w-2xl">{org.description || 'No description provided.'}</p>
        </div>

        <div className="flex items-center gap-3">
          <button
            onClick={() => setIsSettingsModalOpen(true)}
            className="bg-surface hover:bg-surface-hover border border-glass-border flex items-center gap-2 "
          >
            <Settings size={20} />
            Settings
          </button>
          <button
            onClick={() => setIsMeetingModalOpen(true)}
            className="primary flex items-center gap-2"
          >
            <Plus size={20} />
            New Meeting
          </button>
        </div>
      </div>

      {/* Stats Overview */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="card p-6 flex items-center gap-4">
          <div className="w-12 h-12 bg-success/10 text-success rounded-full flex items-center justify-center">
            <Clock size={24} />
          </div>
          <div>
            <p className="text-text-muted text-sm font-medium uppercase tracking-wider text-[10px]">Total Meetings</p>
            <p className="text-2xl font-bold">{meetings.length}</p>
          </div>
        </div>
        <div className="card p-6 flex items-center gap-4">
          <div className="w-12 h-12 bg-primary/10 text-primary rounded-full flex items-center justify-center">
            <Users size={24} />
          </div>
          <div>
            <p className="text-text-muted text-sm font-medium uppercase tracking-wider text-[10px]">Members</p>
            <p className="text-2xl font-bold">{members.length}</p>
          </div>
        </div>
        <div className="card p-6 flex items-center gap-4">
          <div className="w-12 h-12 bg-warning/10 text-warning rounded-full flex items-center justify-center">
            <TrendingUp size={24} />
          </div>
          <div>
            <p className="text-text-muted text-sm font-medium uppercase tracking-wider text-[10px]">Default Wage</p>
            <p className="text-2xl font-bold">${org.default_wage}/hr</p>
          </div>
        </div>
      </div>

      {/* Main Content Sections */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Meetings List */}
        <div className="lg:col-span-2 space-y-6">
          <div className="flex items-center justify-between">
            <h2 className="text-2xl font-bold">Recent Meetings</h2>
          </div>

          <div className="space-y-4">
            {meetings.length === 0 ? (
              <div className="card py-12 text-center text-text-muted italic border-dashed">
                No meetings recorded yet. Start one to see it here.
              </div>
            ) : (
              meetings.map((meeting) => (
                <Link
                  key={meeting.id}
                  to={`/meeting/${meeting.id}`}
                  className="card p-5 group flex items-center justify-between hover:border-primary/50 transition-all border-glass-border"
                >
                  <div className="flex items-center gap-5">
                    <div className={`w-12 h-12 rounded-xl flex items-center justify-center transition-colors ${meeting.is_active ? 'bg-success/20 text-success shadow-[0_0_15px_rgba(34,197,94,0.3)] animate-pulse' : 'bg-surface text-text-muted group-hover:bg-primary/10'
                      }`}>
                      <Calendar size={24} />
                    </div>
                    <div>
                      <h4 className="font-bold text-lg group-hover:text-primary transition-colors">{meeting.purpose}</h4>
                      <div className="flex items-center gap-4 text-sm text-text-muted mt-1">
                        <div className="flex items-center gap-1">
                          <Clock size={14} />
                          {meeting.is_active ? 'Active' : format(new Date(meeting.created_at), 'MMM d, yyyy')}
                        </div>
                        <div>{meeting.max_attendees} attendees</div>
                      </div>
                    </div>
                  </div>
                  <div className="text-right flex items-center gap-6">
                    <div>
                      <span className="text-[10px] text-text-muted uppercase tracking-wider font-bold">Cost: </span>
                      <span className={`font-bold text-lg ${meeting.is_active ? 'text-success' : ''}`}>
                        ${meeting.total_cost.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                      </span>
                    </div>
                    <ChevronRight className="text-text-muted group-hover:text-primary translate-x-0 group-hover:translate-x-1 transition-all" size={24} />
                  </div>
                </Link>
              ))
            )}
          </div>
        </div>

        {/* Members Sidebar */}
        <div className="space-y-6">
          <div className="flex items-center justify-between">
            <h2 className="text-2xl font-bold">Members</h2>
            <button
              onClick={() => setIsInviteModalOpen(true)}
              className="p-1 px-3 text-xs font-bold bg-primary/10 text-primary rounded-lg hover:bg-primary/20 transition-colors"
            >
              Add Member
            </button>
          </div>

          <div className="card p-0 overflow-hidden divide-y divide-glass-border">
            {members.map((member) => (
              <div key={member.person_id} className="p-4 flex items-center justify-between group">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 bg-surface rounded-full flex items-center justify-center font-bold text-primary border border-glass-border">
                    {member.first_name[0]}{member.last_name ? member.last_name[0] : ''}
                  </div>
                  <div>
                    <p className="font-bold text-sm leading-none mb-1">{member.first_name} {member.last_name}</p>
                    <p className="text-xs text-text-muted truncate max-w-[120px]">{member.email}</p>
                  </div>
                </div>
                <div className="flex items-center gap-3">
                  <div className="text-right">
                    <p className="text-sm font-bold">${member.hourly_wage || org.default_wage}/hr</p>
                    <p className="text-[9px] text-text-muted uppercase font-bold">Rate</p>
                  </div>
                  <button
                    onClick={() => openEditMember(member)}
                    className="p-2 text-text-muted hover:text-primary transition-colors opacity-0 group-hover:opacity-100"
                  >
                    <Edit2 size={14} />
                  </button>
                </div>
              </div>
            ))}
          </div>

          <div
            onClick={() => setIsInviteModalOpen(true)}
            className="card bg-surface/30 border-dashed border-2 flex flex-col items-center justify-center py-8 gap-3 text-center cursor-pointer hover:bg-surface/50 transition-colors"
          >
            <div className="w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center text-primary">
              <Mail size={20} />
            </div>
            <p className="text-xs text-text-muted px-4">Invite more team members to this organization to start tracking their costs.</p>
            <button className="primary text-xs py-2 px-4 flex items-center gap-2">
              <Plus size={14} />
              Invite Member
            </button>
          </div>
        </div>
      </div>

      {/* New Meeting Modal */}
      {isMeetingModalOpen && (
        <div className="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-background/80 backdrop-blur-sm animate-fade-in">
          <div className="card max-w-md w-full shadow-2xl space-y-6">
            <div className="flex justify-between items-center">
              <h3 className="text-2xl font-bold">New Meeting</h3>
              <button
                onClick={() => setIsMeetingModalOpen(false)}
                className="text-text-muted hover:text-text"
              >
                <X size={20} />
              </button>
            </div>

            <form onSubmit={handleCreateMeeting} className="space-y-6">
              <div className="space-y-2">
                <label className="text-sm font-medium">Meeting Purpose</label>
                <textarea
                  required
                  rows={3}
                  className="w-full"
                  placeholder="e.g. Weekly Sync, Design Review, Sprint Planning"
                  value={newMeetingPurpose}
                  onChange={(e) => setNewMeetingPurpose(e.target.value)}
                />
              </div>

              <div className="flex gap-4 pt-4">
                <button
                  type="button"
                  onClick={() => setIsMeetingModalOpen(false)}
                  className="flex-1 bg-surface hover:bg-surface-hover"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={submitting}
                  className="flex-1 primary flex items-center justify-center gap-2"
                >
                  {submitting ? <Loader2 className="animate-spin" size={18} /> : 'Create & Open'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Settings Modal */}
      {isSettingsModalOpen && (
        <div className="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-background/80 backdrop-blur-sm animate-fade-in">
          <div className="card max-w-md w-full shadow-2xl space-y-6">
            <div className="flex justify-between items-center">
              <h3 className="text-2xl font-bold flex items-center gap-2">
                <Settings size={24} />
                Organization Settings
              </h3>
              <button onClick={() => setIsSettingsModalOpen(false)} className="text-text-muted">
                <X size={20} />
              </button>
            </div>

            <form onSubmit={handleUpdateSettings} className="space-y-6">
              <div className="space-y-2">
                <label className="text-sm font-medium">Name</label>
                <input
                  type="text"
                  required
                  className="w-full"
                  value={editOrgName}
                  onChange={(e) => setEditOrgName(e.target.value)}
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">Description</label>
                <textarea
                  rows={3}
                  className="w-full"
                  value={editOrgDescription}
                  onChange={(e) => setEditOrgDescription(e.target.value)}
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">Default Hourly Wage ($)</label>
                <input
                  type="number"
                  step="0.01"
                  className="w-full"
                  value={editOrgDefaultWage}
                  onChange={(e) => setEditOrgDefaultWage(Number(e.target.value))}
                />
              </div>

              <div className="flex gap-4 pt-4">
                <button type="button" onClick={() => setIsSettingsModalOpen(false)} className="flex-1 bg-surface">Cancel</button>
                <button type="submit" disabled={submitting} className="flex-1 primary">
                  {submitting ? <Loader2 className="animate-spin" size={18} /> : 'Save Changes'}
                </button>
              </div>
            </form>

            {/* Danger Zone */}
            <div className="border border-danger/30 rounded-xl p-4 space-y-3 bg-danger/5">
              <div className="flex items-center gap-2 text-danger">
                <AlertTriangle size={16} />
                <p className="text-sm font-bold uppercase tracking-wider">Danger Zone</p>
              </div>
              <p className="text-xs text-text-muted">Permanently delete this organization and all its data. This action cannot be undone.</p>
              <button
                type="button"
                onClick={() => {
                  setDeleteConfirmName('');
                  setIsSettingsModalOpen(false);
                  setIsDeleteConfirmModalOpen(true);
                }}
                className="flex items-center gap-2 text-sm font-bold text-danger border border-danger/40 hover:bg-danger/10 transition-colors px-4 py-2 rounded-lg w-full justify-center"
              >
                <Trash2 size={16} />
                Delete Organization
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Invite Modal */}
      {isInviteModalOpen && (
        <div className="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-background/80 backdrop-blur-sm animate-fade-in">
          <div className="card max-w-md w-full shadow-2xl space-y-6">
            <div className="flex justify-between items-center">
              <h3 className="text-2xl font-bold flex items-center gap-2">
                <Users size={24} />
                Add New Member
              </h3>
              <button onClick={() => setIsInviteModalOpen(false)} className="text-text-muted text-sm">
                <X size={20} />
              </button>
            </div>

            <form onSubmit={handleInviteMember} className="space-y-6">
              <div className="space-y-2">
                <label className="text-sm font-medium">Member Email</label>
                <input
                  type="email"
                  required
                  placeholder="Enter the person's email address"
                  className="w-full"
                  value={inviteEmail}
                  onChange={(e) => setInviteEmail(e.target.value)}
                />
                <p className="text-[10px] text-text-muted">Enter the email address of an existing user to add them to your organization.</p>
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">Specific Hourly Wage (Optional)</label>
                <div className="relative">
                  <DollarSign className="absolute left-3 top-1/2 -translate-y-1/2 text-text-muted" size={16} />
                  <input
                    type="number"
                    step="0.01"
                    className="w-full pl-10"
                    placeholder={org.default_wage.toString()}
                    value={inviteWage}
                    onChange={(e) => setInviteWage(Number(e.target.value))}
                  />
                </div>
              </div>

              <div className="flex gap-4 pt-4">
                <button type="button" onClick={() => setIsInviteModalOpen(false)} className="flex-1 bg-surface">Cancel</button>
                <button type="submit" disabled={submitting} className="flex-1 primary">
                  {submitting ? <Loader2 className="animate-spin" size={18} /> : 'Add to Team'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Edit Member Modal */}
      {isEditMemberModalOpen && (
        <div className="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-background/80 backdrop-blur-sm animate-fade-in">
          <div className="card max-w-md w-full shadow-2xl space-y-6">
            <div className="flex justify-between items-center">
              <h3 className="text-2xl font-bold flex items-center gap-2">
                <Edit2 size={24} />
                Edit Member Rate
              </h3>
              <button onClick={() => setIsEditMemberModalOpen(false)} className="text-text-muted">
                <X size={20} />
              </button>
            </div>

            <form onSubmit={handleUpdateMemberWage} className="space-y-6">
              <div className="flex items-center gap-3 p-4 bg-surface rounded-xl">
                <div className="w-12 h-12 bg-primary/10 rounded-full flex items-center justify-center font-bold text-primary">
                  {editingMember?.first_name[0]}{editingMember?.last_name ? editingMember.last_name[0] : ''}
                </div>
                <div>
                  <p className="font-bold">{editingMember?.first_name} {editingMember?.last_name}</p>
                  <p className="text-xs text-text-muted">{editingMember?.email}</p>
                </div>
              </div>

              <div className="space-y-2">
                <label className="text-sm font-medium">Hourly Wage ($)</label>
                <div className="relative">
                  <DollarSign className="absolute left-3 top-1/2 -translate-y-1/2 text-text-muted" size={16} />
                  <input
                    type="number"
                    step="0.01"
                    required
                    className="w-full pl-10"
                    value={editMemberWage}
                    onChange={(e) => setEditMemberWage(Number(e.target.value))}
                  />
                </div>
              </div>

              <div className="flex gap-4 pt-4">
                <button type="button" onClick={() => setIsEditMemberModalOpen(false)} className="flex-1 bg-surface">Cancel</button>
                <button type="submit" disabled={submitting} className="flex-1 primary">
                  {submitting ? <Loader2 className="animate-spin" size={18} /> : 'Update Rate'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Delete Confirmation Modal */}
      {isDeleteConfirmModalOpen && (
        <div className="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-background/80 backdrop-blur-sm animate-fade-in">
          <div className="card max-w-md w-full shadow-2xl space-y-6 border border-danger/30">
            <div className="flex justify-between items-center">
              <h3 className="text-2xl font-bold flex items-center gap-2 text-danger">
                <AlertTriangle size={24} />
                Delete Organization
              </h3>
              <button onClick={() => setIsDeleteConfirmModalOpen(false)} className="text-text-muted hover:text-text">
                <X size={20} />
              </button>
            </div>

            <div className="bg-danger/10 border border-danger/20 rounded-xl p-4 space-y-1">
              <p className="text-sm font-bold text-danger">This action is permanent and cannot be undone.</p>
              <p className="text-xs text-text-muted">
                Deleting <span className="font-bold text-text">{org.name}</span> will permanently remove all meetings, members, and settings associated with this organization.
              </p>
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium">
                Type <span className="font-bold text-danger">{org.name}</span> to confirm
              </label>
              <input
                type="text"
                className="w-full border-danger/40 focus:border-danger"
                placeholder={org.name}
                value={deleteConfirmName}
                onChange={(e) => setDeleteConfirmName(e.target.value)}
              />
            </div>

            <div className="flex gap-4 pt-2">
              <button
                type="button"
                onClick={() => setIsDeleteConfirmModalOpen(false)}
                className="flex-1 bg-surface hover:bg-surface-hover"
              >
                Cancel
              </button>
              <button
                type="button"
                disabled={deleteConfirmName !== org.name || submitting}
                onClick={handleDeleteOrganization}
                className="flex-1 flex items-center justify-center gap-2 font-bold text-white bg-danger hover:bg-danger/80 disabled:opacity-40 disabled:cursor-not-allowed transition-all px-4 py-2 rounded-xl"
              >
                {submitting ? <Loader2 className="animate-spin" size={18} /> : (
                  <>
                    <Trash2 size={16} />
                    Delete Forever
                  </>
                )}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default OrganizationPage;
