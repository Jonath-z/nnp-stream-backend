package models

import "time"

type MuxWebhookPayload struct {
	Type           string             `json:"type"`
	Object         WebhookObject      `json:"object"`
	ID             string             `json:"id"`
	Environment    WebhookEnvironment `json:"environment"`
	Data           WebhookData        `json:"data"`
	CreatedAt      time.Time          `json:"created_at"`
	AccessorSource *string            `json:"accessor_source"`
	Accessor       *string            `json:"accessor"`
	RequestID      *string            `json:"request_id"`
}

type WebhookObject struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type WebhookEnvironment struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type WebhookData struct {
	Tracks              []WebhookTrack `json:"tracks"`
	Status              string         `json:"status"`
	MaxStoredResolution string         `json:"max_stored_resolution"`
	MaxStoredFrameRate  float64        `json:"max_stored_frame_rate"`
	ID                  string         `json:"id"`
	Duration            float64        `json:"duration"`
	CreatedAt           time.Time      `json:"created_at"`
	AspectRatio         string         `json:"aspect_ratio"`
}

type WebhookTrack struct {
	Type             string  `json:"type"`
	MaxWidth         *int    `json:"max_width,omitempty"`
	MaxHeight        *int    `json:"max_height,omitempty"`
	MaxFrameRate     float64 `json:"max_frame_rate,omitempty"`
	ID               string  `json:"id"`
	Duration         float64 `json:"duration"`
	MaxChannels      *int    `json:"max_channels,omitempty"`
	MaxChannelLayout string  `json:"max_channel_layout,omitempty"`
}
