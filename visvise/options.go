package visvise

// Gen360Options defines optional parameters for Gen360
type Gen360Options struct {
	Name              string            // optional, task name (auto-generated if empty)
	AlgorithmModel    string            // optional, algorithm model name; auto-selected if empty
	OutputModelFormat OutputModelFormat // optional, output format (default fbx)
	FaceType          FaceType          // optional, face type (default triangle)
	EnableAPose       *bool             // optional, enable A-Pose
	Style             string            // optional, style type (VISVISE proprietary models only)
	BackView          FileInput         // optional, back view to improve quality
	LeftView          FileInput         // optional, left view
	RightView         FileInput         // optional, right view
}

// NewGen360Options creates Gen360Options with common defaults
func NewGen360Options() *Gen360Options {
	return &Gen360Options{
		Name:              "gen_360",
		OutputModelFormat: OutputModelFormatFBX,
		FaceType:          FaceTypeTriangle,
	}
}

// SetName sets the task name
func (o *Gen360Options) SetName(name string) *Gen360Options {
	o.Name = name
	return o
}

// SetAlgorithmModel sets the algorithm model
func (o *Gen360Options) SetAlgorithmModel(model string) *Gen360Options {
	o.AlgorithmModel = model
	return o
}

// SetEnableAPose enables A-Pose generation
func (o *Gen360Options) SetEnableAPose(enable bool) *Gen360Options {
	o.EnableAPose = &enable
	return o
}

// SetStyle sets the style type
func (o *Gen360Options) SetStyle(style string) *Gen360Options {
	o.Style = style
	return o
}

// SetBackView sets the back view
func (o *Gen360Options) SetBackView(view FileInput) *Gen360Options {
	o.BackView = view
	return o
}

// SetLeftView sets the left view
func (o *Gen360Options) SetLeftView(view FileInput) *Gen360Options {
	o.LeftView = view
	return o
}

// SetRightView sets the right view
func (o *Gen360Options) SetRightView(view FileInput) *Gen360Options {
	o.RightView = view
	return o
}

// SetOutputModelFormat sets the output model format
func (o *Gen360Options) SetOutputModelFormat(format OutputModelFormat) *Gen360Options {
	o.OutputModelFormat = format
	return o
}

// SetFaceType sets the face type
func (o *Gen360Options) SetFaceType(faceType FaceType) *Gen360Options {
	o.FaceType = faceType
	return o
}

// GenHighModelOptions defines optional parameters for GenHighModel
type GenHighModelOptions struct {
	Name              string            // optional, task name (auto-generated if empty)
	AlgorithmModel    string            // optional, algorithm model name; auto-selected if empty
	OutputModelFormat OutputModelFormat // optional, output format (default fbx)
	FaceType          FaceType          // optional, face type (default triangle)
	FaceNum           *int              // optional, target face count (1000-1500000)
	BackView          FileInput         // optional, back view to improve quality
	LeftView          FileInput         // optional, left view
	RightView         FileInput         // optional, right view
}

// NewGenHighModelOptions creates GenHighModelOptions with common defaults
func NewGenHighModelOptions() *GenHighModelOptions {
	return &GenHighModelOptions{
		Name:              "gen_high_model",
		OutputModelFormat: OutputModelFormatFBX,
		FaceType:          FaceTypeTriangle,
	}
}

// SetName sets the task name
func (o *GenHighModelOptions) SetName(name string) *GenHighModelOptions {
	o.Name = name
	return o
}

// SetAlgorithmModel sets the algorithm model
func (o *GenHighModelOptions) SetAlgorithmModel(model string) *GenHighModelOptions {
	o.AlgorithmModel = model
	return o
}

// SetOutputModelFormat sets the output model format
func (o *GenHighModelOptions) SetOutputModelFormat(format OutputModelFormat) *GenHighModelOptions {
	o.OutputModelFormat = format
	return o
}

