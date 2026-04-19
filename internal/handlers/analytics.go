package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nnp-stream-backend/internal"
)

// HandleListAssets godoc
// @Summary List video assets
// @Description Returns a paginated list of all Mux video assets
// @Tags analytics
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /analytics/assets [get]
func HandleListAssets(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	assets, err := internal.ListAssets(int32(page), int32(limit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assets)
}

// HandleGetAsset godoc
// @Summary Get asset details
// @Description Returns details for a specific Mux video asset including duration, status, resolution, and playback IDs
// @Tags analytics
// @Produce json
// @Param assetID path string true "Mux Asset ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /analytics/assets/{assetID} [get]
func HandleGetAsset(c *gin.Context) {
	assetID := c.Param("assetID")

	asset, err := internal.GetAsset(assetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, asset)
}

// HandleListVideoViews godoc
// @Summary List video views
// @Description Returns a paginated list of video views within a timeframe
// @Tags analytics
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(25)
// @Param timeframe query string false "Timeframe filter (e.g. 24:hours, 7:days)" default(24:hours)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /analytics/views [get]
func HandleListVideoViews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "25"))
	timeframeParam := c.DefaultQuery("timeframe", "24:hours")

	timeframe := parseTimeframe(timeframeParam)

	views, err := internal.ListVideoViews(timeframe, int32(limit), int32(page))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, views)
}

// HandleGetVideoView godoc
// @Summary Get video view details
// @Description Returns detailed information for a specific video view session
// @Tags analytics
// @Produce json
// @Param videoViewID path string true "Video View ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /analytics/views/{videoViewID} [get]
func HandleGetVideoView(c *gin.Context) {
	videoViewID := c.Param("videoViewID")

	view, err := internal.GetVideoView(videoViewID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, view)
}

// HandleGetOverallMetrics godoc
// @Summary Get overall metric values
// @Description Returns aggregate values for a specific metric (e.g. views, video_startup_time, rebuffer_percentage)
// @Tags analytics
// @Produce json
// @Param metric query string false "Metric ID (views, video_startup_time, rebuffer_percentage, rebuffer_frequency, exits_before_video_start, playback_failure_percentage, player_startup_time)" default(video_startup_time)
// @Param timeframe query string false "Timeframe filter" default(24:hours)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /analytics/metrics/overall [get]
func HandleGetOverallMetrics(c *gin.Context) {
	metricID := c.DefaultQuery("metric", "video_startup_time")
	timeframeParam := c.DefaultQuery("timeframe", "24:hours")

	timeframe := parseTimeframe(timeframeParam)

	metrics, err := internal.GetOverallMetrics(metricID, timeframe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// HandleGetMetricTimeseries godoc
// @Summary Get metric timeseries data
// @Description Returns metric values over time, useful for charting trends on the dashboard
// @Tags analytics
// @Produce json
// @Param metric query string false "Metric ID" default(views)
// @Param timeframe query string false "Timeframe filter" default(24:hours)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /analytics/metrics/timeseries [get]
func HandleGetMetricTimeseries(c *gin.Context) {
	metricID := c.DefaultQuery("metric", "views")
	timeframeParam := c.DefaultQuery("timeframe", "24:hours")

	timeframe := parseTimeframe(timeframeParam)

	timeseries, err := internal.GetMetricTimeseries(metricID, timeframe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, timeseries)
}

// HandleGetMetricBreakdown godoc
// @Summary Get metric breakdown
// @Description Returns metric values broken down by a dimension (e.g. browser, country, operating_system, device)
// @Tags analytics
// @Produce json
// @Param metric query string false "Metric ID" default(views)
// @Param group_by query string false "Dimension to group by (browser, country, operating_system, device)" default(browser)
// @Param timeframe query string false "Timeframe filter" default(24:hours)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /analytics/metrics/breakdown [get]
func HandleGetMetricBreakdown(c *gin.Context) {
	metricID := c.DefaultQuery("metric", "views")
	groupBy := c.DefaultQuery("group_by", "browser")
	timeframeParam := c.DefaultQuery("timeframe", "24:hours")

	timeframe := parseTimeframe(timeframeParam)

	breakdown, err := internal.ListMetricBreakdown(metricID, groupBy, timeframe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, breakdown)
}

func parseTimeframe(raw string) []string {
	parts := strings.Split(raw, ",")
	var timeframe []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			timeframe = append(timeframe, trimmed)
		}
	}
	return timeframe
}
