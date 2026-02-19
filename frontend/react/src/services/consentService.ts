import api from './api';

export interface CookieConsent {
    id: string;
    person_id?: string;
    session_id: string;
    necessary_cookies: boolean;
    analytics_cookies: boolean;
    marketing_cookies: boolean;
    functional_cookies: boolean;
    consent_version: string;
    consent_date: string;
}

export interface UpdateConsentRequest {
    session_id: string;
    analytics_cookies: boolean;
    marketing_cookies: boolean;
    functional_cookies: boolean;
}

const getSessionId = () => {
    let sessionId = localStorage.getItem('session_id');
    if (!sessionId) {
        sessionId = crypto.randomUUID();
        localStorage.setItem('session_id', sessionId);
    }
    return sessionId;
};

export const consentService = {
    getConsent: async (): Promise<CookieConsent | null> => {
        try {
            const sessionId = getSessionId();
            const { data } = await api.get(`/consent?session_id=${sessionId}`);
            return data;
        } catch (error) {
            return null;
        }
    },

    updateConsent: async (preferences: Omit<UpdateConsentRequest, 'session_id'>): Promise<CookieConsent> => {
        const sessionId = getSessionId();
        const { data } = await api.post('/consent', {
            session_id: sessionId,
            ...preferences,
        });
        return data;
    },

    getHistory: async (): Promise<CookieConsent[]> => {
        const sessionId = getSessionId();
        const { data } = await api.get(`/consent/history?session_id=${sessionId}`);
        return data;
    },

    syncConsent: async (): Promise<void> => {
        try {
            const sessionId = getSessionId();
            await api.post(`/consent/sync?session_id=${sessionId}`);
        } catch (error) {
            console.error('Failed to sync cookie consent', error);
        }
    },
};
