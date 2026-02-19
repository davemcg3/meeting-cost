package service

import (
	"github.com/google/uuid"
)

// EventType defines the type of event being broadcasted.
type EventType string

const (
	EventMeetingStarted     EventType = "meeting:started"
	EventMeetingStopped     EventType = "meeting:stopped"
	EventAttendeeCount      EventType = "meeting:attendee_count"
	EventAverageWage        EventType = "meeting:average_wage"
	EventMeetingCost        EventType = "meeting:cost"
	EventMeetingParticipant EventType = "meeting:participant"
)

// MeetingEvent represents a message broadcasted via websocket.
type MeetingEvent struct {
	Type      EventType   `json:"type"`
	MeetingID uuid.UUID   `json:"meeting_id"`
	Payload   interface{} `json:"payload"`
}
