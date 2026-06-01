package internal

import (
	"encoding/json"
	"log"
	"os"

	"github.com/nnp-stream-backend/models"
	"github.com/nnp-stream-backend/models/event"
	"github.com/supabase-community/supabase-go"
)

func SupabaseClient() (*supabase.Client, error) {
	return supabase.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_API_KEY"), &supabase.ClientOptions{})
}

const (
	TABLE_DRAFTS = "drafts"
	TABLE_USERS  = "users"
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

func CreateUser(payload event.ClerkUserEvent) string {
	user := models.CreateUserDto{
		Email:       payload.Data.EmailAddresses[0].EmailAddress,
		ClerkUserId: payload.Data.ID,
	}

	client, err := SupabaseClient()
	if err != nil {
		log.Fatalf("Error initializing Supabase client: %v", err)
	}

	_, _, uploadError := client.From(TABLE_USERS).Insert(&user, true, "", "", "").Execute()
	if uploadError != nil {
		log.Fatalf("Error uploading draft asset: %v", uploadError)
	}
	log.Printf("Created user: %+v", user)

	return user.Email
}

func GetUserByClerkUserId(clerkUserId string) (*models.User, error) {
	client, err := SupabaseClient()
	if err != nil {
		log.Fatalf("Error initializing Supabase client: %v", err)
		return nil, err
	}

	userByte, _, err_ := client.From(TABLE_USERS).Select("*", "", false).Eq("clerk_id", clerkUserId).Execute()
	if err_ != nil {
		return nil, err_
	}

	var users []models.User
	if err := json.Unmarshal(userByte, &users); err != nil {
		return nil, err
	}

	if len(users) == 0 {
		log.Fatalf("Error getting user by clerk user: %v", err_)
		return nil, err_
	}

	return &users[0], nil
}