// SetFaceType sets the face type
func (o *GenHighModelOptions) SetFaceType(faceType FaceType) *GenHighModelOptions {
	o.FaceType = faceType
	return o
}

// SetFaceNum sets the target face count
func (o *GenHighModelOptions) SetFaceNum(faceNum int) *GenHighModelOptions {
	o.FaceNum = &faceNum
	return o
}

// SetBackView sets the back view
func (o *GenHighModelOptions) SetBackView(view FileInput) *GenHighModelOptions {
	o.BackView = view
	return o
}

// SetLeftView sets the left view
func (o *GenHighModelOptions) SetLeftView(view FileInput) *GenHighModelOptions {
	o.LeftView = view
	return o
}

// SetRightView sets the right view
func (o *GenHighModelOptions) SetRightView(view FileInput) *GenHighModelOptions {
	o.RightView = view
	return o
}

// GenMidModelOptions defines optional parameters for GenMidModel
type GenMidModelOptions struct {
	Name              string            // optional, task name (auto-generated if empty)
	AlgorithmModel    string            // optional, algorithm model name
	OutputModelFormat OutputModelFormat // optional, output format (default fbx)
	FaceType          FaceType          // optional, face type (default triangle)
	SegmentModelID    string            // optional, 2D segmentation asset ID
}

// NewGenMidModelOptions creates GenMidModelOptions with common defaults
func NewGenMidModelOptions() *GenMidModelOptions {
	return &GenMidModelOptions{
		Name:              "gen_mid_model",
		OutputModelFormat: OutputModelFormatFBX,
		FaceType:          FaceTypeTriangle,
	}
}

// SetName sets the task name
func (o *GenMidModelOptions) SetName(name string) *GenMidModelOptions {
	o.Name = name
	return o
}

// SetAlgorithmModel sets the algorithm model
func (o *GenMidModelOptions) SetAlgorithmModel(model string) *GenMidModelOptions {
	o.AlgorithmModel = model
	return o
}

// SetOutputModelFormat sets the output model format
func (o *GenMidModelOptions) SetOutputModelFormat(format OutputModelFormat) *GenMidModelOptions {
	o.OutputModelFormat = format
	return o
}

// SetFaceType sets the face type
func (o *GenMidModelOptions) SetFaceType(faceType FaceType) *GenMidModelOptions {
	o.FaceType = faceType
	return o
}

// SetSegmentModelID sets the 2D segmentation asset ID
func (o *GenMidModelOptions) SetSegmentModelID(id string) *GenMidModelOptions {
	o.SegmentModelID = id
	return o
}

// GenLowModelOptions defines optional parameters for GenLowModel
type GenLowModelOptions struct {
	Name              string            // optional, task name (auto-generated if empty)
	AlgorithmModel    string            // optional, algorithm model name
	OutputModelFormat OutputModelFormat // optional, output format (default fbx)
	FaceType          FaceType          // optional, face type (default triangle)
	BackView          FileInput         // optional, back view
	LeftView          FileInput         // optional, left view
	RightView         FileInput         // optional, right view
}

// NewGenLowModelOptions creates GenLowModelOptions with common defaults
func NewGenLowModelOptions() *GenLowModelOptions {
	return &GenLowModelOptions{
		Name:              "gen_low_model",
		OutputModelFormat: OutputModelFormatFBX,
		FaceType:          FaceTypeTriangle,
	}
}

// SetName sets the task name
func (o *GenLowModelOptions) SetName(name string) *GenLowModelOptions {
	o.Name = name
	return o
}

// SetAlgorithmModel sets the algorithm model
func (o *GenLowModelOptions) SetAlgorithmModel(model string) *GenLowModelOptions {
	o.AlgorithmModel = model
	return o
}

// SetOutputModelFormat sets the output model format
func (o *GenLowModelOptions) SetOutputModelFormat(format OutputModelFormat) *GenLowModelOptions {
	o.OutputModelFormat = format
	return o
}

