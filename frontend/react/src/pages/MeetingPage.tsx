import React, { useEffect, useState, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { meetingService } from '../services/meetingService';
import { useMeetingSocket } from '../hooks/useMeetingSocket';
import { Meeting, MeetingCost } from '../types';
import { 
  Play, 
  Square, 
  Users, 
  DollarSign, 
  Clock, 
  ChevronLeft,
  Loader2,
  TrendingUp,
  AlertCircle
} from 'lucide-react';
import { motion } from 'framer-motion';

const MeetingPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [meeting, setMeeting] = useState<Meeting | null>(null);
  const [cost, setCost] = useState<MeetingCost | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [attendeeCount, setAttendeeCount] = useState(0);
  
  const { lastEvent } = useMeetingSocket(id);
  const costIntervalRef = useRef<number | null>(null);

  useEffect(() => {
    if (id) {
      fetchMeetingData();
    }
  }, [id]);

  useEffect(() => {
    if (lastEvent) {
      if (lastEvent.type === 'meeting:started' || lastEvent.type === 'meeting:stopped') {
        fetchMeetingData();
      } else if (lastEvent.type === 'meeting:cost') {
        setAttendeeCount(lastEvent.payload.attendee_count);
        fetchCost();
      }
    }
  }, [lastEvent]);

  useEffect(() => {
    if (meeting?.is_active) {
      startCostTicker();
    } else {
      stopCostTicker();
    }
    return () => stopCostTicker();
  }, [meeting?.is_active, cost?.cost_per_second]);

  const fetchMeetingData = async () => {
    try {
      const [m, c] = await Promise.all([
        meetingService.getMeeting(id!),
        meetingService.getCost(id!)
      ]);
      setMeeting(m);
      setCost(c);
      setAttendeeCount(m.max_attendees);
    } catch (err) {
      setError('Failed to load meeting data');
    } finally {
      setLoading(false);
    }
  };

  const fetchCost = async () => {
    try {
      const c = await meetingService.getCost(id!);
      setCost(c);
    } catch (err) {
      console.error('Failed to fetch cost', err);
    }
  };

  const startCostTicker = () => {
    stopCostTicker();
    if (!cost?.cost_per_second) return;

    costIntervalRef.current = window.setInterval(() => {
      setCost(prev => {
        if (!prev) return null;
        return {
          ...prev,
          total_cost: prev.total_cost + prev.cost_per_second,
          total_duration: prev.total_duration + 1
        };
      });
    }, 1000);
  };

  const stopCostTicker = () => {
    if (costIntervalRef.current) {
      clearInterval(costIntervalRef.current);
      costIntervalRef.current = null;
    }
  };

  const handleStart = async () => {
    try {
      await meetingService.startMeeting(id!);
      fetchMeetingData();
    } catch (err) {
      setError('Failed to start meeting');
    }
  };

  const handleStop = async () => {
    try {
      await meetingService.stopMeeting(id!);
      fetchMeetingData();
    } catch (err) {
      setError('Failed to stop meeting');
    }
  };

  const handleUpdateAttendees = async (delta: number) => {
    const newCount = Math.max(0, attendeeCount + delta);
    try {
      await meetingService.updateAttendees(id!, newCount);
      setAttendeeCount(newCount);
    } catch (err) {
      console.error('Failed to update attendees');
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center py-20">
        <Loader2 className="animate-spin text-primary" size={48} />
      </div>
    );
  }

  if (error || !meeting) {
    return (
      <div className="card max-w-lg mx-auto mt-10 text-center space-y-4">
        <AlertCircle className="text-danger mx-auto" size={48} />
        <h2 className="text-2xl font-bold">Error</h2>
        <p className="text-text-muted">{error || 'Meeting not found'}</p>
        <button onClick={() => navigate('/dashboard')} className="primary">Back to Dashboard</button>
      </div>
    );
  }

  return (
    <div className="space-y-8 animate-fade-in">
      <div className="flex items-center justify-between">
        <button 
          onClick={() => navigate('/dashboard')}
          className="flex items-center gap-2 p-0 bg-transparent hover:text-primary transition-colors text-text-muted"
        >
          <ChevronLeft size={20} />
          Back to Dashboard
        </button>
        
        <div className="flex items-center gap-3">
          {meeting.is_active ? (
            <div className="flex items-center gap-2 px-3 py-1 bg-success/10 text-success rounded-full text-sm font-bold animate-pulse-glow">
              <div className="w-2 h-2 bg-success rounded-full"></div>
              LIVE
            </div>
          ) : (
            <div className="px-3 py-1 bg-surface text-text-muted rounded-full text-sm font-bold">
              IDLE
            </div>
          )}
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Main Cost View */}
        <div className="lg:col-span-2 space-y-8">
          <motion.div 
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="card p-10 flex flex-col items-center justify-center text-center space-y-6 overflow-hidden relative"
          >
            <div className="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-primary to-success opacity-50"></div>
            
            <p className="text-text-muted font-medium tracking-widest uppercase text-sm">Targeting Purpose</p>
            <h1 className="text-3xl md:text-5xl font-bold max-w-xl">{meeting.purpose}</h1>

            <div className="py-8">
              <p className="text-text-muted mb-2 text-lg font-medium">Accumulated Cost</p>
              <div className="flex items-baseline justify-center gap-2">
                <span className="text-4xl md:text-6xl font-light text-primary">$</span>
                <span className="text-6xl md:text-9xl font-bold font-mono-numbers tracking-tighter">
                  {cost?.total_cost.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                </span>
              </div>
            </div>

            <div className="grid grid-cols-3 gap-12 w-full pt-8 border-t border-glass-border">
              <div className="space-y-1">
                <p className="text-text-muted text-xs uppercase tracking-wider font-bold">Duration</p>
                <div className="flex items-center justify-center gap-2 text-xl font-bold font-mono-numbers">
                  <Clock className="text-primary" size={20} />
                  {formatDuration(cost?.total_duration || 0)}
                </div>
              </div>
              <div className="space-y-1">
                <p className="text-text-muted text-xs uppercase tracking-wider font-bold">Avg. Wage</p>
                <div className="flex items-center justify-center gap-2 text-xl font-bold font-mono-numbers">
                  <DollarSign className="text-primary" size={20} />
                  {cost?.cost_per_hour.toLocaleString(undefined, { maximumFractionDigits: 0 })}/hr
                </div>
              </div>
              <div className="space-y-1">
                <p className="text-text-muted text-xs uppercase tracking-wider font-bold">Burn Rate</p>
                <div className="flex items-center justify-center gap-2 text-xl font-bold font-mono-numbers text-danger">
                  <TrendingUp size={20} />
                  ${cost?.cost_per_minute.toLocaleString(undefined, { minimumFractionDigits: 2 })}/min
                </div>
              </div>
            </div>
          </motion.div>

          <div className="flex items-center gap-4">
            {!meeting.is_active ? (
              <button 
                onClick={handleStart}
                className="primary flex-1 py-4 text-xl flex items-center justify-center gap-3 bg-success hover:bg-success/90"
              >
                <Play fill="currentColor" size={24} />
                Start Meeting
              </button>
            ) : (
              <button 
                onClick={handleStop}
                className="primary flex-1 py-4 text-xl flex items-center justify-center gap-3 bg-danger hover:bg-danger/90"
              >
                <Square fill="currentColor" size={24} />
                End Meeting
              </button>
            )}
          </div>
        </div>

        {/* Sidebar Controls */}
        <div className="space-y-6">
          <div className="card space-y-6">
            <div className="flex items-center justify-between">
              <h3 className="text-xl font-bold flex items-center gap-2">
                <Users className="text-primary" size={20} />
                Attendees
              </h3>
              <span className="bg-primary/10 text-primary px-3 py-1 rounded-lg font-bold text-lg">
                {attendeeCount}
              </span>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <button 
                onClick={() => handleUpdateAttendees(-1)}
                className="bg-surface hover:bg-surface-hover border border-border p-4 text-2xl font-bold"
              >
                -1
              </button>
              <button 
                onClick={() => handleUpdateAttendees(1)}
                className="bg-surface hover:bg-surface-hover border border-border p-4 text-2xl font-bold"
              >
                +1
              </button>
              <button 
                onClick={() => handleUpdateAttendees(-5)}
                className="bg-surface hover:bg-surface-hover border border-border p-4 text-sm font-bold"
              >
                -5
              </button>
              <button 
                onClick={() => handleUpdateAttendees(5)}
                className="bg-surface hover:bg-surface-hover border border-border p-4 text-sm font-bold"
              >
                +5
              </button>
            </div>

            <p className="text-xs text-text-muted text-center italic">
              Changes to attendee count are tracked in real-time and reflected in the cost calculation.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

const formatDuration = (seconds: number): string => {
  const h = Math.floor(seconds / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  const s = seconds % 60;
  return [h, m, s].map(v => v.toString().padStart(2, '0')).join(':');
};

export default MeetingPage;
