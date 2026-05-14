package visvise

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tencentyun/cos-go-sdk-v5"
)

// FileInput represents the input file type
type FileInput interface{}

// StringFileInput is a string file path or COS URL
type StringFileInput string

// BytesFileInput is raw bytes content
type BytesFileInput []byte

// isCosURL checks if a string is a VISVISE platform COS URL
func isCosURL(s string) bool {
	return strings.HasPrefix(s, "https://") && strings.Contains(s, ".myqcloud.com") && strings.Contains(s, ".cos.")
}

// isLocalFile checks if a string is an existing local file path
func isLocalFile(s string) bool {
	info, err := os.Stat(s)
	return err == nil && !info.IsDir()
}

// sniffExtension detects file type from binary content using magic bytes
func sniffExtension(data []byte, defaultExt string) string {
	if len(data) == 0 {
		if !strings.HasPrefix(defaultExt, ".") {
			defaultExt = "." + defaultExt
		}
		return defaultExt
	}

	// PNG image
	if len(data) >= 8 && bytes.Equal(data[:8], []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}) {
		return ".png"
	}
	// JPEG image
	if len(data) >= 3 && bytes.Equal(data[:3], []byte{0xFF, 0xD8, 0xFF}) {
		return ".jpg"
	}
	// GIF image
	if len(data) >= 6 && (bytes.Equal(data[:6], []byte("GIF87a")) || bytes.Equal(data[:6], []byte("GIF89a"))) {
		return ".gif"
	}
	// BMP image
	if len(data) >= 2 && bytes.Equal(data[:2], []byte("BM")) {
		return ".bmp"
	}
	// WebP image
	if len(data) >= 12 && bytes.Equal(data[:4], []byte("RIFF")) && bytes.Equal(data[8:12], []byte("WEBP")) {
		return ".webp"
	}
	// TIFF image
	if len(data) >= 4 && (bytes.Equal(data[:4], []byte{0x49, 0x49, 0x2A, 0x00}) || bytes.Equal(data[:4], []byte{0x4D, 0x4D, 0x00, 0x2A})) {
		return ".tiff"
	}

	// FBX binary (Kaydara FBX Binary first 16 chars: "Kaydara FBX Bina")
	if len(data) >= 16 && string(data[:16]) == "Kaydara FBX Bina" {
		return ".fbx"
	}
	// GLB (binary glTF)
	if len(data) >= 4 && bytes.Equal(data[:4], []byte("glTF")) {
		return ".glb"
	}
	// FBX ASCII (starts with comment or FBXHeaderExtension)
	if len(data) >= 1 && data[0] == ';' && bytes.Contains(data[:200], []byte("FBX")) {
		return ".fbx"
	}
	if len(data) >= 18 && string(data[:18]) == "FBXHeaderExtension" {
		return ".fbx"
	}
	// OBJ text format
	head := bytes.TrimLeft(data[:64], " \t\r\n")
	if bytes.HasPrefix(head, []byte("v ")) ||
		bytes.HasPrefix(head, []byte("vn ")) ||
		bytes.HasPrefix(head, []byte("vt ")) ||
		bytes.HasPrefix(head, []byte("f ")) ||
		bytes.HasPrefix(head, []byte("o ")) ||
		bytes.HasPrefix(head, []byte("g ")) ||
		bytes.HasPrefix(head, []byte("mtllib ")) {
		return ".obj"
	}

	// MP4/MOV video (ftyp box at offset 4)
	if len(data) >= 12 && bytes.Equal(data[4:8], []byte("ftyp")) {
		brand := data[8:12]
		if bytes.Equal(brand, []byte("qt  ")) || bytes.Equal(brand, []byte("moov")) {
			return ".mov"
		}
		return ".mp4"
	}
	// WebM video
	if len(data) >= 4 && bytes.Equal(data[:4], []byte{0x1A, 0x45, 0xDF, 0xA3}) {
		return ".webm"
	}
	// AVI video
	if len(data) >= 12 && bytes.Equal(data[:4], []byte("RIFF")) && bytes.Equal(data[8:12], []byte("AVI ")) {
		return ".avi"
	}

	// ZIP archive
	if isZip(data) {
		return ".zip"
	}

	// glTF JSON text
	stripped := bytes.TrimLeft(data[:128], " \t\r\n")
	if bytes.HasPrefix(stripped, []byte("{")) && bytes.Contains(data[:256], []byte(`"asset"`)) && bytes.Contains(data[:512], []byte(`"version"`)) {
		return ".gltf"
	}
	// OBJ comment line
	if bytes.HasPrefix(stripped, []byte("#")) && (bytes.Contains(data[:512], []byte("\nv ")) || bytes.Contains(data[:512], []byte("\no "))) {
		return ".obj"
	}

	// Fallback to default
	if !strings.HasPrefix(defaultExt, ".") {
		defaultExt = "." + defaultExt
	}
	return defaultExt
}

// genRandomFilename generates a unique filename with given extension
func genRandomFilename(suffix string) string {
	if !strings.HasPrefix(suffix, ".") {
		suffix = "." + suffix
	}
	return uuid.New().String() + suffix
}

// Client is the main entry point for the VISVISE Weaver SDK
type Client struct {
	http *HTTPClient
	api  *VisviseAPI
}