// SetFaceType sets the face type
func (o *GenLowModelOptions) SetFaceType(faceType FaceType) *GenLowModelOptions {
	o.FaceType = faceType
	return o
}

// SetBackView sets the back view
func (o *GenLowModelOptions) SetBackView(view FileInput) *GenLowModelOptions {
	o.BackView = view
	return o
}

// SetLeftView sets the left view
func (o *GenLowModelOptions) SetLeftView(view FileInput) *GenLowModelOptions {
	o.LeftView = view
	return o
}

// SetRightView sets the right view
func (o *GenLowModelOptions) SetRightView(view FileInput) *GenLowModelOptions {
	o.RightView = view
	return o
}

// GenMeshRefineOptions defines optional parameters for GenMeshRefine
type GenMeshRefineOptions struct {
	Name             string            // optional, task name (auto-generated if empty)
	AlgorithmModel   string            // optional, algorithm model name
	InputModelFormat OutputModelFormat // optional, input model format (default fbx)
	Mode             *MeshRefineMode   // optional, MeshRefineModeOptimize(1) or MeshRefineModeDensify(2)
	ColorModel       FileInput         // optional, color model for texture preservation
}

// NewGenMeshRefineOptions creates GenMeshRefineOptions with common defaults
func NewGenMeshRefineOptions() *GenMeshRefineOptions {
	return &GenMeshRefineOptions{
		Name:             "gen_mesh_refine",
		InputModelFormat: OutputModelFormatFBX,
	}
}

// SetName sets the task name
func (o *GenMeshRefineOptions) SetName(name string) *GenMeshRefineOptions {
	o.Name = name
	return o
}

// SetAlgorithmModel sets the algorithm model
func (o *GenMeshRefineOptions) SetAlgorithmModel(model string) *GenMeshRefineOptions {
	o.AlgorithmModel = model
	return o
}

// SetInputModelFormat sets the input model format
func (o *GenMeshRefineOptions) SetInputModelFormat(format OutputModelFormat) *GenMeshRefineOptions {
	o.InputModelFormat = format
	return o
}

// SetMode sets the refine mode
func (o *GenMeshRefineOptions) SetMode(mode MeshRefineMode) *GenMeshRefineOptions {
	o.Mode = &mode
	return o
}

// SetColorModel sets the color model
func (o *GenMeshRefineOptions) SetColorModel(model FileInput) *GenMeshRefineOptions {
	o.ColorModel = model
	return o
}

// GenRetopologyOptions defines optional parameters for GenRetopology
type GenRetopologyOptions struct {
	Name              string            // optional, task name (auto-generated if empty)
	AlgorithmModel    string            // optional, algorithm model name
	OutputModelFormat OutputModelFormat // optional, output format (default fbx)
	FaceType          FaceType          // optional, face type (default quad)
	DetailLevel       *DetailLevel      // optional, for Hunyuan models: DetailLevel.LOW/MEDIUM/HIGH
	FaceNum           *int              // optional, for VISVISE models: target face count
}

// NewGenRetopologyOptions creates GenRetopologyOptions with common defaults
func NewGenRetopologyOptions() *GenRetopologyOptions {
	return &GenRetopologyOptions{
		Name:              "gen_retopology",
		OutputModelFormat: OutputModelFormatFBX,
		FaceType:          FaceTypeQuad,
	}
}

// SetName sets the task name
func (o *GenRetopologyOptions) SetName(name string) *GenRetopologyOptions {
	o.Name = name
	return o
}

// SetAlgorithmModel sets the algorithm model
func (o *GenRetopologyOptions) SetAlgorithmModel(model string) *GenRetopologyOptions {
	o.AlgorithmModel = model
	return o
}

// SetOutputModelFormat sets the output model format
func (o *GenRetopologyOptions) SetOutputModelFormat(format OutputModelFormat) *GenRetopologyOptions {
	o.OutputModelFormat = format
	return o
}

