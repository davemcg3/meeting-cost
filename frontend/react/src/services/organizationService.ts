import api from './api';
import { Organization, MemberDTO } from '../types';

export const organizationService = {
  getOrganization: async (id: string): Promise<Organization> => {
    const { data } = await api.get(`/organizations/${id}`);
    return data;
  },

  listOrganizations: async (): Promise<Organization[]> => {
    const { data } = await api.get('/organizations');
    return data;
  },

  getMembers: async (id: string): Promise<MemberDTO[]> => {
    const { data } = await api.get(`/organizations/${id}/members`);
    return data;
  },

  addMember: async (orgId: string, memberIdentifier: string, wage?: number): Promise<void> => {
    const payload: any = { wage };
    if (memberIdentifier.includes('@')) {
      payload.email = memberIdentifier;
    } else {
      payload.person_id = memberIdentifier;
    }
    await api.post(`/organizations/${orgId}/members`, payload);
  },

  updateWage: async (orgId: string, memberId: string, wage: number): Promise<void> => {
    await api.patch(`/organizations/${orgId}/members/${memberId}/wage`, { wage });
  },

  removeMember: async (orgId: string, memberId: string): Promise<void> => {
    await api.delete(`/organizations/${orgId}/members/${memberId}`);
  },

  getMeetings: async (orgId: string): Promise<any[]> => {
    const { data } = await api.get(`/meetings?organization_id=${orgId}`);
    return data;
  },
};