// NewClient creates a new VISVISE client
func NewClient(appID, secretKey, uid string, env Environment, timeout int) *Client {
	httpClient := NewHTTPClient(appID, secretKey, uid, env, timeout)
	api := NewVisviseAPI(httpClient)
	return &Client{
		http: httpClient,
		api:  api,
	}
}

// NewClientDefault creates a new VISVISE client with default settings (production environment)
func NewClientDefault(appID, secretKey, uid string) *Client {
	return NewClient(appID, secretKey, uid, EnvProd, 30)
}

// GetAPI returns the underlying API instance
func (c *Client) GetAPI() *VisviseAPI {
	return c.api
}

// resolveFile resolves a file input to a COS URL
func (c *Client) resolveFile(source FileInput, isTemp bool) (string, error) {
	// If source is a string
	if s, ok := source.(string); ok {
		// Check if it's a local file path
		if isLocalFile(s) {
			return c.uploadFile(s, filepath.Base(s), isTemp)
		}
		// Not a local file, check if it's a COS URL
		if isCosURL(s) {
			return s, nil
		}
		// Not a local file and not a COS URL
		return "", fmt.Errorf("local file not found: %s. If passing a COS URL, it must be a VISVISE platform COS URL (format: https://{bucket}.cos.{region}.myqcloud.com/...)", s)
	}

	// If source is bytes
	if b, ok := source.([]byte); ok {
		filename := genRandomFilename(sniffExtension(b, ".bin"))
		return c.uploadBytes(b, filename, isTemp)
	}

	// If source is io.Reader
	if r, ok := source.(io.Reader); ok {
		data, err := io.ReadAll(r)
		if err != nil {
			return "", fmt.Errorf("failed to read data: %w", err)
		}
		filename := genRandomFilename(sniffExtension(data, ".bin"))
		return c.uploadBytes(data, filename, isTemp)
	}

	return "", fmt.Errorf("unsupported file input type")
}

func (c *Client) uploadFile(path string, filename string, isTemp bool) (string, error) {
	cred, err := c.api.GetCosCred(isTemp)
	if err != nil {
		return "", err
	}

	// Read file content
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return c.uploadWithCred(cred, data, filename, isTemp)
}

func (c *Client) uploadBytes(data []byte, filename string, isTemp bool) (string, error) {
	cred, err := c.api.GetCosCred(isTemp)
	if err != nil {
		return "", err
	}

	return c.uploadWithCred(cred, data, filename, isTemp)
}

func (c *Client) uploadWithCred(cred *GetCosCredResult, data []byte, filename string, isTemp bool) (string, error) {
	if filename == "" {
		filename = fmt.Sprintf("upload_%d.bin", time.Now().UnixNano()%1000000)
	}

	cosKey := strings.TrimRight(cred.PathPrefix, "/") + "/" + filename

	// Initialize COS client with session token transport
	bucketURL, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", cred.Bucket, cred.Region))
	baseURL := &cos.BaseURL{BucketURL: bucketURL}

	httpClient := &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     cred.Cred.TmpSecretID,
			SecretKey:    cred.Cred.TmpSecretKey,
			SessionToken: cred.Cred.SessionToken,
		},
	}

	client := cos.NewClient(baseURL, httpClient)

	// Upload file
	opt := &cos.ObjectPutOptions{}
	if isTemp {
		opt.ObjectPutHeaderOptions = &cos.ObjectPutHeaderOptions{
			XCosStorageClass: "STANDARD_IA",
		}
	}

	_, err := client.Object.Put(context.Background(), cosKey, bytes.NewReader(data), opt)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	cosURL := fmt.Sprintf("https://%s.cos.%s.myqcloud.com/%s", cred.Bucket, cred.Region, cosKey)
	return cosURL, nil
}

// resolveModelFile resolves a model file to a COS URL, handling zip packaging automatically
func (c *Client) resolveModelFile(source FileInput, isTemp bool) (string, error) {
	// If source is a string
	if s, ok := source.(string); ok {
		// Check if it's a local file path
		if isLocalFile(s) {
			ext := strings.ToLower(filepath.Ext(s))
			modelExtensions := map[string]bool{".fbx": true, ".obj": true, ".glb": true, ".gltf": true}

			if modelExtensions[ext] {
				// Model file needs to be zipped
				data, err := os.ReadFile(s)
				if err != nil {
					return "", fmt.Errorf("failed to read file: %w", err)
				}
				srcFilename := filepath.Base(s)
				stem := strings.TrimSuffix(srcFilename, filepath.Ext(srcFilename))
				zipFilename := stem + ".zip"
				return c.uploadZip(data, srcFilename, zipFilename, isTemp)
			}
			// .zip or other format, upload directly
			return c.uploadFile(s, filepath.Base(s), isTemp)
		}
		// Not a local file, check if it's a COS URL
		if isCosURL(s) {
			return s, nil
		}
		// Not a local file and not a COS URL
		return "", fmt.Errorf("local file not found: %s. If passing a COS URL, it must be a VISVISE platform COS URL (format: https://{bucket}.cos.{region}.myqcloud.com/...)", s)
	}

	var data []byte

	switch v := source.(type) {
	case []byte:
		data = v

	case io.Reader:
		var err error
		data, err = io.ReadAll(v)
		if err != nil {
			return "", fmt.Errorf("failed to read data: %w", err)
		}
	}

	// Check if data is already a zip
	if isZip(data) {
		filename := genRandomFilename(".zip")
		return c.uploadBytes(data, filename, isTemp)
	}

	// Not a zip: detect extension (fbx/obj/glb etc), default to .fbx and package
	filename := genRandomFilename(sniffExtension(data, ".fbx"))
	stem := strings.TrimSuffix(filename, filepath.Ext(filename))
	zipFilename := stem + ".zip"
	return c.uploadZip(data, filename, zipFilename, isTemp)
}

