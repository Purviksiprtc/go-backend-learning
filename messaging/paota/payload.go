package paota

import "time"

type UserEventPayload struct {
	Event     string    `json:"event"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	Data      UserData  `json:"data"`
}

type UserData struct {
	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}
