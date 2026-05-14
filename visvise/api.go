package visvise

// VisviseAPI provides atomic API methods
type VisviseAPI struct {
	http *HTTPClient
}

// NewVisviseAPI creates a new VisviseAPI instance
func NewVisviseAPI(httpClient *HTTPClient) *VisviseAPI {
	return &VisviseAPI{http: httpClient}
}

// GetCosCred retrieves COS temporary credentials for direct file upload
func (api *VisviseAPI) GetCosCred(isTemp bool) (*GetCosCredResult, error) {
	body := map[string]interface{}{}
	if isTemp {
		body["is_temp"] = true
	}

	data, err := api.http.Post("openapi/weaver/resource/get_cos_cred", body)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	return parseCosCredResult(data), nil
}

func parseCosCredResult(data interface{}) *GetCosCredResult {
	if m, ok := data.(map[string]interface{}); ok {
		result := &GetCosCredResult{}
		if v, ok := m["start_time"].(float64); ok {
			result.StartTime = int64(v)
		}
		if v, ok := m["expired_time"].(float64); ok {
			result.ExpiredTime = int64(v)
		}
		if v, ok := m["bucket"].(string); ok {
			result.Bucket = v
		}
		if v, ok := m["region"].(string); ok {
			result.Region = v
		}
		if v, ok := m["path_prefix"].(string); ok {
			result.PathPrefix = v
		}
		if credRaw, ok := m["cred"].(map[string]interface{}); ok {
			result.Cred = CosCred{
				TmpSecretID:  getString(credRaw, "tmp_secret_id", ""),
				TmpSecretKey: getString(credRaw, "tmp_secret_key", ""),
				SessionToken: getString(credRaw, "session_token", ""),
			}
		}
		return result
	}
	return nil
}