func isZip(data []byte) bool {
	return len(data) >= 4 && data[0] == 0x50 && data[1] == 0x4B && data[2] == 0x03 && data[3] == 0x04
}

func (c *Client) uploadZip(data []byte, innerFilename, zipFilename string, isTemp bool) (string, error) {
	// Create zip in memory
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)

	_, err := zw.Create(innerFilename)
	if err != nil {
		return "", fmt.Errorf("failed to create zip entry: %w", err)
	}

	w, err := zw.Create(innerFilename)
	if err != nil {
		return "", fmt.Errorf("failed to create zip entry: %w", err)
	}
	_, err = w.Write(data)
	if err != nil {
		return "", fmt.Errorf("failed to write to zip: %w", err)
	}

	if err := zw.Close(); err != nil {
		return "", fmt.Errorf("failed to close zip: %w", err)
	}

	return c.uploadBytes(buf.Bytes(), zipFilename, isTemp)
}

// buildModelZip builds a model zip with JSON parameters
func (c *Client) buildModelZip(source FileInput, jsonData map[string]interface{}, filename string) ([]byte, string, error) {
	var data []byte
	var srcFilename string

	switch v := source.(type) {
	case string:
		if !isLocalFile(v) {
			return nil, "", fmt.Errorf("model_path only accepts local file path or binary content, not COS URL: %s", v)
		}
		var err error
		data, err = os.ReadFile(v)
		if err != nil {
			return nil, "", fmt.Errorf("failed to read file: %w", err)
		}
		srcFilename = filename
		if srcFilename == "" {
			srcFilename = filepath.Base(v)
		}

	case []byte:
		data = v
		srcFilename = filename
		if srcFilename == "" {
			// Auto-detect extension for bytes input
			srcFilename = genRandomFilename(sniffExtension(data, ".fbx"))
		}

	default:
		var err error
		data, err = io.ReadAll(v.(io.Reader))
		if err != nil {
			return nil, "", fmt.Errorf("failed to read data: %w", err)
		}
		srcFilename = filename
		if srcFilename == "" {
			if name, ok := v.(interface{ Name() string }); ok {
				srcFilename = filepath.Base(name.Name())
			}
		}
		if srcFilename == "" {
			// Auto-detect extension for reader input
			srcFilename = genRandomFilename(sniffExtension(data, ".fbx"))
		}
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	var newFiles map[string][]byte

	// If data is already a zip, extract and repack
	if isZip(data) {
		zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			return nil, "", fmt.Errorf("failed to read zip: %w", err)
		}

		existing := make(map[string][]byte)
		for _, f := range zr.File {
			rc, err := f.Open()
			if err != nil {
				return nil, "", fmt.Errorf("failed to open zip entry: %w", err)
			}
			content, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return nil, "", fmt.Errorf("failed to read zip entry: %w", err)
			}
			existing[f.Name] = content
		}

		// Find model file (first non-.json entry)
		var modelEntry string
		for name := range existing {
			if !strings.HasSuffix(strings.ToLower(name), ".json") {
				modelEntry = name
				break
			}
		}

		if modelEntry == "" {
			return nil, "", errors.New("no model file found in zip (non-.json entry)")
		}

		stem := strings.TrimSuffix(modelEntry, filepath.Ext(modelEntry))
		jsonEntry := stem + ".json"

		newFiles = make(map[string][]byte)
		for k, v := range existing {
			if k != jsonEntry {
				newFiles[k] = v
			}
		}
		newFiles[jsonEntry] = jsonBytes

	} else {
		// Single model file, package directly
		stem := strings.TrimSuffix(srcFilename, filepath.Ext(srcFilename))
		jsonEntry := stem + ".json"

		newFiles = map[string][]byte{
			srcFilename: data,
			jsonEntry:   jsonBytes,
		}
	}

	// Build zip
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)
	for name, content := range newFiles {
		w, err := zw.Create(name)
		if err != nil {
			return nil, "", fmt.Errorf("failed to create zip entry: %w", err)
		}
		_, err = w.Write(content)
		if err != nil {
			return nil, "", fmt.Errorf("failed to write to zip: %w", err)
		}
	}
	if err := zw.Close(); err != nil {
		return nil, "", fmt.Errorf("failed to close zip: %w", err)
	}

	// Determine zip filename
	var zipFilename string
	if len(newFiles) == 2 {
		for name := range newFiles {
			if strings.HasSuffix(strings.ToLower(name), ".zip") {
				zipFilename = name
				break
			}
		}
	}
	if zipFilename == "" {
		zipFilename = "model.zip"
	}

	return buf.Bytes(), zipFilename, nil
}

