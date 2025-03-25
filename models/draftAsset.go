package models

type UploadDraftAssetDTO struct {
	MuxAssetId string `json:"mux_asset_id,omitempty"`
}

type DraftAsset struct {
	Id         string `json:"id"`
	CreatedAt  string `json:"created_at"`
	MuxAssetId string `json:"mux_asset_id"`
}