// SetFaceType sets the face type
func (o *GenRetopologyOptions) SetFaceType(faceType FaceType) *GenRetopologyOptions {
	o.FaceType = faceType
	return o
}

// SetDetailLevel sets the detail level for Hunyuan models
func (o *GenRetopologyOptions) SetDetailLevel(level DetailLevel) *GenRetopologyOptions {
	o.DetailLevel = &level
	return o
}

// SetFaceNum sets the target face count for VISVISE models
func (o *GenRetopologyOptions) SetFaceNum(faceNum int) *GenRetopologyOptions {
	o.FaceNum = &faceNum
	return o
}

// GenLODOptions defines optional parameters for GenLOD
type GenLODOptions struct {
	Name              string            // optional, task name (auto-generated if empty)
	AlgorithmModel    string            // optional, algorithm model name
	OutputModelFormat OutputModelFormat // optional, output format (default fbx)
	GenTimes          int               // optional, number of generations (default 3)
}

// NewGenLODOptions creates GenLODOptions with common defaults
func NewGenLODOptions() *GenLODOptions {
	return &GenLODOptions{
		Name:              "gen_lod",
		OutputModelFormat: OutputModelFormatFBX,
		GenTimes:          3,
	}
}

// SetName sets the task name
func (o *GenLODOptions) SetName(name string) *GenLODOptions {
	o.Name = name
	return o
}

// SetAlgorithmModel sets the algorithm model
func (o *GenLODOptions) SetAlgorithmModel(model string) *GenLODOptions {
	o.AlgorithmModel = model
	return o
}

// SetOutputModelFormat sets the output model format
func (o *GenLODOptions) SetOutputModelFormat(format OutputModelFormat) *GenLODOptions {
	o.OutputModelFormat = format
	return o
}

// SetGenTimes sets the number of generations
func (o *GenLODOptions) SetGenTimes(times int) *GenLODOptions {
	o.GenTimes = times
	return o
}

// GenUVOptions defines optional parameters for GenUV
type GenUVOptions struct {
	Name                string // optional, task name (auto-generated if empty)
	AlgorithmModel      string // optional, algorithm model name
	EnableAutoSmoothing *bool  // optional, enable auto-smoothing
}

// NewGenUVOptions creates GenUVOptions with common defaults
func NewGenUVOptions() *GenUVOptions {
	return &GenUVOptions{
		Name: "gen_uv",
	}
}

// SetName sets the task name
func (o *GenUVOptions) SetName(name string) *GenUVOptions {
	o.Name = name
	return o
}

// SetAlgorithmModel sets the algorithm model
func (o *GenUVOptions) SetAlgorithmModel(model string) *GenUVOptions {
	o.AlgorithmModel = model
	return o
}

// SetEnableAutoSmoothing enables auto-smoothing
func (o *GenUVOptions) SetEnableAutoSmoothing(enable bool) *GenUVOptions {
	o.EnableAutoSmoothing = &enable
	return o
}

// GenTextureOptions defines optional parameters for GenTexture
type GenTextureOptions struct {
	Name           string // optional, task name (auto-generated if empty)
	AlgorithmModel string // optional, algorithm model name
	InputView      *View  // optional, reference view (required with prompt)
	Resolution     *int   // optional, resolution (e.g. 1024, 2048)
	UnwarpUV       *bool  // optional, also unwrap UV
	Prompt         string // optional, text prompt for texture
}

// NewGenTextureOptions creates GenTextureOptions with common defaults
func NewGenTextureOptions() *GenTextureOptions {
	return &GenTextureOptions{
		Name: "gen_texture",
	}
}

// SetName sets the task name
func (o *GenTextureOptions) SetName(name string) *GenTextureOptions {
	o.Name = name
	return o
}