// resolveAlgorithmModel resolves the algorithm model name
func (c *Client) resolveAlgorithmModel(algorithmModel string, nodeType NodeType, subType *int) (string, error) {
	if algorithmModel != "" {
		return algorithmModel, nil
	}

	models, err := c.api.ListAlgorithmModel(int(nodeType), subType)
	if err != nil {
		return "", err
	}

	if len(models) == 0 {
		return "", fmt.Errorf("no available algorithm model for node_type=%d. Please apply through the platform or specify algorithm_model manually", nodeType)
	}

	log.Printf("algorithm_model not specified, auto-selecting first available model: %s", models[0])
	return models[0], nil
}

// WaitModel polls and waits for model generation to complete
func (c *Client) WaitModel(modelID string, opts *WaitOptions) (*ModelInfo, error) {
	if opts == nil {
		opts = DefaultWaitOptions()
	}

	start := time.Now()
	interval := opts.GetInterval()

	for {
		elapsed := time.Since(start)
		if elapsed >= time.Duration(opts.Timeout)*time.Second {
			return nil, NewPollingTimeoutError(modelID, opts.Timeout)
		}

		models, _, err := c.api.GetModelList([]string{modelID}, nil, nil, "", 10, 1)
		if err != nil {
			// Network errors: log and continue retry
			if _, ok := err.(*NetworkError); ok {
				log.Printf("API error during polling (will retry): %v", err)
				time.Sleep(interval)
				continue
			}
			// Other errors: return immediately
			return nil, err
		}

		if len(models) == 0 {
			log.Printf("Model %s not found, continuing to wait...", modelID)
			time.Sleep(interval)
			continue
		}

		model := models[0]

		if model.Status == int(ModelStatusSuccess) {
			log.Printf("Model %s generated successfully (time: %ds) output_model=%s",
				modelID, model.TimeCost, model.OutputModel)
			return &model, nil
		}

		if model.Status == int(ModelStatusFailed) {
			reason := ""
			code := -1
			if model.FailedReason != nil {
				reason = model.FailedReason.Reason
				code = model.FailedReason.Code
			}
			return nil, NewModelGenerationError(
				fmt.Sprintf("Model generation failed: %s", reason),
				code,
				modelID,
				"",
			)
		}

		// PENDING or RUNNING
		statusName := "pending"
		if model.Status == int(ModelStatusRunning) {
			statusName = "generating"
		}
		log.Printf("Model %s %s, remaining time: %ds (waited: %.0fs)",
			modelID, statusName, model.RemainingTime, elapsed.Seconds())
		time.Sleep(interval)
	}
}

// ════════════════════════════════════════════════════════════════════
// High-level methods: gen_xxx (simplified with Options)
// ════════════════════════════════════════════════════════════════════

// Gen360 generates multi-view images from an image
// Simplified version using Gen360Options
func (c *Client) Gen360(mainView FileInput, opts *Gen360Options) (string, error) {
	if opts == nil {
		opts = NewGen360Options()
	}

	// Upload views
	mainURL, err := c.resolveFile(mainView, false)
	if err != nil {
		return "", err
	}

	view := &View{MainView: mainURL}

	if opts.BackView != nil {
		backURL, err := c.resolveFile(opts.BackView, false)
		if err != nil {
			return "", err
		}
		view.BackView = backURL
	}

	if opts.LeftView != nil {
		leftURL, err := c.resolveFile(opts.LeftView, false)
		if err != nil {
			return "", err
		}
		view.LeftView = leftURL
	}

	if opts.RightView != nil {
		rightURL, err := c.resolveFile(opts.RightView, false)
		if err != nil {
			return "", err
		}
		view.RightView = rightURL
	}

	resolvedModel, err := c.resolveAlgorithmModel(opts.AlgorithmModel, NodeTypeImgTo360, nil)
	if err != nil {
		return "", err
	}

	img360 := map[string]interface{}{
		"algorithm_model":     resolvedModel,
		"output_model_format": opts.OutputModelFormat,
		"face_type":           opts.FaceType,
	}
	if opts.EnableAPose != nil {
		img360["enable_a_pose"] = *opts.EnableAPose
	}
	if opts.Style != "" {
		img360["style"] = opts.Style
	}

	return c.api.GenMultiViews(opts.Name, view, map[string]interface{}{
		"image_gen_360_params": img360,
	})
}

// GenHighModel generates a high-detail 3D model from images
// Simplified version using GenHighModelOptions
func (c *Client) GenHighModel(mainView FileInput, opts *GenHighModelOptions) (string, error) {
	if opts == nil {
		opts = NewGenHighModelOptions()
	}

	view := &View{}

	mainURL, err := c.resolveFile(mainView, false)
	if err != nil {
		return "", err
	}
	view.MainView = mainURL

	if opts.BackView != nil {
		backURL, err := c.resolveFile(opts.BackView, false)
		if err != nil {
			return "", err
		}
		view.BackView = backURL
	}

	if opts.LeftView != nil {
		leftURL, err := c.resolveFile(opts.LeftView, false)
		if err != nil {
			return "", err
		}
		view.LeftView = leftURL
	}

	if opts.RightView != nil {
		rightURL, err := c.resolveFile(opts.RightView, false)
		if err != nil {
			return "", err
		}
		view.RightView = rightURL
	}

	resolvedModel, err := c.resolveAlgorithmModel(opts.AlgorithmModel, NodeTypeImgTo3DHigh, nil)
	if err != nil {
		return "", err
	}

	imgParams := map[string]interface{}{
		"algorithm_model":     resolvedModel,
		"output_model_format": opts.OutputModelFormat,
		"face_type":           opts.FaceType,
	}
	if opts.FaceNum != nil {
		imgParams["face_num"] = *opts.FaceNum
	}

	modelIDs, err := c.api.Gen3DModel(opts.Name, int(NodeTypeImgTo3DHigh),
		map[string]interface{}{"image_gen_model_params": imgParams},
		view, "", "", "")
	if err != nil {
		return "", err
	}

	if len(modelIDs) > 0 {
		return modelIDs[0], nil
	}
	return "", nil
}

