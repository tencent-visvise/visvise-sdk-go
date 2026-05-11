# VISVISE Weaver Go SDK

中文 | **[English](README_EN.md)**

Go SDK for the VISVISE Weaver OpenAPI. It provides:

- All atomic API methods (1:1 mapping to OpenAPI endpoints)
- High-level `Gen*()` methods for each node type (auto-upload files + create tasks)
- `WaitModel()` async polling helper

---

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Client Initialization](#client-initialization)
- [Enum Constants](#enum-constants)
- [High-Level Methods](#high-level-methods)
  - [Gen360 — Image to 360](#gen360--image-to-360)
  - [GenHighModel — Image to High-poly](#genhighmodel--image-to-high-poly)
  - [GenMidModel — Image to Mid-poly](#genmidmodel--image-to-mid-poly)
  - [GenLowModel — Image to Low-poly](#genlowmodel--image-to-low-poly)
  - [GenMeshRefine — Mesh Refinement](#genmeshrefine--mesh-refinement)
  - [GenRetopology — Retopology](#genretopology--retopology)
  - [GenLOD — LOD](#genlod--lod)
  - [GenUV — UV Unwrap](#genuv--uv-unwrap)
  - [GenTexture — Texture Generation](#gentexture--texture-generation)
  - [GenRigging — Rigging](#genrigging--rigging)
  - [GenSkinning — Skinning](#genskinning--skinning)
  - [GenVideoMotion — Video to Animation](#genvideomotion--video-to-animation)
  - [GenTextMotion — Text to Animation](#gentextmotion--text-to-animation)
  - [GenPose — Image to Pose](#genpose--image-to-pose)
  - [GenSegment2D — 2D Segmentation](#gensegment2d--2d-segmentation)
  - [WaitModel — Wait for Completion](#waitmodel--wait-for-completion)
- [Atomic API Methods](#atomic-api-methods)
- [Exceptions](#exceptions)
- [Full Workflow Examples](#full-workflow-examples)

---

## Installation

Install using `go get`:

```bash
go get github.com/tencent-visvise/visvise-sdk-go
```

---

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/visvise/visvise-sdk-go/visvise"
)

func main() {
    client := visvise.NewClient(
        "your_app_id",
        "your_secret_key",
        "your_uid",
        visvise.EnvProd,
        30,
    )

    // 1) Image-to-360: upload local image, generate multi-views
    mvModelID, err := client.Gen360(
        "character.png", // local path, SDK uploads automatically
        "",
        "my_360",
        "character.png",
        nil,
        "",
        nil,
        "",
        nil,
        "",
        nil,
        "",
    )
    if err != nil {
        panic(err)
    }

    // 2) Wait for completion, fetch multi-view output
    mvInfo, err := client.WaitModel(mvModelID, &visvise.WaitOptions{Interval: 3, Timeout: 300})
    if err != nil {
        panic(err)
    }
    outputView := mvInfo.ImageGen360Output.OutputView

    // 3) Image-to-high-poly (pass COS URLs directly)
    highModelID, err := client.GenHighModel(
        outputView.MainView,
        "",
        visvise.OutputModelFormatFBX,
        1, // FaceType.TRIANGLE
        "high_model",
        "",
        nil,
        outputView.BackView,
        "",
        outputView.LeftView,
        "",
        outputView.RightView,
        "",
    )
    if err != nil {
        panic(err)
    }

    // 4) Wait for completion
    modelInfo, err := client.WaitModel(highModelID, &visvise.WaitOptions{Timeout: 900})
    if err != nil {
        panic(err)
    }
    fmt.Println("Output model:", modelInfo.OutputModel)
}
```

---

## Client Initialization

```go
import "github.com/visvise/visvise-sdk-go/visvise"

client := visvise.NewClient(
    "your_app_id",     // required, assigned by platform
    "your_secret_key", // required, assigned by platform
    "your_uid",        // required, the user ID of the account that obtained the key
    visvise.EnvProd,  // optional, default: production
    30,               // optional, per-request HTTP timeout in seconds (default 30)
)
```

| Parameter | Required | Description |
|---|---|---|
| `appID` | ✅ | Client identifier assigned by the platform |
| `secretKey` | ✅ | Signing key assigned by the platform |
| `uid` | ✅ | User ID, obtained from the account that registered the key |
| `env` | — | Environment: `EnvProd` (default) / `EnvTest` / `EnvDev` |
| `timeout` | — | Per-request HTTP timeout in seconds (default 30) |

---

## Enum Constants

The SDK exposes the following constants. Prefer them over hard-coded numbers/strings:

```go
import "github.com/visvise/visvise-sdk-go/visvise"

// Face type
visvise.FaceTypeTriangle // 1 - triangle faces
visvise.FaceTypeQuad     // 2 - quad faces

// Detail level (for retopology)
visvise.DetailLevelLow    // 1 - low
visvise.DetailLevelMedium // 2 - medium
visvise.DetailLevelHigh   // 3 - high

// Output model format
visvise.OutputModelFormatFBX // "fbx"
visvise.OutputModelFormatOBJ // "obj"
visvise.OutputModelFormatGLB // "glb"

// Mesh refine mode
visvise.MeshRefineModeOptimize // 1 - mesh optimization
visvise.MeshRefineModeDensify  // 2 - mesh densification

// 2D segmentation split type
visvise.SegmentSplitTypeFrontView // 1 - front-view split (default)
visvise.SegmentSplitTypeFourView  // 2 - four-view split

// 2D segmentation granularity
visvise.SegmentGranularityCoarse  // 1 - coarse
visvise.SegmentGranularityMedium  // 2 - medium (default)
visvise.SegmentGranularityFine    // 3 - fine
```

---

## High-Level Methods

High-level methods bundle "COS upload + async task creation" into a single call. Pass either a local file path or a VISVISE COS URL; each method returns a `model_id`.

> **About `algorithmModel`:** Every `Gen*` method's `algorithmModel` is optional. When omitted, the SDK calls `ListAlgorithmModel` and uses the first available model for the current account.

### Gen360 — Image to 360

Generate 360-degree multi-views from a single image.

```go
modelID, err := client.Gen360(
    "path/to/character.png",  // required, main view (local path or VISVISE COS URL)
    "",                       // optional, e.g. "hunyuan3D-MultiView-v3.0"; auto-selected if omitted
    "gen_360",                // optional, task name
    "character.png",          // required, filename
    nil,                      // optional, enable A-Pose (True/False)
    "",                       // optional, style (VISVISE proprietary models only)
    nil,                      // optional, back view to improve quality
    "",                       // optional, back_view filename
    nil,                      // optional, left view
    "",                       // optional, left_view filename
    nil,                      // optional, right view
    "",                       // optional, right_view filename
)
```

---

### GenHighModel — Image to High-poly

Generate a high-poly 3D model from images / multi-views (node_type=3).

```go
modelID, err := client.GenHighModel(
    "path/to/main.png",                  // required, main view
    "",                                   // optional, e.g. "hunyuan3D-v3.1"; auto-selected if omitted
    visvise.OutputModelFormatFBX,         // optional, output format (default fbx)
    visvise.FaceTypeTriangle,             // optional, face type (default triangle)
    "gen_high_model",                     // optional, task name
    "main.png",                           // required, filename
    nil,                                  // optional, target face count (1000-1500000); auto-tuned if omitted
    nil,                                  // optional, back view to improve quality
    "",
    nil,                                  // optional, left view
    "",
    nil,                                  // optional, right view
    "",
)
```

---

### GenMidModel — Image to Mid-poly

Mid-poly generation requires all four views (node_type=11).

```go
modelID, err := client.GenMidModel(
    "path/to/main.png",                  // required (all four views are required)
    "path/to/back.png",                   // required
    "path/to/left.png",                   // required
    "path/to/right.png",                  // required
    "",                                   // optional, e.g. "VISVISE-MeshGen-V1.0.0"
    visvise.OutputModelFormatFBX,         // optional, output format
    visvise.FaceTypeTriangle,             // optional, face type
    "gen_mid_model",                      // optional, task name
    "main.png",                           // required, filename
    "back.png",
    "left.png",
    "right.png",
    "",                                   // optional, 2D segmentation asset ID (mid-poly only)
)
```

---

### GenLowModel — Image to Low-poly

Low-poly only needs the main view (node_type=13).

```go
modelID, err := client.GenLowModel(
    "path/to/main.png",                  // required, main view
    "",                                   // optional, e.g. "Tripo-v1.0-fast"
    visvise.OutputModelFormatFBX,         // optional, output format
    visvise.FaceTypeTriangle,             // optional, face type
    "gen_low_model",                      // optional, task name
    "main.png",                           // required, filename
    nil,                                  // optional, back view
    "",
    nil,                                  // optional, left view
    "",
    nil,                                  // optional, right view
    "",
)
```

---

### GenMeshRefine — Mesh Refinement

Mesh-line refinement (node_type=10).

```go
modelID, err := client.GenMeshRefine(
    "path/to/model.fbx",               // required, input model
    "",                                 // optional, e.g. "VISVISE-MeshRefine-V1.0.0"
    visvise.OutputModelFormatFBX,       // optional, input format (default fbx)
    "gen_mesh_refine",                  // optional, task name
    "model.fbx",                        // required, filename
)
```

---

### GenRetopology — Retopology

Retopology of high-poly models (node_type=1).

> Note: pass `detailLevel` for Hunyuan models, `faceNum` for VISVISE proprietary models — choose one.

```go
modelID, err := client.GenRetopology(
    "path/to/model.fbx",               // required, input model
    "",                                 // optional, e.g. "hunyuan3D-RTP-v1.5"
    visvise.OutputModelFormatFBX,       // optional, output format
    visvise.FaceTypeQuad,               // optional, face type (default quad)
    "gen_retopology",                    // optional, task name
    "model.fbx",                         // required, filename
    nil,                                 // optional (Hunyuan): DetailLevel.LOW/MEDIUM/HIGH
    nil,                                 // optional (VISVISE): target face count
)
```

---

### GenLOD — LOD

Generate level-of-detail meshes (node_type=2), with multi-shot support.

```go
reduceFaces := []visvise.ReduceFace{
    {ReduceLevel: 1, ReducePercent: 50, FaceType: visvise.FaceTypeQuad},
    {ReduceLevel: 2, ReducePercent: 25, FaceType: visvise.FaceTypeQuad},
}

modelIDs, err := client.GenLOD(
    "path/to/model.fbx",               // required, input model
    "",                                 // optional, e.g. "VISVISE-LOD-V1.0.0"
    visvise.OutputModelFormatFBX,       // optional, output format
    "gen_lod",                           // optional, task name
    "model.fbx",                         // required, filename
    reduceFaces,                         // required, reduction config list
)
```

---

### GenUV — UV Unwrap

Automatic UV unwrap (node_type=9).

```go
modelID, err := client.GenUV(
    "path/to/model.fbx",               // required, input model
    "",                                 // optional, e.g. "hunyuan3D-UV-v2.0"
    visvise.OutputModelFormatFBX,       // optional, output format
    "gen_uv",                            // optional, task name
    "model.fbx",                         // required, filename
    nil,                                 // optional, enable auto-smoothing
)
```

---

### GenTexture — Texture Generation

Generate textures for a model (node_type=8).

> At least one of `inputView.MainView` or `prompt` is required; both can be supplied together.

```go
view := &visvise.View{MainView: "path/to/ref.png"}

modelID, err := client.GenTexture(
    "path/to/model.fbx",               // required, input model
    "",                                 // optional, e.g. "hunyuan3D-TEX-v2.0"
    visvise.OutputModelFormatFBX,       // optional, output format
    "gen_texture",                       // optional, task name
    "model.fbx",                         // required, filename
    view,                                // optional, reference view (or use prompt instead)
    nil,                                 // optional, resolution (e.g. 1024 / 2048)
    nil,                                 // optional, also unwrap UV
    "",                                  // optional, text prompt for the texture
)
```

---

### GenRigging — Rigging

Auto-rigging (node_type=5). The SDK packages the raw model + JSON parameters into a zip automatically — no manual zipping required.

```go
modelID, err := client.GenRigging(
    "path/to/model.fbx",               // required, raw model (auto-zipped by SDK)
    "",                                 // optional, e.g. "VISVISE-GoRigging-V1.0.0"
    visvise.OutputModelFormatFBX,       // optional, output format
    "gen_rigging",                       // optional, task name
    "model.fbx",                         // required, filename
    "humanoid",                          // optional, "humanoid" (default) or "tetrapod"
    "",                                  // optional, template skeleton
)
```

---

### GenSkinning — Skinning

Auto-skinning (node_type=6). The SDK packages the rigged model + selection JSON into a zip automatically.

```go
meshNames := []string{"Body_Mesh", "Hair_Mesh"}
jointNames := []string{"Bip001", "Bip001 Pelvis"}

modelID, err := client.GenSkinning(
    "path/to/rigged_model.fbx",        // required, rigged model
    "",                                 // optional, e.g. "VISVISE-GoSkinning-V1.0.0"
    visvise.OutputModelFormatFBX,       // optional, output format
    "gen_skinning",                      // optional, task name
    "rigged_model.fbx",                 // required, filename
    meshNames,                           // required, meshes to skin
    jointNames,                          // required, joints to skin
)
```

---

### GenVideoMotion — Video to Animation

Drive a 3D model from motion extracted from a video (node_type=4).

```go
modelID, err := client.GenVideoMotion(
    "path/to/model.fbx",               // required, input model
    "",                                 // optional, e.g. "VISVISE-FramingAI-Base-V1.5.0"
    visvise.OutputModelFormatFBX,       // optional, output format
    "gen_video_motion",                  // optional, task name
    "model.fbx",                         // required, filename
    "path/to/dance.mp4",                 // required, source video
    nil,                                 // optional, enable hand capture
    nil,                                 // optional, enable multi-person capture
    nil,                                 // optional, rotation axis-angle [x, y, z] (radians)
)
```

---

### GenTextMotion — Text to Animation

Generate animation from a text prompt; returns 4 candidate models (node_type=4).

```go
modelIDs, err := client.GenTextMotion(
    "path/to/model.fbx",               // required, input model
    "a person breakdancing",            // required, animation prompt
    "",                                 // optional, e.g. "VISVISE-TextMotion-V1.1.0"
    visvise.OutputModelFormatFBX,       // optional, output format
    "gen_text_motion",                   // optional, task name
    "model.fbx",                         // required, filename
)
// modelIDs contains 4 IDs, wait for whichever you prefer
```

---

### GenPose — Image to Pose

Generate pose models from reference images (up to 10).

```go
inputImages := []visvise.FileInput{
    "path/to/pose_ref_1.png",
    "path/to/pose_ref_2.png",
}
imageFilenames := []string{"pose_ref_1.png", "pose_ref_2.png"}

modelIDs, err := client.GenPose(
    "path/to/model.fbx",               // required, FBX model
    inputImages,                        // required, reference images (1-10)
    "",                                 // optional, e.g. "VISVISE-PosingAI-V1.0.0"
    visvise.OutputModelFormatFBX,       // optional, output format
    "gen_pose",                          // optional, task name
    "model.fbx",                         // required, model filename
    imageFilenames,                      // optional, filenames for each input_images entry
)
```

---

### GenSegment2D — 2D Segmentation

Component segmentation over multi-views from Gen360 (node_type=14, SSE protocol). The resulting `model_id` can be passed as `segmentModelID` for `GenMidModel` / `GenLowModel`.

```go
onThinking := func(content string) {
    fmt.Println("[thinking]", content)
}

segModelID, err := client.GenSegment2D(
    "Model202604xxxxxx",               // required (or pass input_view), the Gen360 model_id
    "",                                 // optional, e.g. "VISVISE-Seg2D-V1.0.0"
    "gen_segment_2d",                   // optional, task name
    nil,                                // optional (or pass model_id_360), input view
    nil,                                // optional, SegmentSplitTypeFrontView(1, default) / FourView(2)
    nil,                                // optional, SegmentGranularityCoarse(1) / Medium(2, default) / Fine(3)
    "",                                 // optional, natural-language splitting rule (max 200 chars)
    nil,                                // optional, callback for thinking events
)
// Use the result as segmentModelID for GenMidModel / GenLowModel
```

---

### WaitModel — Wait for Completion

Poll until an async task finishes; returns `ModelInfo`.

```go
modelInfo, err := client.WaitModel(
    "Model2026033100192028",
    &visvise.WaitOptions{
        Interval: 2,  // poll interval in seconds (default 2)
        Timeout: 600, // max wait in seconds (default 600)
    },
)

fmt.Println(modelInfo.OutputModel) // output model URL
fmt.Println(modelInfo.TimeCost)    // elapsed seconds
```

**Exceptions:**

- `PollingTimeoutError` — raised when the timeout is reached
- `ModelGenerationError` — raised when the task fails (status=4)
- `InvalidParamsError` — raised immediately on parameter errors during polling (no retry)
- Other network/business errors are **silently retried**

---

## Atomic API Methods

Access low-level endpoints via `client.GetAPI().xxx()`:

```go
api := client.GetAPI()

// Get temporary upload credentials
cred, err := api.GetCosCred()

// Query remaining quota
quota, err := api.GetUserQuota()
fmt.Println(quota.Quota) // remaining count

// Fetch model list
models, total, err := api.GetModelList(
    []string{"Model2026..."},
    nil, nil, "", 10, 1,
)

// Fetch algorithm models for a node type
algModels, err := api.ListAlgorithmModel(&visvise.ListAlgorithmModelParams{NodeType: 4, SubType: 1})

// Get download URL
url, err := api.DownloadModel("Model2026...")

// Delete a single model
err = api.DeleteModel("Model2026...")

// Batch delete
err = api.BatchDeleteModel([]string{"Model2026...", "Model2026..."})

// Remove background
outURL, err := api.RemoveBg("https://cos.../image.png")

// Text-to-motion prompt suggestions
prompts, err := api.GetText2MotionPromptList("en")
```

---

## Exceptions

All SDK exceptions inherit from `WeaverError`; you can catch the base class or any subclass.

| Exception | Code | Description |
|---|---|---|
| `WeaverError` | any | Base exception |
| `NetworkError` | — | Connection / timeout errors |
| `SignatureError` | 400 | Signature failure |
| `InvalidParamsError` | 120008 | Invalid request parameters |
| `UserNotFoundError` | 120017 | User not found |
| `PermissionDeniedError` | 120018 | Permission denied |
| `QuotaExceededError` | 120020 | Daily quota exceeded |
| `ProjectPermissionError` | 120027 | Project permission missing |
| `ServerNetworkError` | 120028 | Server network error |
| `ServerTimeoutError` | 120032 | Server processing timeout |
| `RateLimitError` | 120040 | Too many requests |
| `ModelGenerationError` | — | Task failed (status=4) |
| `PollingTimeoutError` | — | WaitModel timed out |

```go
import (
    "fmt"
    "github.com/visvise/visvise-sdk-go/visvise"
)

client := visvise.NewClient("...", "...", "...", visvise.EnvProd, 30)

modelID, err := client.Gen360("image.png", "", "test", "image.png", nil, "", nil, "", nil, "", nil, "")
if err != nil {
    if _, ok := err.(*visvise.QuotaExceededError); ok {
        fmt.Println("Daily quota exceeded; please try again tomorrow")
    } else if _, ok := err.(*visvise.PollingTimeoutError); ok {
        fmt.Println("Timeout waiting for model")
    } else {
        fmt.Printf("API error: %v\n", err)
    }
}
```

---

## Full Workflow Examples

### Example 1: Image → High-poly (Gen360 + GenHighModel)

```go
package main

import (
    "fmt"
    "github.com/visvise/visvise-sdk-go/visvise"
)

func main() {
    client := visvise.NewClient("...", "...", "...", visvise.EnvProd, 30)

    // Step 1: Image-to-360
    fmt.Println("Step 1: generating multi-views...")
    mvID, _ := client.Gen360("character.png", "", "gen_360", "character.png", nil, "", nil, "", nil, "", nil, "")
    mv, _ := client.WaitModel(mvID, &visvise.WaitOptions{Interval: 3, Timeout: 300})
    views := mv.ImageGen360Output.OutputView

    // Step 2: High-poly model
    fmt.Println("Step 2: generating high-poly model...")
    highID, _ := client.GenHighModel(
        views.MainView, "", visvise.OutputModelFormatFBX,
        visvise.FaceTypeTriangle, "gen_high_model", "",
        nil, views.BackView, "", views.LeftView, "", views.RightView, "",
    )
    highModel, _ := client.WaitModel(highID, &visvise.WaitOptions{Timeout: 900})
    fmt.Println("High-poly download URL:", highModel.OutputModel)
}
```

---

### Example 2: Animation pipeline (Rigging → Skinning → Animation)

```go
package main

import (
    "fmt"
    "github.com/visvise/visvise-sdk-go/visvise"
)

func main() {
    client := visvise.NewClient("...", "...", "...", visvise.EnvProd, 30)

    // Step 1: Rigging
    rigID, _ := client.GenRigging(
        "character.fbx", "", visvise.OutputModelFormatFBX,
        "gen_rigging", "character.fbx", "humanoid", "",
    )
    rig, _ := client.WaitModel(rigID, &visvise.WaitOptions{Timeout: 600})
    fmt.Println("Rigged model:", rig.OutputModel)

    // Step 2: Skinning
    skinID, _ := client.GenSkinning(
        "rigged_character.fbx", "", visvise.OutputModelFormatFBX,
        "gen_skinning", "rigged_character.fbx",
        []string{"Body_Mesh"}, []string{"Bip001", "Bip001 Pelvis"},
    )
    client.WaitModel(skinID, &visvise.WaitOptions{Timeout: 600})

    // Step 3: Video-to-animation
    animID, _ := client.GenVideoMotion(
        "skinned_model.fbx", "", visvise.OutputModelFormatFBX,
        "gen_video_motion", "skinned_model.fbx", "dance.mp4", nil, nil, nil,
    )
    anim, _ := client.WaitModel(animID, &visvise.WaitOptions{Timeout: 900})
    fmt.Println("Animation download URL:", anim.OutputModel)
}
```

---

### Example 3: LOD generation (with multi-shot)

```go
package main

import (
    "fmt"
    "github.com/visvise/visvise-sdk-go/visvise"
)

func main() {
    client := visvise.NewClient("...", "...", "...", visvise.EnvProd, 30)

    reduceFaces := []visvise.ReduceFace{
        {ReduceLevel: 1, ReducePercent: 50, FaceType: visvise.FaceTypeQuad},
        {ReduceLevel: 2, ReducePercent: 25, FaceType: visvise.FaceTypeQuad},
    }

    modelIDs, _ := client.GenLOD(
        "high_model.fbx", "", visvise.OutputModelFormatFBX,
        "gen_lod", "high_model.fbx", reduceFaces,
    )

    // Wait for all variants
    for _, mid := range modelIDs {
        r, _ := client.WaitModel(mid, &visvise.WaitOptions{Timeout: 300})
        fmt.Println(r.ModelID, r.OutputModel)
    }
}
```

---

## License

MIT License