// SetAlgorithmModel sets the algorithm model
func (o *GenTextureOptions) SetAlgorithmModel(model string) *GenTextureOptions {
	o.AlgorithmModel = model
	return o
}

// SetInputView sets the reference view
func (o *GenTextureOptions) SetInputView(view *View) *GenTextureOptions {
	o.InputView = view
	return o
}

// SetResolution sets the texture resolution
func (o *GenTextureOptions) SetResolution(res int) *GenTextureOptions {
	o.Resolution = &res
	return o
}

// SetUnwarpUV enables UV unwrapping
func (o *GenTextureOptions) SetUnwarpUV(unwarp bool) *GenTextureOptions {
	o.UnwarpUV = &unwarp
	return o
}

// SetPrompt sets the text prompt
func (o *GenTextureOptions) SetPrompt(prompt string) *GenTextureOptions {
	o.Prompt = prompt
	return o
}

// GenRiggingOptions defines optional parameters for GenRigging
type GenRiggingOptions struct {
	Name             string       // optional, task name (auto-generated if empty)
	AlgorithmModel   string       // optional, algorithm model name
	MeshCategory     MeshCategory // optional, MeshCategoryHumanoid (default) or MeshCategoryTetrapod
	TemplateSkeleton FileInput    // optional, template skeleton
}

// NewGenRiggingOptions creates GenRiggingOptions with common defaults
func NewGenRiggingOptions() *GenRiggingOptions {
	return &GenRiggingOptions{
		Name:         "gen_rigging",
		MeshCategory: MeshCategoryHumanoid,
	}
}

// SetName sets the task name
func (o *GenRiggingOptions) SetName(name string) *GenRiggingOptions {
	o.Name = name
	return o
}

// SetAlgorithmModel sets the algorithm model
func (o *GenRiggingOptions) SetAlgorithmModel(model string) *GenRiggingOptions {
	o.AlgorithmModel = model
	return o
}

// SetMeshCategory sets the mesh category
func (o *GenRiggingOptions) SetMeshCategory(category MeshCategory) *GenRiggingOptions {
	o.MeshCategory = category
	return o
}

// SetTemplateSkeleton sets the template skeleton
func (o *GenRiggingOptions) SetTemplateSkeleton(skeleton FileInput) *GenRiggingOptions {
	o.TemplateSkeleton = skeleton
	return o
}

// GenSkinningOptions defines optional parameters for GenSkinning
type GenSkinningOptions struct {
	Name           string   // optional, task name (auto-generated if empty)
	AlgorithmModel string   // optional, algorithm model name
	MeshNames      []string // required, meshes to skin
	JointNames     []string // required, joints to skin
}

// NewGenSkinningOptions creates GenSkinningOptions
func NewGenSkinningOptions(meshNames, jointNames []string) *GenSkinningOptions {
	return &GenSkinningOptions{
		Name:       "gen_skinning",
		MeshNames:  meshNames,
		JointNames: jointNames,
	}
}

// SetName sets the task name
func (o *GenSkinningOptions) SetName(name string) *GenSkinningOptions {
	o.Name = name
	return o
}

// SetAlgorithmModel sets the algorithm model
func (o *GenSkinningOptions) SetAlgorithmModel(model string) *GenSkinningOptions {
	o.AlgorithmModel = model
	return o
}

// SetMeshNames sets the mesh names to skin
func (o *GenSkinningOptions) SetMeshNames(names []string) *GenSkinningOptions {
	o.MeshNames = names
	return o
}

// SetJointNames sets the joint names to skin
func (o *GenSkinningOptions) SetJointNames(names []string) *GenSkinningOptions {
	o.JointNames = names
	return o
}

// GenVideoMotionOptions defines optional parameters for GenVideoMotion
type GenVideoMotionOptions struct {
	Name              string            // optional, task name (auto-generated if empty)
	AlgorithmModel    string            // optional, algorithm model name
	OutputModelFormat OutputModelFormat // optional, output format (default fbx)
	WithHand          *bool             // optional, enable hand capture
	MultipleTrack     *bool             // optional, enable multi-person capture
	RotateAxisAngle   []float64         // optional, rotation axis-angle [x, y, z] (radians)
}

