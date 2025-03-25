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

func GenerateMuxUploadURL() (muxgo.UploadResponse, error) {

	car := muxgo.CreateAssetRequest{PlaybackPolicy: []muxgo.PlaybackPolicy{muxgo.PUBLIC}}
	cur := muxgo.CreateUploadRequest{NewAssetSettings: car, Timeout: 3600, CorsOrigin: "*"}
	return MuxClient().DirectUploadsApi.CreateDirectUpload(cur)
}