// GetUserQuota retrieves the remaining generation count for the current API key
func (api *VisviseAPI) GetUserQuota() (*UserQuota, error) {
	data, err := api.http.Post("openapi/weaver/resource/get_user_quota", map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	if m, ok := data.(map[string]interface{}); ok {
		return &UserQuota{
			Quota:    int(getFloat64(m, "quota", 0)),
			ServerTS: int64(getFloat64(m, "server_ts", 0)),
		}, nil
	}
	return nil, nil
}

// Gen3DModel creates a 3D generation task (async)
func (api *VisviseAPI) Gen3DModel(
	name string,
	nodeType int,
	params map[string]interface{},
	inputView *View,
	inputModel string,
	inputModelFormat string,
	inputVideo string,
) ([]string, error) {
	body := map[string]interface{}{
		"name":      name,
		"node_type": nodeType,
		"params":    params,
	}

	if inputView != nil {
		body["input_view"] = inputView.ToMap()
	}
	if inputModel != "" {
		body["input_model"] = inputModel
	}
	if inputModelFormat != "" {
		body["input_model_format"] = inputModelFormat
	}
	if inputVideo != "" {
		body["input_video"] = inputVideo
	}

	data, err := api.http.Post("openapi/weaver/resource/gen_3d_model", body)
	if err != nil {
		return nil, err
	}

	if m, ok := data.(map[string]interface{}); ok {
		if modelIDs, ok := m["model_ids"].([]interface{}); ok {
			result := make([]string, len(modelIDs))
			for i, id := range modelIDs {
				if s, ok := id.(string); ok {
					result[i] = s
				}
			}
			return result, nil
		}
	}
	return nil, nil
}

// GenMultiViews generates multi-view images from a single image (async)
func (api *VisviseAPI) GenMultiViews(name string, inputView *View, params map[string]interface{}) (string, error) {
	body := map[string]interface{}{
		"name":       name,
		"input_view": inputView.ToMap(),
		"params":     params,
	}

	data, err := api.http.Post("openapi/weaver/resource/gen_multi_views", body)
	if err != nil {
		return "", err
	}

	if m, ok := data.(map[string]interface{}); ok {
		if modelID, ok := m["model_id"].(string); ok {
			return modelID, nil
		}
	}
	return "", nil
}

// GetModelList retrieves the model asset list
func (api *VisviseAPI) GetModelList(
	modelIDList []string,
	nodeTypeList []int,
	statusList []int,
	keyword string,
	limit int,
	page int,
) ([]ModelInfo, int, error) {
	body := map[string]interface{}{
		"limit": limit,
		"page":  page,
	}

	if len(modelIDList) > 0 {
		body["model_id_list"] = modelIDList
	}
	if len(nodeTypeList) > 0 {
		body["node_type_list"] = nodeTypeList
	}
	if len(statusList) > 0 {
		body["status_list"] = statusList
	}
	if keyword != "" {
		body["keyword"] = keyword
	}

	data, err := api.http.Post("openapi/weaver/resource/get_model_list", body)
	if err != nil {
		return nil, 0, err
	}

	models, totalCount := ParseModelList(data.(map[string]interface{}))
	return models, totalCount, nil
}

// ListAlgorithmModel retrieves the list of algorithm models for a given node type
func (api *VisviseAPI) ListAlgorithmModel(nodeType int, subType *int) ([]string, error) {
	body := map[string]interface{}{
		"node_type": nodeType,
	}
	if subType != nil {
		body["type"] = *subType
	}

	data, err := api.http.Post("openapi/weaver/resource/list_algorithm_model", body)
	if err != nil {
		return nil, err
	}

	if m, ok := data.(map[string]interface{}); ok {
		if modelList, ok := m["model_list"].([]interface{}); ok {
			result := make([]string, len(modelList))
			for i, model := range modelList {
				if s, ok := model.(string); ok {
					result[i] = s
				}
			}
			return result, nil
		}
	}
	return nil, nil
}

// DownloadModel generates a signed download URL for a model asset
func (api *VisviseAPI) DownloadModel(modelID string) (string, error) {
	body := map[string]interface{}{
		"model_id": modelID,
	}

	data, err := api.http.Post("openapi/weaver/resource/download_model", body)
	if err != nil {
		return "", err
	}

	if s, ok := data.(string); ok {
		return s, nil
	}
	return "", nil
}

// DeleteModel deletes a single model asset
func (api *VisviseAPI) DeleteModel(modelID string) error {
	body := map[string]interface{}{
		"model_id": modelID,
	}
	_, err := api.http.Post("openapi/weaver/resource/delete_model", body)
	return err
}

// BatchDeleteModel batch deletes model assets
func (api *VisviseAPI) BatchDeleteModel(modelIDs []string) error {
	body := map[string]interface{}{
		"model_ids": modelIDs,
	}
	_, err := api.http.Post("openapi/weaver/resource/batch_delete_model", body)
	return err
}

// RemoveBackground removes the background from an image
func (api *VisviseAPI) RemoveBackground(imageURL string) (string, error) {
	body := map[string]interface{}{
		"image_url": imageURL,
	}

	data, err := api.http.Post("openapi/weaver/resource/remove_background", body)
	if err != nil {
		return "", err
	}

	if m, ok := data.(map[string]interface{}); ok {
		if imageURL, ok := m["image_url"].(string); ok {
			return imageURL, nil
		}
	}
	return "", nil
}

// BatchGenPose batch generates poses from images (async)
func (api *VisviseAPI) BatchGenPose(
	name string,
	inputModel string,
	inputImages []string,
	params map[string]interface{},
) ([]string, error) {
	body := map[string]interface{}{
		"name":         name,
		"input_model":  inputModel,
		"input_images": inputImages,
		"params":       params,
	}

	data, err := api.http.Post("openapi/weaver/resource/batch_gen_pose", body)
	if err != nil {
		return nil, err
	}

	if m, ok := data.(map[string]interface{}); ok {
		if modelIDs, ok := m["model_ids"].([]interface{}); ok {
			result := make([]string, len(modelIDs))
			for i, id := range modelIDs {
				if s, ok := id.(string); ok {
					result[i] = s
				}
			}
			return result, nil
		}
	}
	return nil, nil
}

// GetText2MotionPromptList retrieves the text-to-motion prompt demo list
func (api *VisviseAPI) GetText2MotionPromptList(language string) ([]string, error) {
	body := map[string]interface{}{
		"language": language,
	}

	data, err := api.http.Post("openapi/weaver/demo/get_text2motion_prompt_list", body)
	if err != nil {
		return nil, err
	}

	if m, ok := data.(map[string]interface{}); ok {
		if promptList, ok := m["prompt_list"].([]interface{}); ok {
			result := make([]string, len(promptList))
			for i, p := range promptList {
				if s, ok := p.(string); ok {
					result[i] = s
				}
			}
			return result, nil
		}
	}
	return nil, nil
}

// InitSegment initializes 2D segmentation (SSE streaming interface)
func (api *VisviseAPI) InitSegment(
	name string,
	algorithmModel string,
	modelID string,
	inputView *View,
	splitType *SegmentSplitType,
	granularity *SegmentGranularity,
	prompt string,
	readTimeout int,
) (*SSEIterator, error) {
	body := map[string]interface{}{
		"name":            name,
		"algorithm_model": algorithmModel,
	}

	if modelID != "" {
		body["model_id"] = modelID
	}
	if inputView != nil {
		body["input_view"] = inputView.ToMap()
	}
	if splitType != nil {
		body["split_type"] = *splitType
	}
	if granularity != nil {
		body["granularity"] = *granularity
	}
	if prompt != "" {
		body["prompt"] = prompt
	}

	return api.http.PostSSE("openapi/weaver/component/init_segment", body, readTimeout)
}