// NewGenVideoMotionOptions creates GenVideoMotionOptions with common defaults
func NewGenVideoMotionOptions() *GenVideoMotionOptions {
	return &GenVideoMotionOptions{
		Name:              "gen_video_motion",
		OutputModelFormat: OutputModelFormatFBX,
	}
}

// SetName sets the task name
func (o *GenVideoMotionOptions) SetName(name string) *GenVideoMotionOptions {
	o.Name = name
	return o
}

// SetAlgorithmModel sets the algorithm model
func (o *GenVideoMotionOptions) SetAlgorithmModel(model string) *GenVideoMotionOptions {
	o.AlgorithmModel = model
	return o
}

// SetOutputModelFormat sets the output model format
func (o *GenVideoMotionOptions) SetOutputModelFormat(format OutputModelFormat) *GenVideoMotionOptions {
	o.OutputModelFormat = format
	return o
}

// SetWithHand enables hand capture
func (o *GenVideoMotionOptions) SetWithHand(enable bool) *GenVideoMotionOptions {
	o.WithHand = &enable
	return o
}

// SetMultipleTrack enables multi-person capture
func (o *GenVideoMotionOptions) SetMultipleTrack(enable bool) *GenVideoMotionOptions {
	o.MultipleTrack = &enable
	return o
}

// SetRotateAxisAngle sets the rotation axis-angle
func (o *GenVideoMotionOptions) SetRotateAxisAngle(x, y, z float64) *GenVideoMotionOptions {
	o.RotateAxisAngle = []float64{x, y, z}
	return o
}

// GenTextMotionOptions defines optional parameters for GenTextMotion
type GenTextMotionOptions struct {
	Name              string            // optional, task name (auto-generated if empty)
	AlgorithmModel    string            // optional, algorithm model name
	OutputModelFormat OutputModelFormat // optional, output format (default fbx)
}

// NewGenTextMotionOptions creates GenTextMotionOptions with common defaults
func NewGenTextMotionOptions() *GenTextMotionOptions {
	return &GenTextMotionOptions{
		Name:              "gen_text_motion",
		OutputModelFormat: OutputModelFormatFBX,
	}
}

// SetName sets the task name
func (o *GenTextMotionOptions) SetName(name string) *GenTextMotionOptions {
	o.Name = name
	return o
}

// SetAlgorithmModel sets the algorithm model
func (o *GenTextMotionOptions) SetAlgorithmModel(model string) *GenTextMotionOptions {
	o.AlgorithmModel = model
	return o
}

// SetOutputModelFormat sets the output model format
func (o *GenTextMotionOptions) SetOutputModelFormat(format OutputModelFormat) *GenTextMotionOptions {
	o.OutputModelFormat = format
	return o
}

// GenPoseOptions defines optional parameters for GenPose
type GenPoseOptions struct {
	Name              string            // optional, task name (auto-generated if empty)
	AlgorithmModel    string            // optional, algorithm model name
	OutputModelFormat OutputModelFormat // optional, output format (default fbx)
}

// NewGenPoseOptions creates GenPoseOptions with common defaults
func NewGenPoseOptions() *GenPoseOptions {
	return &GenPoseOptions{
		Name:              "gen_pose",
		OutputModelFormat: OutputModelFormatFBX,
	}
}

// SetName sets the task name
func (o *GenPoseOptions) SetName(name string) *GenPoseOptions {
	o.Name = name
	return o
}

// SetAlgorithmModel sets the algorithm model
func (o *GenPoseOptions) SetAlgorithmModel(model string) *GenPoseOptions {
	o.AlgorithmModel = model
	return o
}

