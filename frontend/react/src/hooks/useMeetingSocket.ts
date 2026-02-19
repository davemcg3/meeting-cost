import { useEffect, useRef, useState } from 'react';
import { MeetingEvent } from '../types';

export const useMeetingSocket = (meetingId: string | undefined) => {
  const [lastEvent, setLastEvent] = useState<MeetingEvent | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const socketRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    if (!meetingId) return;

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = import.meta.env.VITE_WS_URL || 'localhost:8080';
    const url = `${protocol}//${host}/ws/meetings/${meetingId}`;

    const connect = () => {
      const ws = new WebSocket(url);
      socketRef.current = ws;

      ws.onopen = () => {
        console.log('Meeting socket connected');
        setIsConnected(true);
      };

      ws.onmessage = (event) => {
        try {
          const data: MeetingEvent = JSON.parse(event.data);
          setLastEvent(data);
        } catch (err) {
          console.error('Failed to parse websocket message', err);
        }
      };

      ws.onclose = () => {
        console.log('Meeting socket closed');
        setIsConnected(false);
        // Attempt to reconnect after some time
        setTimeout(connect, 3000);
      };

      ws.onerror = (err) => {
        console.error('Meeting socket error', err);
        ws.close();
      };
    };

    connect();

    return () => {
      if (socketRef.current) {
        socketRef.current.close();
      }
    };
  }, [meetingId]);

  return { lastEvent, isConnected };
};
