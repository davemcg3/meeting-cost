export interface User {
  id: string;
  email: string;
  firstName: string;
  lastName?: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  user: User;
}

export interface Organization {
  id: string;
  name: string;
  slug: string;
  description: string;
  default_wage: number;
  use_blended_wage: boolean;
  created_at: string;
  member_count: number;
}

export interface Meeting {
  id: string;
  organization_id: string;
  purpose: string;
  started_at?: string;
  stopped_at?: string;
  is_active: boolean;
  total_cost: number;
  total_duration: number;
  max_attendees: number;
  created_at: string;
}

export interface MeetingCost {
  total_cost: number;
  total_duration: number;
  cost_per_second: number;
  cost_per_minute: number;
  cost_per_hour: number;
}

export interface MeetingEvent {
  type: string;
  meeting_id: string;
  payload: any;
}

export interface MemberDTO {
  person_id: string;
  organization_id: string;
  email: string;
  first_name: string;
  last_name: string;
  is_active: boolean;
  hourly_wage?: number;
  joined_at: string;
}