// GenMidModel generates a mid-detail 3D model from images
// Simplified version using GenMidModelOptions
func (c *Client) GenMidModel(mainView, backView, leftView, rightView FileInput, opts *GenMidModelOptions) (string, error) {
	if opts == nil {
		opts = NewGenMidModelOptions()
	}

	view := &View{}

	// Resolve main view
	mainURL, err := c.resolveFile(mainView, false)
	if err != nil {
		return "", err
	}
	view.MainView = mainURL

	// Resolve back view (required)
	backURL, err := c.resolveFile(backView, false)
	if err != nil {
		return "", err
	}
	view.BackView = backURL

	// Resolve left view (required)
	leftURL, err := c.resolveFile(leftView, false)
	if err != nil {
		return "", err
	}
	view.LeftView = leftURL

	// Resolve right view (required)
	rightURL, err := c.resolveFile(rightView, false)
	if err != nil {
		return "", err
	}
	view.RightView = rightURL

	resolvedModel, err := c.resolveAlgorithmModel(opts.AlgorithmModel, NodeTypeImgTo3DMid, nil)
	if err != nil {
		return "", err
	}

	imgParams := map[string]interface{}{
		"algorithm_model":     resolvedModel,
		"output_model_format": opts.OutputModelFormat,
		"face_type":           opts.FaceType,
	}
	if opts.SegmentModelID != "" {
		imgParams["segment_model_id"] = opts.SegmentModelID
	}

	modelIDs, err := c.api.Gen3DModel(opts.Name, int(NodeTypeImgTo3DMid),
		map[string]interface{}{"image_gen_model_params": imgParams},
		view, "", "", "")
	if err != nil {
		return "", err
	}

	if len(modelIDs) > 0 {
		return modelIDs[0], nil
	}
	return "", nil
}

// GenLowModel generates a low-detail 3D model from images
// Simplified version using GenLowModelOptions
func (c *Client) GenLowModel(mainView FileInput, opts *GenLowModelOptions) (string, error) {
	if opts == nil {
		opts = NewGenLowModelOptions()
	}

	view := &View{}

	mainURL, err := c.resolveFile(mainView, false)
	if err != nil {
		return "", err
	}
	view.MainView = mainURL

	if opts.BackView != nil {
		backURL, err := c.resolveFile(opts.BackView, false)
		if err != nil {
			return "", err
		}
		view.BackView = backURL
	}

	if opts.LeftView != nil {
		leftURL, err := c.resolveFile(opts.LeftView, false)
		if err != nil {
			return "", err
		}
		view.LeftView = leftURL
	}

	if opts.RightView != nil {
		rightURL, err := c.resolveFile(opts.RightView, false)
		if err != nil {
			return "", err
		}
		view.RightView = rightURL
	}

	resolvedModel, err := c.resolveAlgorithmModel(opts.AlgorithmModel, NodeTypeImgTo3DLow, nil)
	if err != nil {
		return "", err
	}

	imgParams := map[string]interface{}{
		"algorithm_model":     resolvedModel,
		"output_model_format": opts.OutputModelFormat,
		"face_type":           opts.FaceType,
	}

	modelIDs, err := c.api.Gen3DModel(opts.Name, int(NodeTypeImgTo3DLow),
		map[string]interface{}{"image_gen_model_params": imgParams},
		view, "", "", "")
	if err != nil {
		return "", err
	}

	if len(modelIDs) > 0 {
		return modelIDs[0], nil
	}
	return "", nil
}

// GenMeshRefine performs mesh refinement/optimization
// Simplified version using GenMeshRefineOptions
func (c *Client) GenMeshRefine(modelPath FileInput, opts *GenMeshRefineOptions) (string, error) {
	if opts == nil {
		opts = NewGenMeshRefineOptions()
	}

	cosURL, err := c.resolveModelFile(modelPath, false)
	if err != nil {
		return "", err
	}

	resolvedModel, err := c.resolveAlgorithmModel(opts.AlgorithmModel, NodeTypeMeshRefine, nil)
	if err != nil {
		return "", err
	}

	params := map[string]interface{}{
		"algorithm_model":    resolvedModel,
		"input_model_format": opts.InputModelFormat,
	}
	if opts.Mode != nil {
		params["mode"] = *opts.Mode
	}
	if opts.ColorModel != nil {
		colorURL, err := c.resolveModelFile(opts.ColorModel, false)
		if err != nil {
			return "", err
		}
		params["color_model"] = colorURL
	}

	modelIDs, err := c.api.Gen3DModel(opts.Name, int(NodeTypeMeshRefine),
		map[string]interface{}{"mesh_refine_params": params},
		nil, cosURL, "", "")
	if err != nil {
		return "", err
	}

	if len(modelIDs) > 0 {
		return modelIDs[0], nil
	}
	return "", nil
}

