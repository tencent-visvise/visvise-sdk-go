package visvise

import (
	"encoding/json"
	"time"
)

// View represents the multi-view structure
type View struct {
	MainView string `json:"main_view"`
	BackView string `json:"back_view,omitempty"`
	LeftView string `json:"left_view,omitempty"`
	RightView string `json:"right_view,omitempty"`
}

// ToMap converts View to map
func (v *View) ToMap() map[string]interface{} {
	m := map[string]interface{}{"main_view": v.MainView}
	if v.BackView != "" {
		m["back_view"] = v.BackView
	}
	if v.LeftView != "" {
		m["left_view"] = v.LeftView
	}
	if v.RightView != "" {
		m["right_view"] = v.RightView
	}
	return m
}

// ReduceFace represents the LOD single level reduce face configuration
type ReduceFace struct {
	ReduceLevel   int `json:"reduce_level"`
	ReducePercent int `json:"reduce_percent"`
	FaceType      int `json:"face_type"` // 1: Triangle, 2: Quad
}

// ToMap converts ReduceFace to map
func (r *ReduceFace) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"reduce_level":   r.ReduceLevel,
		"reduce_percent": r.ReducePercent,
		"face_type":      r.FaceType,
	}
}

// CosCred represents the COS temporary credentials
type CosCred struct {
	TmpSecretID  string `json:"tmp_secret_id"`
	TmpSecretKey string `json:"tmp_secret_key"`
	SessionToken string `json:"session_token"`
}

// GetCosCredResult represents the get_cos_cred API response
type GetCosCredResult struct {
	Cred        CosCred `json:"cred"`
	StartTime   int64   `json:"start_time"`
	ExpiredTime int64   `json:"expired_time"`
	Bucket      string  `json:"bucket"`
	Region      string  `json:"region"`
	PathPrefix  string  `json:"path_prefix"`
}

// UserQuota represents the get_user_quota API response
type UserQuota struct {
	Quota    int   `json:"quota"`
	ServerTS int64 `json:"server_ts"`
}

// FailedReason represents the generation failure reason
type FailedReason struct {
	Code   int    `json:"code"`
	Reason string `json:"reason"`
}

// LODFile represents a single LOD level output
type LODFile struct {
	ReduceLevel   int    `json:"reduce_level"`
	DownloadURL   string `json:"download_url"`
	PreviewImg    string `json:"preview_img,omitempty"`
}

// LODOutput represents the LOD output file collection
type LODOutput struct {
	LODFiles []LODFile `json:"lod_files"`
	ZipFile  string    `json:"zip_file,omitempty"`
}

// ImageGen360Output represents the image to 360 output result
type ImageGen360Output struct {
	OutputView          *View  `json:"output_view,omitempty"`
	HorizontalViewVideo string `json:"horizontal_view_video,omitempty"`
}

// Text2Motion represents a single text-to-motion output
type Text2Motion struct {
	OutputModel string `json:"output_model,omitempty"`
	PreviewImg  string `json:"preview_img,omitempty"`
}

// FramingAIOutput represents the Framing AI output result
type FramingAIOutput struct {
	Text2MotionResult []Text2Motion `json:"text2_motion_result"`
}

// ModelInfo represents the model asset information
type ModelInfo struct {
	ModelID           string           `json:"model_id"`
	Name              string           `json:"name"`
	Status            int              `json:"status"`
	NodeType          int              `json:"node_type"`
	CreateTS          int64            `json:"create_ts,omitempty"`
	CreateUser        string           `json:"create_user,omitempty"`
	PreviewImg        string           `json:"preview_img,omitempty"`
	OutputModel       string           `json:"output_model,omitempty"`
	InputModel        string           `json:"input_model,omitempty"`
	InputVideo        string           `json:"input_video,omitempty"`
	TimeCost          int              `json:"time_cost,omitempty"`
	RemainingTime     int              `json:"remaining_time,omitempty"`
	WaitTime          int              `json:"wait_time,omitempty"`
	FailedReason      *FailedReason    `json:"failed_reason,omitempty"`
	LODOutput         *LODOutput       `json:"lod_output,omitempty"`
	ImageGen360Output *ImageGen360Output `json:"image_gen_360_output,omitempty"`
	FramingAIOutput   *FramingAIOutput `json:"framing_ai_output,omitempty"`
	Params            interface{}      `json:"params,omitempty"`
	InputView         interface{}      `json:"input_view,omitempty"`
	AlgorithmModel    string           `json:"algorithm_model,omitempty"`
}

// IsSuccess returns true if the model generation succeeded
func (m *ModelInfo) IsSuccess() bool {
	return m.Status == int(ModelStatusSuccess)
}

// IsFailed returns true if the model generation failed
func (m *ModelInfo) IsFailed() bool {
	return m.Status == int(ModelStatusFailed)
}

// IsPending returns true if the model is pending or running
func (m *ModelInfo) IsPending() bool {
	return m.Status == int(ModelStatusPending) || m.Status == int(ModelStatusRunning)
}

// ParseJSON parses a JSON byte slice into the appropriate struct
func ParseJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// ParseModelList parses the model list response
func ParseModelList(data map[string]interface{}) ([]ModelInfo, int) {
	var models []ModelInfo
	modelListRaw, ok := data["model_list"]
	if ok {
		if modelList, ok := modelListRaw.([]interface{}); ok {
			for _, m := range modelList {
				if mMap, ok := m.(map[string]interface{}); ok {
					modelInfo := parseModelInfo(mMap)
					models = append(models, modelInfo)
				}
			}
		}
	}
	totalCount := 0
	if count, ok := data["total_count"].(float64); ok {
		totalCount = int(count)
	}
	return models, totalCount
}

