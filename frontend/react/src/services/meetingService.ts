import api from './api';
import { Meeting, MeetingCost } from '../types';

export const meetingService = {
  getMeeting: async (id: string): Promise<Meeting> => {
    const { data } = await api.get(`/meetings/${id}`);
    return data;
  },

  startMeeting: async (id: string): Promise<void> => {
    await api.post(`/meetings/${id}/start`);
  },

  stopMeeting: async (id: string): Promise<void> => {
    await api.post(`/meetings/${id}/stop`);
  },

  updateAttendees: async (id: string, count: number): Promise<void> => {
    await api.patch(`/meetings/${id}/attendees`, { count });
  },

  getCost: async (id: string): Promise<MeetingCost> => {
    const { data } = await api.get(`/meetings/${id}/cost`);
    return data;
  },

  createMeeting: async (orgId: string, purpose: string): Promise<Meeting> => {
    const { data } = await api.post(`/meetings`, {
      organization_id: orgId,
      purpose,
    });
    return data;
  },
};