// SetOutputModelFormat sets the output model format
func (o *GenPoseOptions) SetOutputModelFormat(format OutputModelFormat) *GenPoseOptions {
	o.OutputModelFormat = format
	return o
}

// GenSegment2DOptions defines optional parameters for GenSegment2D
type GenSegment2DOptions struct {
	Name           string              // optional, task name (auto-generated if empty)
	AlgorithmModel string              // optional, algorithm model name
	InputView      *View               // optional, input view (with model_id_360)
	SplitType      *SegmentSplitType   // optional, split type
	Granularity    *SegmentGranularity // optional, granularity
	Prompt         string              // optional, natural language prompt
	OnThinking     ThinkingCallback    // optional, callback for thinking events
	ReadTimeout    int
}

// NewGenSegment2DOptions creates GenSegment2DOptions with common defaults
func NewGenSegment2DOptions() *GenSegment2DOptions {
	return &GenSegment2DOptions{
		Name: "gen_segment_2d",
	}
}

// SetName sets the task name
func (o *GenSegment2DOptions) SetName(name string) *GenSegment2DOptions {
	o.Name = name
	return o
}

// SetAlgorithmModel sets the algorithm model
func (o *GenSegment2DOptions) SetAlgorithmModel(model string) *GenSegment2DOptions {
	o.AlgorithmModel = model
	return o
}

// SetInputView sets the input view
func (o *GenSegment2DOptions) SetInputView(view *View) *GenSegment2DOptions {
	o.InputView = view
	return o
}

// SetSplitType sets the split type
func (o *GenSegment2DOptions) SetSplitType(splitType SegmentSplitType) *GenSegment2DOptions {
	o.SplitType = &splitType
	return o
}

// SetGranularity sets the granularity
func (o *GenSegment2DOptions) SetGranularity(granularity SegmentGranularity) *GenSegment2DOptions {
	o.Granularity = &granularity
	return o
}

// SetPrompt sets the prompt
func (o *GenSegment2DOptions) SetPrompt(prompt string) *GenSegment2DOptions {
	o.Prompt = prompt
	return o
}

// SetOnThinking sets the thinking callback
func (o *GenSegment2DOptions) SetOnThinking(callback ThinkingCallback) *GenSegment2DOptions {
	o.OnThinking = callback
	return o
}

// SetReadTimeout sets the readTimeout
func (o *GenSegment2DOptions) SetReadTimeout(readTimeout int) *GenSegment2DOptions {
	o.ReadTimeout = readTimeout
	return o
}

// ══════════════════════════════════════════════════════════════════════════════
// Client Options
// ══════════════════════════════════════════════════════════════════════════════

// ClientOptions defines optional parameters for creating a Client.
// Use NewClientOptions() to create with default values, then chain setters.
type ClientOptions struct {
	Env     Environment // optional, environment (default EnvProd)
	Timeout int         // optional, HTTP timeout in seconds (default 30)
	Debug   bool        // optional, enable debug logging (default false)
}

// NewClientOptions creates ClientOptions with default values.
// Defaults: Env=EnvProd, Timeout=30, Debug=false.
//
// 示例:
//
//	opts := visvise.NewClientOptions().SetEnv(visvise.EnvDev).SetDebug(true)
//	client := visvise.NewClient("app_id", "secret_key", "uid", opts)
func NewClientOptions() *ClientOptions {
	return &ClientOptions{
		Env:     EnvProd,
		Timeout: 30,
		Debug:   false,
	}
}

// SetEnv sets the environment
func (o *ClientOptions) SetEnv(env Environment) *ClientOptions {
	o.Env = env
	return o
}

// SetTimeout sets the HTTP timeout in seconds
func (o *ClientOptions) SetTimeout(seconds int) *ClientOptions {
	o.Timeout = seconds
	return o
}

// SetDebug enables debug logging
func (o *ClientOptions) SetDebug(enabled bool) *ClientOptions {
	o.Debug = enabled
	return o
}