func parseModelInfo(m map[string]interface{}) ModelInfo {
	info := ModelInfo{}
	if v, ok := m["model_id"].(string); ok {
		info.ModelID = v
	}
	if v, ok := m["name"].(string); ok {
		info.Name = v
	}
	if v, ok := m["status"].(float64); ok {
		info.Status = int(v)
	}
	if v, ok := m["node_type"].(float64); ok {
		info.NodeType = int(v)
	}
	if v, ok := m["create_ts"].(float64); ok {
		info.CreateTS = int64(v)
	}
	if v, ok := m["create_user"].(string); ok {
		info.CreateUser = v
	}
	if v, ok := m["preview_img"].(string); ok {
		info.PreviewImg = v
	}
	if v, ok := m["output_model"].(string); ok {
		info.OutputModel = v
	}
	if v, ok := m["input_model"].(string); ok {
		info.InputModel = v
	}
	if v, ok := m["input_video"].(string); ok {
		info.InputVideo = v
	}
	if v, ok := m["time_cost"].(float64); ok {
		info.TimeCost = int(v)
	}
	if v, ok := m["remaining_time"].(float64); ok {
		info.RemainingTime = int(v)
	}
	if v, ok := m["wait_time"].(float64); ok {
		info.WaitTime = int(v)
	}
	if v, ok := m["failed_reason"].(map[string]interface{}); ok {
		info.FailedReason = &FailedReason{
			Code:   int(getFloat64(v, "code", -1)),
			Reason: getString(v, "reason", ""),
		}
	}
	if v, ok := m["lod_output"].(map[string]interface{}); ok {
		info.LODOutput = parseLODOutput(v)
	}
	if v, ok := m["image_gen_360_output"].(map[string]interface{}); ok {
		info.ImageGen360Output = parseImageGen360Output(v)
	}
	if v, ok := m["framing_ai_output"].(map[string]interface{}); ok {
		info.FramingAIOutput = parseFramingAIOutput(v)
	}
	if v, ok := m["algorithm_model"].(string); ok {
		info.AlgorithmModel = v
	}
	return info
}

func parseLODOutput(m map[string]interface{}) *LODOutput {
	output := &LODOutput{}
	if v, ok := m["zip_file"].(string); ok {
		output.ZipFile = v
	}
	if lodFilesRaw, ok := m["lod_files"].([]interface{}); ok {
		for _, f := range lodFilesRaw {
			if fMap, ok := f.(map[string]interface{}); ok {
				lf := LODFile{
					ReduceLevel: int(getFloat64(fMap, "reduce_level", 0)),
					DownloadURL: getString(fMap, "download_url", ""),
					PreviewImg:  getString(fMap, "preview_img", ""),
				}
				output.LODFiles = append(output.LODFiles, lf)
			}
		}
	}
	return output
}

func parseImageGen360Output(m map[string]interface{}) *ImageGen360Output {
	output := &ImageGen360Output{}
	if v, ok := m["horizontal_view_video"].(string); ok {
		output.HorizontalViewVideo = v
	}
	if ov, ok := m["output_view"].(map[string]interface{}); ok {
		view := &View{
			MainView: getString(ov, "main_view", ""),
			BackView: getString(ov, "back_view", ""),
			LeftView: getString(ov, "left_view", ""),
			RightView: getString(ov, "right_view", ""),
		}
		if view.MainView != "" || view.BackView != "" {
			output.OutputView = view
		}
	}
	return output
}

func parseFramingAIOutput(m map[string]interface{}) *FramingAIOutput {
	output := &FramingAIOutput{}
	if resultsRaw, ok := m["text2_motion_result"].([]interface{}); ok {
		for _, r := range resultsRaw {
			if rMap, ok := r.(map[string]interface{}); ok {
				t2m := Text2Motion{
					OutputModel: getString(rMap, "output_model", ""),
					PreviewImg:  getString(rMap, "preview_img", ""),
				}
				output.Text2MotionResult = append(output.Text2MotionResult, t2m)
			}
		}
	}
	return output
}

func getString(m map[string]interface{}, key, defaultVal string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return defaultVal
}

func getFloat64(m map[string]interface{}, key string, defaultVal float64) float64 {
	if v, ok := m[key].(float64); ok {
		return v
	}
	return defaultVal
}

// SSEResult represents an SSE event frame
type SSEResult struct {
	Event string
	Data  interface{}
}

// WaitOptions represents the options for waiting for model completion
type WaitOptions struct {
	Interval float64 // Polling interval in seconds, default 2.0
	Timeout  int     // Timeout in seconds, default 600
}

// DefaultWaitOptions returns the default wait options
func DefaultWaitOptions() *WaitOptions {
	return &WaitOptions{
		Interval: 2.0,
		Timeout:  600,
	}
}

// GetInterval returns the interval, using default if 0
func (w *WaitOptions) GetInterval() time.Duration {
	if w.Interval <= 0 {
		w.Interval = 2.0
	}
	return time.Duration(w.Interval * float64(time.Second))
}
