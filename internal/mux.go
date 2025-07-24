package internal

import (
	"os"

	muxgo "github.com/muxinc/mux-go"
)

func MuxClient() *muxgo.APIClient {
	return muxgo.NewAPIClient(
		muxgo.NewConfiguration(
			muxgo.WithBasicAuth(os.Getenv("MUX_TOKEN_ID"), os.Getenv("MUX_TOKEN_SECRET")),
		))
}

/**
	Generate a new upload URL for a new asset.
*/
func GenerateMuxUploadURL() (muxgo.UploadResponse, error) {
	car := muxgo.CreateAssetRequest{PlaybackPolicy: []muxgo.PlaybackPolicy{muxgo.PUBLIC}}
	cur := muxgo.CreateUploadRequest{NewAssetSettings: car, Timeout: 3600, CorsOrigin: "*"}
	return MuxClient().DirectUploadsApi.CreateDirectUpload(cur)
}


func CreateLiveStream() (muxgo.LiveStreamResponse, error) {
	liveStream := muxgo.CreateLiveStreamRequest{
		PlaybackPolicy: []muxgo.PlaybackPolicy{muxgo.PUBLIC},
		LowLatency: true,
		NewAssetSettings: muxgo.CreateAssetRequest{
			PlaybackPolicy: []muxgo.PlaybackPolicy{muxgo.PUBLIC},
		},
	}
	return MuxClient().LiveStreamsApi.CreateLiveStream(liveStream)
}

func DeleteLiveStream(liveStreamID string) error {
	return MuxClient().LiveStreamsApi.DeleteLiveStream(liveStreamID)
}