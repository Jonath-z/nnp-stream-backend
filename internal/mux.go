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

func ListAssets(page int32, limit int32) (muxgo.ListAssetsResponse, error) {
	params := muxgo.ListAssetsParams{
		Limit: limit,
		Page:  page,
	}
	return MuxClient().AssetsApi.ListAssets(muxgo.WithParams(&params))
}

func GetAsset(assetID string) (muxgo.AssetResponse, error) {
	return MuxClient().AssetsApi.GetAsset(assetID)
}

func ListVideoViews(timeframe []string, limit int32, page int32) (muxgo.ListVideoViewsResponse, error) {
	params := muxgo.ListVideoViewsParams{
		Limit:          limit,
		Page:           page,
		Timeframe:      timeframe,
		OrderDirection: "desc",
	}
	return MuxClient().VideoViewsApi.ListVideoViews(muxgo.WithParams(&params))
}

func GetVideoView(videoViewID string) (muxgo.VideoViewResponse, error) {
	return MuxClient().VideoViewsApi.GetVideoView(videoViewID)
}

func GetOverallMetrics(metricID string, timeframe []string) (muxgo.GetOverallValuesResponse, error) {
	params := muxgo.GetOverallValuesParams{
		Timeframe: timeframe,
	}
	return MuxClient().MetricsApi.GetOverallValues(metricID, muxgo.WithParams(&params))
}

func GetMetricTimeseries(metricID string, timeframe []string) (muxgo.GetMetricTimeseriesDataResponse, error) {
	params := muxgo.GetMetricTimeseriesDataParams{
		Timeframe: timeframe,
	}
	return MuxClient().MetricsApi.GetMetricTimeseriesData(metricID, muxgo.WithParams(&params))
}

func ListMetricBreakdown(metricID string, groupBy string, timeframe []string) (muxgo.ListBreakdownValuesResponse, error) {
	params := muxgo.ListBreakdownValuesParams{
		GroupBy:   groupBy,
		Timeframe: timeframe,
	}
	return MuxClient().MetricsApi.ListBreakdownValues(metricID, muxgo.WithParams(&params))
}