// GenRetopology performs re-topology on a model
// Simplified version using GenRetopologyOptions
func (c *Client) GenRetopology(modelPath FileInput, opts *GenRetopologyOptions) (string, error) {
	if opts == nil {
		opts = NewGenRetopologyOptions()
	}

	cosURL, err := c.resolveModelFile(modelPath, false)
	if err != nil {
		return "", err
	}

	resolvedModel, err := c.resolveAlgorithmModel(opts.AlgorithmModel, NodeTypeReTopology, nil)
	if err != nil {
		return "", err
	}

	params := map[string]interface{}{
		"algorithm_model":     resolvedModel,
		"output_model_format": opts.OutputModelFormat,
		"face_type":           opts.FaceType,
	}
	if opts.DetailLevel != nil {
		params["detail_level"] = *opts.DetailLevel
	}
	if opts.FaceNum != nil {
		params["face_num"] = *opts.FaceNum
	}

	modelIDs, err := c.api.Gen3DModel(opts.Name, int(NodeTypeReTopology),
		map[string]interface{}{"re_topology_params": params},
		nil, cosURL, "", "")
	if err != nil {
		return "", err
	}

	if len(modelIDs) > 0 {
		return modelIDs[0], nil
	}
	return "", nil
}

// GenLOD generates LOD (Level of Detail) models
// Simplified version using GenLODOptions
func (c *Client) GenLOD(modelPath FileInput, reduceFaces []ReduceFace, opts *GenLODOptions) ([]string, error) {
	if opts == nil {
		opts = NewGenLODOptions()
	}

	cosURL, err := c.resolveModelFile(modelPath, false)
	if err != nil {
		return nil, err
	}

	resolvedModel, err := c.resolveAlgorithmModel(opts.AlgorithmModel, NodeTypeLOD, nil)
	if err != nil {
		return nil, err
	}

	reduceFacesData := make([]map[string]interface{}, len(reduceFaces))
	for i, rf := range reduceFaces {
		reduceFacesData[i] = rf.ToMap()
	}

	return c.api.Gen3DModel(opts.Name, int(NodeTypeLOD),
		map[string]interface{}{"lod_params": map[string]interface{}{
			"algorithm_model":     resolvedModel,
			"output_model_format": opts.OutputModelFormat,
			"reduce_faces":        reduceFacesData,
			"gen_times":           opts.GenTimes,
		}},
		nil, cosURL, "", "")
}

// GenUV performs UV unwrapping on a model
// Simplified version using GenUVOptions
func (c *Client) GenUV(modelPath FileInput, opts *GenUVOptions) (string, error) {
	if opts == nil {
		opts = NewGenUVOptions()
	}

	cosURL, err := c.resolveModelFile(modelPath, false)
	if err != nil {
		return "", err
	}

	resolvedModel, err := c.resolveAlgorithmModel(opts.AlgorithmModel, NodeTypeUV, nil)
	if err != nil {
		return "", err
	}

	uvParams := map[string]interface{}{
		"algorithm_model": resolvedModel,
	}
	if opts.EnableAutoSmoothing != nil {
		uvParams["enable_auto_smoothing"] = *opts.EnableAutoSmoothing
	}

	modelIDs, err := c.api.Gen3DModel(opts.Name, int(NodeTypeUV),
		map[string]interface{}{"uv_params": uvParams},
		nil, cosURL, "", "")
	if err != nil {
		return "", err
	}

	if len(modelIDs) > 0 {
		return modelIDs[0], nil
	}
	return "", nil
}

// GenTexture generates textures for a model
// Simplified version using GenTextureOptions
func (c *Client) GenTexture(modelPath FileInput, opts *GenTextureOptions) (string, error) {
	if opts == nil {
		opts = NewGenTextureOptions()
	}

	mainViewURL := ""
	if opts.InputView != nil {
		mainViewURL = opts.InputView.MainView
	}
	if mainViewURL == "" && opts.Prompt == "" {
		return "", errors.New("gen_texture requires either input_view.main_view or prompt")
	}

	cosURL, err := c.resolveModelFile(modelPath, false)
	if err != nil {
		return "", err
	}

	resolvedModel, err := c.resolveAlgorithmModel(opts.AlgorithmModel, NodeTypeTexture, nil)
	if err != nil {
		return "", err
	}

	texParams := map[string]interface{}{
		"algorithm_model": resolvedModel,
	}
	if opts.Resolution != nil {
		texParams["resolution"] = *opts.Resolution
	}
	if opts.UnwarpUV != nil {
		texParams["unwarp_uv"] = *opts.UnwarpUV
	}
	if opts.Prompt != "" {
		texParams["prompt"] = opts.Prompt
	}

	var resolvedView *View
	if opts.InputView != nil {
		resolvedView = &View{}
		if opts.InputView.MainView != "" {
			url, err := c.resolveFile(opts.InputView.MainView, false)
			if err != nil {
				return "", err
			}
			resolvedView.MainView = url
		}
		if opts.InputView.BackView != "" {
			url, err := c.resolveFile(opts.InputView.BackView, false)
			if err != nil {
				return "", err
			}
			resolvedView.BackView = url
		}
		if opts.InputView.LeftView != "" {
			url, err := c.resolveFile(opts.InputView.LeftView, false)
			if err != nil {
				return "", err
			}
			resolvedView.LeftView = url
		}
		if opts.InputView.RightView != "" {
			url, err := c.resolveFile(opts.InputView.RightView, false)
			if err != nil {
				return "", err
			}
			resolvedView.RightView = url
		}
	}

	modelIDs, err := c.api.Gen3DModel(opts.Name, int(NodeTypeTexture),
		map[string]interface{}{"tex_params": texParams},
		resolvedView, cosURL, "", "")
	if err != nil {
		return "", err
	}

	if len(modelIDs) > 0 {
		return modelIDs[0], nil
	}
	return "", nil
}

