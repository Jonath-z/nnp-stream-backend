package models

type UploadDraftPayload struct {
	AssetID    string
	PlaybackId string
	Duration   float64
}

type UploadDraftAssetDTO struct {
	MuxAssetId string  `json:"mux_asset_id,omitempty"`
	PlaybackID string  `json:"playback_id,omitempty"`
	Duration   float64 `json:"duration"`
}

type DraftAsset struct {
	Id         string `json:"id"`
	CreatedAt  string `json:"created_at"`
	MuxAssetId string `json:"mux_asset_id"`
}
