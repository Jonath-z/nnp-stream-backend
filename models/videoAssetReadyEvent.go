package models

type VideoAssetReadyEvent struct {
	Type           string      `json:"type"`
	RequestID      interface{} `json:"request_id"`
	Object         Object      `json:"object"`
	ID             string      `json:"id"`
	Environment    Environment `json:"environment"`
	Data           Data        `json:"data"`
	CreatedAt      string      `json:"created_at"`
	Attempts       []any       `json:"attempts"`
	AccessorSource interface{} `json:"accessor_source"`
	Accessor       interface{} `json:"accessor"`
}

type Object struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type Environment struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type Data struct {
	VideoQuality            string                  `json:"video_quality"`
	UploadID                string                  `json:"upload_id"`
	Tracks                  []Track                 `json:"tracks"`
	Status                  string                  `json:"status"`
	ResolutionTier          string                  `json:"resolution_tier"`
	Progress                Progress                `json:"progress"`
	PlaybackIDs             []PlaybackID            `json:"playback_ids"`
	NonStandardInputReasons NonStandardInputReasons `json:"non_standard_input_reasons"`
	MP4Support              string                  `json:"mp4_support"`
	MaxStoredResolution     string                  `json:"max_stored_resolution"`
	MaxStoredFrameRate      float64                 `json:"max_stored_frame_rate"`
	MaxResolutionTier       string                  `json:"max_resolution_tier"`
	MasterAccess            string                  `json:"master_access"`
	IngestType              string                  `json:"ingest_type"`
	ID                      string                  `json:"id"`
	EncodingTier            string                  `json:"encoding_tier"`
	Duration                float64                 `json:"duration"`
	CreatedAt               int64                   `json:"created_at"`
	AspectRatio             string                  `json:"aspect_ratio"`
}

type Track struct {
	Type         string  `json:"type"`
	MaxWidth     int     `json:"max_width"`
	MaxHeight    int     `json:"max_height"`
	MaxFrameRate float64 `json:"max_frame_rate"`
	ID           string  `json:"id"`
}

type Progress struct {
	State string `json:"state"`
}

type PlaybackID struct {
	Policy string `json:"policy"`
	ID     string `json:"id"`
}

type NonStandardInputReasons struct {
	VideoCodec                string `json:"video_codec"`
	UnexpectedVideoParameters string `json:"unexpected_video_parameters"`
}