// GenRigging performs skeleton rigging on a model
// Simplified version using GenRiggingOptions
func (c *Client) GenRigging(modelPath FileInput, opts *GenRiggingOptions) (string, error) {
	if opts == nil {
		opts = NewGenRiggingOptions()
	}

	resolvedModel, err := c.resolveAlgorithmModel(opts.AlgorithmModel, NodeTypeRigging, nil)
	if err != nil {
		return "", err
	}

	jsonData := map[string]interface{}{
		"config": map[string]interface{}{
			"mesh_category": opts.MeshCategory,
			"algo_name":     resolvedModel,
		},
	}

	zipBytes, _, err := c.buildModelZip(modelPath, jsonData, "")
	if err != nil {
		return "", err
	}

	cosURL, err := c.uploadBytes(zipBytes, "", false)
	if err != nil {
		return "", err
	}

	goRiggingParams := map[string]interface{}{
		"algorithm_model": resolvedModel,
	}
	if opts.TemplateSkeleton != nil {
		skeletonURL, err := c.resolveModelFile(opts.TemplateSkeleton, false)
		if err != nil {
			return "", err
		}
		goRiggingParams["template_skeleton"] = skeletonURL
	}

	modelIDs, err := c.api.Gen3DModel(opts.Name, int(NodeTypeRigging),
		map[string]interface{}{"go_rigging_params": goRiggingParams},
		nil, cosURL, "", "")
	if err != nil {
		return "", err
	}

	if len(modelIDs) > 0 {
		return modelIDs[0], nil
	}
	return "", nil
}

// GenSkinning performs skinning on a rigged model
// Simplified version using GenSkinningOptions
func (c *Client) GenSkinning(modelPath FileInput, opts *GenSkinningOptions) (string, error) {
	if opts == nil {
		return "", errors.New("gen_skinning requires opts with mesh_names and joint_names")
	}

	if len(opts.MeshNames) == 0 || len(opts.JointNames) == 0 {
		return "", errors.New("gen_skinning requires mesh_names and joint_names")
	}

	resolvedModel, err := c.resolveAlgorithmModel(opts.AlgorithmModel, NodeTypeSkinning, nil)
	if err != nil {
		return "", err
	}

	jsonData := map[string]interface{}{
		"config": map[string]interface{}{
			"algo_name": resolvedModel,
		},
		"selection": map[string]interface{}{
			"mesh_names":  opts.MeshNames,
			"joint_names": opts.JointNames,
		},
	}

	zipBytes, _, err := c.buildModelZip(modelPath, jsonData, "")
	if err != nil {
		return "", err
	}

	cosURL, err := c.uploadBytes(zipBytes, "", false)
	if err != nil {
		return "", err
	}

	modelIDs, err := c.api.Gen3DModel(opts.Name, int(NodeTypeSkinning),
		map[string]interface{}{},
		nil, cosURL, "", "")
	if err != nil {
		return "", err
	}

	if len(modelIDs) > 0 {
		return modelIDs[0], nil
	}
	return "", nil
}

// GenVideoMotion generates animation from video
// Simplified version using GenVideoMotionOptions
func (c *Client) GenVideoMotion(modelPath, videoPath FileInput, opts *GenVideoMotionOptions) (string, error) {
	if opts == nil {
		opts = NewGenVideoMotionOptions()
	}

	modelURL, err := c.resolveModelFile(modelPath, false)
	if err != nil {
		return "", err
	}

	videoURL, err := c.resolveFile(videoPath, false)
	if err != nil {
		return "", err
	}

	subType := int(AnimationSubTypeVideo)
	resolvedModel, err := c.resolveAlgorithmModel(opts.AlgorithmModel, NodeTypeAnimation, &subType)
	if err != nil {
		return "", err
	}

	framing := map[string]interface{}{
		"algorithm_model":     resolvedModel,
		"output_model_format": opts.OutputModelFormat,
	}
	if opts.WithHand != nil {
		framing["with_hand"] = *opts.WithHand
	}
	if opts.MultipleTrack != nil {
		framing["multiple_track"] = *opts.MultipleTrack
	}
	if len(opts.RotateAxisAngle) == 3 {
		framing["rotate_axis_angle"] = opts.RotateAxisAngle
	}

	modelIDs, err := c.api.Gen3DModel(opts.Name, int(NodeTypeAnimation),
		map[string]interface{}{"framing_ai_params": framing},
		nil, modelURL, "", videoURL)
	if err != nil {
		return "", err
	}

	if len(modelIDs) > 0 {
		return modelIDs[0], nil
	}
	return "", nil
}

