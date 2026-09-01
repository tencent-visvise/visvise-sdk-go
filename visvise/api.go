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
func (api *VisviseAPI) GetCosCred(isTemp bool, isPublic bool, rtx string) (*GetCosCredResult, error) {
	body := map[string]interface{}{}
	if isTemp {
		body["is_temp"] = true
	}
	if isPublic {
		body["is_public"] = true
	}

	data, err := api.http.Post("openapi/weaver/resource/get_cos_cred", body, rtx)
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
func (api *VisviseAPI) GetUserQuota(rtx string) (*UserQuota, error) {
	data, err := api.http.Post("openapi/weaver/resource/get_user_quota", map[string]interface{}{}, rtx)
	if err != nil {
		return nil, err
	}

	if m, ok := data.(map[string]interface{}); ok {
		return &UserQuota{
			ModelQuota:           int(getFloat64(m, "model_quota", 0)),
			AnimationQuota:       int(getFloat64(m, "animation_quota", 0)),
			ServerTS:             int64(getFloat64(m, "server_ts", 0)),
			ImageProcessingQuota: int(getFloat64(m, "image_processing_quota", 0)),
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
	rtx string,
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

	data, err := api.http.Post("openapi/weaver/resource/gen_3d_model", body, rtx)
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
func (api *VisviseAPI) GenMultiViews(name string, inputView *View, params map[string]interface{}, rtx string) (string, error) {
	body := map[string]interface{}{
		"name":       name,
		"input_view": inputView.ToMap(),
		"params":     params,
	}

	data, err := api.http.Post("openapi/weaver/resource/gen_multi_views", body, rtx)
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
	lastTs uint64,
	modelTypeList []int,
	sorter *Sorter,
	rtx string,
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
	if len(modelTypeList) > 0 {
		body["model_type_list"] = modelTypeList
	}
	if lastTs > 0 {
		body["last_ts"] = lastTs
	}
	if sorter != nil && sorter.Name != "" {
		body["sorter"] = sorter
	}

	data, err := api.http.Post("openapi/weaver/resource/get_model_list", body, rtx)
	if err != nil {
		return nil, 0, err
	}

	if data == nil {
		return nil, 0, nil
	}
	models, totalCount := ParseModelList(data.(map[string]interface{}))
	return models, totalCount, nil
}

// ListAlgorithmModel retrieves the list of algorithm models for a given node type
func (api *VisviseAPI) ListAlgorithmModel(nodeType int, subType *int, rtx string) ([]string, error) {
	body := map[string]interface{}{
		"node_type": nodeType,
	}
	if subType != nil {
		body["type"] = *subType
	}

	data, err := api.http.Post("openapi/weaver/resource/list_algorithm_model", body, rtx)
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
func (api *VisviseAPI) DownloadModel(modelID string, rtx string) (string, error) {
	body := map[string]interface{}{
		"model_id": modelID,
	}

	data, err := api.http.Post("openapi/weaver/resource/download_model", body, rtx)
	if err != nil {
		return "", err
	}

	if s, ok := data.(string); ok {
		return s, nil
	}
	return "", nil
}

// DeleteModel deletes a single model asset
func (api *VisviseAPI) DeleteModel(modelID string, rtx string) error {
	body := map[string]interface{}{
		"model_id": modelID,
	}
	_, err := api.http.Post("openapi/weaver/resource/delete_model", body, rtx)
	return err
}

// BatchDeleteModel batch deletes model assets
func (api *VisviseAPI) BatchDeleteModel(modelIDs []string, rtx string) error {
	body := map[string]interface{}{
		"model_ids": modelIDs,
	}
	_, err := api.http.Post("openapi/weaver/resource/batch_delete_model", body, rtx)
	return err
}

// RegenerateModel regenerates an existing model asset in place.
// Only node_type=AUTO_LUV (2UV, NodeTypeAutoLUV) assets are supported.
// If params is nil, the server reuses the asset's previous generation parameters.
func (api *VisviseAPI) RegenerateModel(modelID string, params map[string]interface{}, rtx string) error {
	body := map[string]interface{}{
		"model_id": modelID,
	}
	if params != nil {
		body["params"] = params
	}
	_, err := api.http.Post("openapi/weaver/resource/regenerate_model", body, rtx)
	return err
}

// RemoveBackground removes the background from an image
func (api *VisviseAPI) RemoveBackground(imageURL string, rtx string) (string, error) {
	body := map[string]interface{}{
		"image_url": imageURL,
	}

	data, err := api.http.Post("openapi/weaver/resource/remove_background", body, rtx)
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

// StyleTransfer stylizes an input image and returns the processed image COS URL.
func (api *VisviseAPI) StyleTransfer(inputView string, styleType StyleType, rtx string) (string, error) {
	body := map[string]interface{}{
		"input_view": inputView,
		"style_type": styleType,
	}

	data, err := api.http.Post("openapi/weaver/resource/style_transfer", body, rtx)
	if err != nil {
		return "", err
	}

	if m, ok := data.(map[string]interface{}); ok {
		if resultImage, ok := m["result_image"].(string); ok && resultImage != "" {
			return resultImage, nil
		}
	}
	return "", nil
}

// PatterAutoRemove automatically removes surface patterns and returns the processed image COS URL.
func (api *VisviseAPI) PatterAutoRemove(inputView string, rtx string) (string, error) {
	body := map[string]interface{}{
		"input_view": inputView,
	}

	data, err := api.http.Post("openapi/weaver/resource/patter_auto_remove", body, rtx)
	if err != nil {
		return "", err
	}

	if m, ok := data.(map[string]interface{}); ok {
		if resultImage, ok := m["result_image"].(string); ok && resultImage != "" {
			return resultImage, nil
		}
	}
	return "", nil
}

// GenPreprocess saves a processed image as a 2D preprocess model asset.
func (api *VisviseAPI) GenPreprocess(
	name, inputView string,
	preprocessType PreprocessType,
	algorithmModel string,
	styleParam *StyleParam,
	removePatternParam *RemovePatternParam,
	rtx string,
) (string, error) {
	body := map[string]interface{}{
		"name":            name,
		"input_view":      inputView,
		"preprocess_type": preprocessType,
	}
	if algorithmModel != "" {
		body["algorithm_model"] = algorithmModel
	}
	if styleParam != nil {
		body["style_param"] = styleParam
	}
	if removePatternParam != nil {
		body["remove_pattern_param"] = removePatternParam
	}

	data, err := api.http.Post("openapi/weaver/resource/gen_preprocess", body, rtx)
	if err != nil {
		return "", err
	}
	if m, ok := data.(map[string]interface{}); ok {
		if modelID, ok := m["model_id"].(string); ok && modelID != "" {
			return modelID, nil
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
	rtx string,
) ([]string, error) {
	body := map[string]interface{}{
		"name":         name,
		"input_model":  inputModel,
		"input_images": inputImages,
		"params":       params,
	}

	data, err := api.http.Post("openapi/weaver/resource/batch_gen_pose", body, rtx)
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
func (api *VisviseAPI) GetText2MotionPromptList(language string, rtx string) ([]string, error) {
	body := map[string]interface{}{
		"language": language,
	}

	data, err := api.http.Post("openapi/weaver/demo/get_text2motion_prompt_list", body, rtx)
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
	rtx string,
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

	return api.http.PostSSE("openapi/weaver/component/init_segment", body, readTimeout, rtx)
}

// BeginSegment enters the segment state and specifies the component to split.
func (api *VisviseAPI) BeginSegment(clientID string, componentLabel int32, viewType SegmentViewType, rtx string) (*OperatorResult, error) {
	body := map[string]interface{}{
		"client_id":       clientID,
		"view_type":       viewType,
		"component_label": componentLabel,
	}
	data, err := api.http.Post("openapi/weaver/component/begin_segment", body, rtx)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	var result OperatorResult
	if err := decodeData(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SegmentComponent marks the region to split out in the segment state (repeatable).
func (api *VisviseAPI) SegmentComponent(
	clientID string,
	viewType SegmentViewType,
	addPixels []*Pixel,
	removePixels []*Pixel,
	rects []*Rect,
	rtx string,
) (*OperatorResult, error) {
	body := map[string]interface{}{
		"client_id": clientID,
		"view_type": viewType,
	}
	if len(addPixels) > 0 {
		body["add_pixels"] = addPixels
	}
	if len(removePixels) > 0 {
		body["remove_pixels"] = removePixels
	}
	if len(rects) > 0 {
		body["rects"] = rects
	}
	data, err := api.http.Post("openapi/weaver/component/segment", body, rtx)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	var result OperatorResult
	if err := decodeData(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ConfirmSegment finalizes the current segmentation result.
func (api *VisviseAPI) ConfirmSegment(clientID string, viewType SegmentViewType, rtx string) (*OperatorResult, error) {
	body := map[string]interface{}{
		"client_id": clientID,
		"view_type": viewType,
	}
	data, err := api.http.Post("openapi/weaver/component/confirm_segment", body, rtx)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	var result OperatorResult
	if err := decodeData(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CancelSegment cancels the current segmentation and reverts to the pre-segment state.
func (api *VisviseAPI) CancelSegment(clientID string, viewType SegmentViewType, rtx string) (*OperatorResult, error) {
	body := map[string]interface{}{
		"client_id": clientID,
		"view_type": viewType,
	}
	data, err := api.http.Post("openapi/weaver/component/cancel_segment", body, rtx)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	var result OperatorResult
	if err := decodeData(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// MergeComponent merges multiple components into one connected component.
func (api *VisviseAPI) MergeComponent(clientID string, componentLabels []int32, viewType SegmentViewType, rtx string) (*MultiViewSegmentResult, error) {
	body := map[string]interface{}{
		"client_id":        clientID,
		"component_labels": componentLabels,
		"view_type":        viewType,
	}
	data, err := api.http.Post("openapi/weaver/component/merge", body, rtx)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	var result MultiViewSegmentResult
	if err := decodeData(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AutoMergeComponent automatically merges all adjacent connected components.
func (api *VisviseAPI) AutoMergeComponent(clientID string, rtx string) (*MultiViewSegmentResult, error) {
	body := map[string]interface{}{
		"client_id": clientID,
	}
	data, err := api.http.Post("openapi/weaver/component/auto_merge", body, rtx)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	var result MultiViewSegmentResult
	if err := decodeData(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// BoundaryAdjust adjusts a component boundary via a painted mask region.
func (api *VisviseAPI) BoundaryAdjust(clientID string, viewType SegmentViewType, componentLabel int32, paintMask string, rtx string) (*OperatorResult, error) {
	body := map[string]interface{}{
		"client_id":       clientID,
		"view_type":       viewType,
		"paint_mask":      paintMask,
		"component_label": componentLabel,
	}
	data, err := api.http.Post("openapi/weaver/component/boundary_adjust", body, rtx)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	var result OperatorResult
	if err := decodeData(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// RenameComponent renames a component, syncing across all views in the four-view stage.
func (api *VisviseAPI) RenameComponent(clientID string, viewType SegmentViewType, componentLabel int32, newName string, rtx string) (*MultiViewSegmentResult, error) {
	body := map[string]interface{}{
		"client_id":       clientID,
		"view_type":       viewType,
		"component_label": componentLabel,
		"new_name":        newName,
	}
	data, err := api.http.Post("openapi/weaver/component/part_rename", body, rtx)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	var result MultiViewSegmentResult
	if err := decodeData(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SaveSegment persists the current segmentation result as a standalone 2D split asset (node_type=14).
func (api *VisviseAPI) SaveSegment(clientID string, name string, algorithmModel string, openedModelID string, rtx string) (*ModelInfo, error) {
	body := map[string]interface{}{
		"client_id":       clientID,
		"name":            name,
		"algorithm_model": algorithmModel,
	}
	if openedModelID != "" {
		body["opened_model_id"] = openedModelID
	}
	data, err := api.http.Post("openapi/weaver/component/save_segment", body, rtx)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	m, ok := data.(map[string]interface{})
	if !ok {
		return nil, nil
	}
	info := parseModelInfo(m)
	return &info, nil
}

// OpenSegment opens a saved 2D split asset for re-editing, returning a new client_id.
func (api *VisviseAPI) OpenSegment(modelID string, rtx string) (*MultiViewSegmentResult, error) {
	body := map[string]interface{}{
		"model_id": modelID,
	}
	data, err := api.http.Post("openapi/weaver/component/open_segment", body, rtx)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	var result MultiViewSegmentResult
	if err := decodeData(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
