package handler

import (
	"context"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"github.com/yourorg/meeting-cost/backend/go/internal/cache"
	"github.com/yourorg/meeting-cost/backend/go/internal/logger"
	"github.com/yourorg/meeting-cost/backend/go/internal/pubsub"
	"github.com/yourorg/meeting-cost/backend/go/internal/service"
)

type WebsocketHandler struct {
	pubsub pubsub.PubSub
	logger logger.Logger
}

func NewWebsocketHandler(ps pubsub.PubSub, l logger.Logger) *WebsocketHandler {
	return &WebsocketHandler{
		pubsub: ps,
		logger: l,
	}
}

// HandleMeetingEvents upgrades the connection and streams meeting events.
func (h *WebsocketHandler) HandleMeetingEvents(c *websocket.Conn) {
	meetingID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		h.logger.Error("invalid meeting id for websocket", "error", err)
		c.WriteJSON(fiber.Map{"error": "invalid meeting id"})
		c.Close()
		return
	}

	// In a real app, we should verify the user has access to this meeting.
	// We can pass the person_id via a token in the query param or Sec-WebSocket-Protocol.
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	channel := cache.ChannelMeetingEvents(meetingID)
	events := h.pubsub.Subscribe(ctx, channel)

	h.logger.Info("websocket client connected", "meeting_id", meetingID)

	// Keep alive / ping loop could be here
	
	for {
		select {
		case msg, ok := <-events:
			if !ok {
				return
			}
			
			// We receive a JSON string from Redis, need to send it to client
			var event service.MeetingEvent
			if err := json.Unmarshal([]byte(msg), &event); err != nil {
				h.logger.Error("failed to unmarshal event from pubsub", "error", err)
				continue
			}

			if err := c.WriteJSON(event); err != nil {
				h.logger.Info("websocket client disconnected", "meeting_id", meetingID)
				return
			}
		}
	}
}