// GenTextMotion generates animation from text prompts
// Simplified version using GenTextMotionOptions
func (c *Client) GenTextMotion(modelPath FileInput, prompt string, opts *GenTextMotionOptions) ([]string, error) {
	if opts == nil {
		opts = NewGenTextMotionOptions()
	}

	modelURL, err := c.resolveModelFile(modelPath, false)
	if err != nil {
		return nil, err
	}

	subType := int(AnimationSubTypeText)
	resolvedModel, err := c.resolveAlgorithmModel(opts.AlgorithmModel, NodeTypeAnimation, &subType)
	if err != nil {
		return nil, err
	}

	return c.api.Gen3DModel(opts.Name, int(NodeTypeAnimation),
		map[string]interface{}{"framing_ai_params": map[string]interface{}{
			"algorithm_model":     resolvedModel,
			"output_model_format": opts.OutputModelFormat,
			"prompt":              prompt,
		}},
		nil, modelURL, "", "")
}

// GenPose generates poses from reference images
// Simplified version using GenPoseOptions
func (c *Client) GenPose(modelPath FileInput, inputImages []FileInput, opts *GenPoseOptions) ([]string, error) {
	if opts == nil {
		opts = NewGenPoseOptions()
	}

	modelURL, err := c.resolveModelFile(modelPath, false)
	if err != nil {
		return nil, err
	}

	uploadedImages := make([]string, len(inputImages))
	for i, img := range inputImages {
		url, err := c.resolveFile(img, false)
		if err != nil {
			return nil, err
		}
		uploadedImages[i] = url
	}

	resolvedModel, err := c.resolveAlgorithmModel(opts.AlgorithmModel, NodeTypeImgToPose, nil)
	if err != nil {
		return nil, err
	}

	return c.api.BatchGenPose(opts.Name, modelURL, uploadedImages, map[string]interface{}{
		"algorithm_model":     resolvedModel,
		"output_model_format": opts.OutputModelFormat,
	})
}

// ThinkingCallback is a callback function for SSE thinking events
type ThinkingCallback func(string)

// GenSegment2D performs 2D segmentation (SSE streaming interface)
// Simplified version using GenSegment2DOptions
func (c *Client) GenSegment2D(modelID360 string, opts *GenSegment2DOptions) (string, error) {
	if opts == nil {
		opts = NewGenSegment2DOptions()
	}

	if modelID360 == "" && opts.InputView == nil {
		return "", errors.New("gen_segment_2d requires either model_id_360 or input_view")
	}

	var resolvedView *View
	if opts.InputView != nil {
		resolvedView = &View{}
		if opts.InputView.MainView != "" {
			url, err := c.resolveFile(opts.InputView.MainView, false)
			if err != nil {
				return "", err
			}
			resolvedView.MainView = url
		}
		if opts.InputView.BackView != "" {
			url, err := c.resolveFile(opts.InputView.BackView, false)
			if err != nil {
				return "", err
			}
			resolvedView.BackView = url
		}
		if opts.InputView.LeftView != "" {
			url, err := c.resolveFile(opts.InputView.LeftView, false)
			if err != nil {
				return "", err
			}
			resolvedView.LeftView = url
		}
		if opts.InputView.RightView != "" {
			url, err := c.resolveFile(opts.InputView.RightView, false)
			if err != nil {
				return "", err
			}
			resolvedView.RightView = url
		}
	}

	resolvedModel, err := c.resolveAlgorithmModel(opts.AlgorithmModel, NodeTypeSegment2D, nil)
	if err != nil {
		return "", err
	}

	iter, err := c.api.InitSegment(opts.Name, resolvedModel, modelID360, resolvedView, opts.SplitType, opts.Granularity, opts.Prompt, opts.ReadTimeout)
	if err != nil {
		return "", err
	}
	defer iter.Close()

	var newModelID string

	for {
		event, err := iter.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		switch event.Event {
		case "pre_create":
			if m, ok := event.Data.(map[string]interface{}); ok {
				if id, ok := m["model_id"].(string); ok {
					newModelID = id
					log.Printf("gen_segment_2d: pre_create model_id=%s", newModelID)
				}
			}
		case "thinking":
			if event.Data != nil {
				if s, ok := event.Data.(string); ok {
					log.Printf("gen_segment_2d thinking: %s", s)
					if opts.OnThinking != nil {
						opts.OnThinking(s)
					}
				}
			}
		case "reply":
			log.Printf("gen_segment_2d: reply received, complete")
			goto done
		case "error":
			var msg string
			var code int = -1
			if m, ok := event.Data.(map[string]interface{}); ok {
				if s, ok := m["msg"].(string); ok {
					msg = s
				}
				if f, ok := m["code"].(float64); ok {
					code = int(f)
				}
			} else {
				msg = fmt.Sprintf("%v", event.Data)
			}
			return "", NewModelGenerationError(
				fmt.Sprintf("2D segmentation failed: %s", msg),
				code,
				newModelID,
				"",
			)
		}
	}

done:
	if newModelID == "" {
		return "", NewModelGenerationError("2D segmentation did not return model_id", -1, "", "")
	}
	return newModelID, nil
}

// UploadFile uploads a local file to COS and returns the COS URL
func (c *Client) UploadFile(path string, filename string, isTemp bool) (string, error) {
	return c.uploadFile(path, filename, isTemp)
}
