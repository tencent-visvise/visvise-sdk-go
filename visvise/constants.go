package visvise

// NodeType represents the node type enum
type NodeType int

const (
	NodeTypeReTopology   NodeType = 1  // Re-topology
	NodeTypeLOD          NodeType = 2  // LOD
	NodeTypeImgTo3DHigh  NodeType = 3  // Image to 3D (High)
	NodeTypeAnimation    NodeType = 4  // Framing AI Animation
	NodeTypeRigging      NodeType = 5  // Skeleton setup
	NodeTypeSkinning     NodeType = 6  // Skinning
	NodeTypeImgTo360     NodeType = 7  // Image to 360
	NodeTypeTexture      NodeType = 8  // Texture
	NodeTypeUV           NodeType = 9  // UV unwrap
	NodeTypeMeshRefine   NodeType = 10 // Mesh refine
	NodeTypeImgTo3DMid   NodeType = 11 // Image to 3D (Mid)
	NodeTypeImgToPose    NodeType = 12 // Image to Pose
	NodeTypeImgTo3DLow   NodeType = 13 // Image to 3D (Low)
	NodeTypeSegment2D    NodeType = 14 // 2D Segmentation
)

// ModelStatus represents the model asset status code
type ModelStatus int

const (
	ModelStatusInvalid  ModelStatus = 0 // Invalid
	ModelStatusPending  ModelStatus = 1 // Waiting for generation
	ModelStatusRunning  ModelStatus = 2 // Generating
	ModelStatusSuccess  ModelStatus = 3 // Generation succeeded
	ModelStatusFailed   ModelStatus = 4 // Generation failed
)

// FaceType represents the face type enum
type FaceType int

const (
	FaceTypeTriangle FaceType = 1 // Triangle
	FaceTypeQuad     FaceType = 2 // Quad
)

// DetailLevel represents the detail level enum for re-topology
type DetailLevel int

const (
	DetailLevelLow    DetailLevel = 1 // Low
	DetailLevelMedium DetailLevel = 2 // Medium
	DetailLevelHigh   DetailLevel = 3 // High
)

// OutputModelFormat represents the output model format
type OutputModelFormat string

const (
	OutputModelFormatFBX OutputModelFormat = "fbx"
	OutputModelFormatOBJ OutputModelFormat = "obj"
	OutputModelFormatGLB OutputModelFormat = "glb"
)

// MeshRefineMode represents the mesh refine mode
type MeshRefineMode int

const (
	MeshRefineModeOptimize MeshRefineMode = 1 // Optimize
	MeshRefineModeDensify  MeshRefineMode = 2 // Densify
)

// SegmentSplitType represents the 2D segmentation split type
type SegmentSplitType int

const (
	SegmentSplitFrontView SegmentSplitType = 1 // Front view (default)
	SegmentSplitFourView  SegmentSplitType = 2 // Four view
)

// SegmentGranularity represents the 2D segmentation granularity
type SegmentGranularity int

const (
	SegmentGranularityCoarse  SegmentGranularity = 1 // Coarse (x50%)
	SegmentGranularityMedium  SegmentGranularity = 2 // Medium (x70%, default)
	SegmentGranularityFine    SegmentGranularity = 3 // Fine (x100%)
)

// AnimationSubType represents the animation sub-type
type AnimationSubType int

const (
	AnimationSubTypeVideo AnimationSubType = 1 // Video to animation
	AnimationSubTypeText  AnimationSubType = 2 // Text to animation
)
