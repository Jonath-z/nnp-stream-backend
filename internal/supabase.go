package internal

import (
	"log"
	"os"

	"github.com/nnp-stream-backend/models"
	"github.com/supabase-community/supabase-go"
)

func SupabaseClient() (*supabase.Client, error) {
	return supabase.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_API_KEY"), &supabase.ClientOptions{})
}

const (
	TABLE_DRAFTS = "drafts"
)

func UploadDraftAsset(payload models.UploadDraftPayload) string {
	draftAsset := models.UploadDraftAssetDTO{
		MuxAssetId: payload.AssetID,
		PlaybackID: payload.PlaybackId,
		Duration:   payload.Duration,
	}

	client, err := SupabaseClient()
	if err != nil {
		log.Fatalf("Error initializing Supabase client: %v", err)
	}
	_, _, uploadingError := client.From(TABLE_DRAFTS).Insert(&draftAsset, true, "", "", "").Execute()
	if uploadingError != nil {
		log.Fatalf("Error uploading draft asset: %v", uploadingError)
	}
	log.Printf("Created draft asset: %+v", draftAsset)
	return draftAsset.MuxAssetId
}
