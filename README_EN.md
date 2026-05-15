# VISVISE Weaver Go SDK

**[中文](README.md)** | English

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
- [Errors](#Errors)
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
    // Create client with optional parameters
    client := visvise.NewClient(
        "your_app_id",
        "your_secret_key",
        "your_uid",
        nil, // Use default config, or pass visvise.NewClientOptions()
    )
    // Enable debug logging
    // client := visvise.NewClient("...", "...", "...", visvise.NewClientOptions().SetDebug(true))

    // 1) Image-to-360: upload local image, generate multi-views
    // name is set in Options (optional, auto-generated if omitted)
    opts := visvise.NewGen360Options().
        SetEnableAPose(true)  // optional
    mvModelID, err := client.Gen360("character.png", opts)
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
    opts = visvise.NewGenHighModelOptions().
        SetBackView(outputView.BackView).
        SetLeftView(outputView.LeftView).
        SetRightView(outputView.RightView)
    highModelID, err := client.GenHighModel(outputView.MainView, opts)
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

// Required parameters
client := visvise.NewClient(
    "your_app_id",     // required, assigned by platform
    "your_secret_key", // required, assigned by platform
    "your_uid",        // required, the user ID of the account that obtained the key
    nil,               // optional parameters, use defaults (EnvProd, Timeout=30, Debug=false)
)

// Optional parameters example
client := visvise.NewClient("...", "...", "...",
    visvise.NewClientOptions().
        SetEnv(visvise.EnvDev).  // Set environment, default is EnvProd
        SetTimeout(60).            // Set timeout, default is 30 seconds
        SetDebug(true))            // Enable debug logging, default is false
```

| Parameter | Required | Description |
|---|---|---|
| `appID` | ✅ | Client identifier assigned by the platform |
| `secretKey` | ✅ | Signing key assigned by the platform |
| `uid` | ✅ | User ID, obtained from the account that registered the key |
| `opts` | — | Optional parameters `*ClientOptions`, nil means use defaults |

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
visvise.SegmentSplitFrontView // 1 - front-view split (default)
visvise.SegmentSplitFourView  // 2 - four-view split

// 2D segmentation granularity
visvise.SegmentGranularityCoarse  // 1 - coarse
visvise.SegmentGranularityMedium  // 2 - medium (default)
visvise.SegmentGranularityFine    // 3 - fine

// Mesh category (for rigging)
visvise.MeshCategoryHumanoid // "humanoid" - humanoid (default)
visvise.MeshCategoryTetrapod // "tetrapod" - tetrapod (four-legged)
```

---

## High-Level Methods

High-level methods bundle "COS upload + async task creation" into a single call. Pass either a local file path or a VISVISE COS URL; each method returns a `model_id`.

All `Gen*()` methods use **Options struct** pattern with fluent API for cleaner optional parameter handling:

> **About `name`:** Every `Gen*` method's `name` is optional. Default values are set in `New*Options()` (e.g., `gen_360`, `gen_high_model`, etc.). Customize via `SetName()`.
>
> **About `algorithmModel`:** Every `Gen*` method's `algorithmModel` is optional. When omitted, the SDK calls `ListAlgorithmModel` and uses the first available model for the current account.

### Gen360 — Image to 360

Generate 360-degree multi-views from a single image.

```go
// Optional params via Options struct, supports fluent API
// Default name="gen_360", output format is fbx, face type is triangle
opts := visvise.NewGen360Options().
    SetName("my_360").                                    // optional, default "gen_360"
    SetAlgorithmModel("hunyuan3D-MultiView-v3.0").        // optional
    SetOutputModelFormat(visvise.OutputModelFormatFBX).  // optional, output format (default fbx)
    SetFaceType(visvise.FaceTypeTriangle).               // optional, face type (default triangle)
    SetEnableAPose(true).                                 // optional, enable A-Pose
    SetStyle("anime").                                    // optional, style type
    SetBackView("path/to/back.png").                      // optional, back view
    SetLeftView("path/to/left.png").                      // optional, left view
    SetRightView("path/to/right.png")                     // optional, right view

modelID, err := client.Gen360("path/to/character.png", opts)
```

---

### GenHighModel — Image to High-poly

Generate a high-poly 3D model from images / multi-views (node_type=3).

```go
opts := visvise.NewGenHighModelOptions().
    SetName("my_high_model").                            // optional, default "gen_high_model"
    SetAlgorithmModel("hunyuan3D-v3.1").                 // optional
    SetOutputModelFormat(visvise.OutputModelFormatFBX). // optional, output format (default fbx)
    SetFaceType(visvise.FaceTypeTriangle).               // optional, face type (default triangle)
    SetFaceNum(500000).                                  // optional, target face count (1000-1500000)
    SetBackView(outputView.BackView).                    // optional, back view
    SetLeftView(outputView.LeftView).                    // optional, left view
    SetRightView(outputView.RightView)                   // optional, right view

modelID, err := client.GenHighModel("path/to/main.png", opts)
```

---

### GenMidModel — Image to Mid-poly

Mid-poly generation requires all four views (node_type=11).

```go
opts := visvise.NewGenMidModelOptions().
    SetName("my_mid_model").                             // optional, default "gen_mid_model"
    SetAlgorithmModel("VISVISE-MeshGen-V1.0.0").          // optional
    SetOutputModelFormat(visvise.OutputModelFormatFBX). // optional, output format
    SetFaceType(visvise.FaceTypeTriangle).               // optional, face type
    SetSegmentModelID("Model2026...")                    // optional, 2D segmentation asset ID

// mainView, backView, leftView, rightView are all required
modelID, err := client.GenMidModel(
    "path/to/main.png",
    "path/to/back.png",
    "path/to/left.png",
    "path/to/right.png",
    opts,
)
```

---

### GenLowModel — Image to Low-poly

Low-poly only needs the main view (node_type=13).

```go
opts := visvise.NewGenLowModelOptions().
    SetName("my_low_model").                             // optional, default "gen_low_model"
    SetAlgorithmModel("Tripo-v1.0-fast").                // optional
    SetOutputModelFormat(visvise.OutputModelFormatFBX). // optional, output format
    SetFaceType(visvise.FaceTypeTriangle).               // optional, face type
    SetBackView("path/to/back.png").                     // optional, back view
    SetLeftView("path/to/left.png").                     // optional, left view
    SetRightView("path/to/right.png")                    // optional, right view

modelID, err := client.GenLowModel("path/to/main.png", opts)
```

---

### GenMeshRefine — Mesh Refinement

Mesh-line refinement (node_type=10).

```go
opts := visvise.NewGenMeshRefineOptions().
    SetName("my_mesh_refine").                          // optional, default "gen_mesh_refine"
    SetAlgorithmModel("VISVISE-MeshRefine-V1.0.0").      // optional
    SetInputModelFormat(visvise.OutputModelFormatFBX).  // optional, input format (default fbx)
    SetMode(visvise.MeshRefineModeOptimize).            // optional, refine mode
    SetColorModel("path/to/color.fbx")                   // optional, color model

modelID, err := client.GenMeshRefine("path/to/model.fbx", opts)
```

---

### GenRetopology — Retopology

Retopology of high-poly models (node_type=1).

> Note: pass `DetailLevel` for Hunyuan models, `FaceNum` for VISVISE proprietary models — choose one.

```go
// Hunyuan model example
opts := visvise.NewGenRetopologyOptions().
    SetName("my_retopo").                               // optional, default "gen_retopology"
    SetAlgorithmModel("hunyuan3D-RTP-v1.5").             // optional
    SetOutputModelFormat(visvise.OutputModelFormatFBX). // optional, output format
    SetFaceType(visvise.FaceTypeQuad).                  // optional, face type (default quad)
    SetDetailLevel(visvise.DetailLevelMedium)           // optional, for Hunyuan models

// VISVISE proprietary model example
opts := visvise.NewGenRetopologyOptions().
    SetAlgorithmModel("VISVISE-RTP-V1.0.0").
    SetFaceNum(50000)                                   // optional, for VISVISE models

modelID, err := client.GenRetopology("path/to/model.fbx", opts)
```

---

### GenLOD — LOD

Generate level-of-detail meshes (node_type=2), with multi-shot support. Default generation times is 3.

```go
reduceFaces := []visvise.ReduceFace{
    {ReduceLevel: 1, ReducePercent: 50, FaceType: visvise.FaceTypeQuad},
    {ReduceLevel: 2, ReducePercent: 25, FaceType: visvise.FaceTypeQuad},
}

opts := visvise.NewGenLODOptions().
    SetName("my_lod").                                  // optional, default "gen_lod"
    SetAlgorithmModel("VISVISE-LOD-V1.0.0").             // optional
    SetOutputModelFormat(visvise.OutputModelFormatFBX).  // optional, output format (default fbx)
    SetGenTimes(3)                                       // optional, number of generations (default 3)

modelIDs, err := client.GenLOD("path/to/model.fbx", reduceFaces, opts)
```

---

### GenUV — UV Unwrap

Automatic UV unwrap (node_type=9).

```go
opts := visvise.NewGenUVOptions().
    SetName("my_uv").                                    // optional, default "gen_uv"
    SetAlgorithmModel("hunyuan3D-UV-v2.0").              // optional
    SetEnableAutoSmoothing(true)                         // optional, enable auto-smoothing

modelID, err := client.GenUV("path/to/model.fbx", opts)
```

---

### GenTexture — Texture Generation

Generate textures for a model (node_type=8).

> At least one of `InputView.MainView` or `Prompt` is required; both can be supplied together.

```go
view := &visvise.View{MainView: "path/to/ref.png"}

opts := visvise.NewGenTextureOptions().
    SetName("my_texture").                              // optional, default "gen_texture"
    SetAlgorithmModel("hunyuan3D-TEX-v2.0").            // optional
    SetInputView(view).                                 // optional, reference view (or use prompt)
    SetResolution(2048).                                // optional, resolution
    SetUnwarpUV(true).                                  // optional, also unwrap UV
    SetPrompt("high quality, realistic")                  // optional, text prompt

modelID, err := client.GenTexture("path/to/model.fbx", opts)
```

---

### GenRigging — Rigging

Auto-rigging (node_type=5). The SDK packages the raw model + JSON parameters into a zip automatically — no manual zipping required.

```go
opts := visvise.NewGenRiggingOptions().
    SetName("my_rigging").                              // optional, default "gen_rigging"
    SetAlgorithmModel("VISVISE-GoRigging-V1.0.0").      // optional
    SetMeshCategory(visvise.MeshCategoryHumanoid).      // optional, humanoid (default) or visvise.MeshCategoryTetrapod
    SetTemplateSkeleton("path/to/skeleton.fbx")          // optional, template skeleton

modelID, err := client.GenRigging("path/to/model.fbx", opts)
```

---

### GenSkinning — Skinning

Auto-skinning (node_type=6). The SDK packages the rigged model + selection JSON into a zip automatically.

```go
meshNames := []string{"Body_Mesh", "Hair_Mesh"}
jointNames := []string{"Bip001", "Bip001 Pelvis"}

opts := visvise.NewGenSkinningOptions(meshNames, jointNames).
    SetName("my_skinning").                             // optional, default "gen_skinning"
    SetAlgorithmModel("VISVISE-GoSkinning-V1.0.0")      // optional

modelID, err := client.GenSkinning("path/to/rigged_model.fbx", opts)
```

---

### GenVideoMotion — Video to Animation

Drive a 3D model from motion extracted from a video (node_type=4).

```go
opts := visvise.NewGenVideoMotionOptions().
    SetName("my_video_motion").                         // optional, default "gen_video_motion"
    SetAlgorithmModel("VISVISE-FramingAI-Base-V1.5.0"). // optional
    SetOutputModelFormat(visvise.OutputModelFormatFBX). // optional, output format (default fbx)
    SetWithHand(true).                                  // optional, enable hand capture
    SetMultipleTrack(false).                            // optional, enable multi-person capture
    SetRotateAxisAngle(0, 0, 0)                         // optional, rotation axis-angle [x, y, z] (radians)

modelID, err := client.GenVideoMotion("path/to/model.fbx", "path/to/dance.mp4", opts)
```

---

### GenTextMotion — Text to Animation

Generate animation from a text prompt; returns 4 candidate models (node_type=4).

```go
opts := visvise.NewGenTextMotionOptions().
    SetName("my_text_motion").                          // optional, default "gen_text_motion"
    SetAlgorithmModel("VISVISE-TextMotion-V1.1.0").     // optional
    SetOutputModelFormat(visvise.OutputModelFormatFBX). // optional, output format

modelIDs, err := client.GenTextMotion("path/to/model.fbx", "a person breakdancing", opts)
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

opts := visvise.NewGenPoseOptions().
    SetName("my_pose").                                 // optional, default "gen_pose"
    SetAlgorithmModel("VISVISE-PosingAI-V1.0.0").       // optional
    SetOutputModelFormat(visvise.OutputModelFormatFBX)  // optional, output format

modelIDs, err := client.GenPose("path/to/model.fbx", inputImages, opts)
```

---

### GenSegment2D — 2D Segmentation

Component segmentation over multi-views from Gen360 (node_type=14, SSE protocol). The resulting `model_id` can be passed as `segmentModelID` for `GenMidModel`.

```go
onThinking := func(content string) {
    fmt.Println("[thinking]", content)
}

opts := visvise.NewGenSegment2DOptions().
    SetName("my_segment").                              // optional, default "gen_segment_2d"
    SetAlgorithmModel("VISVISE-Seg2D-V1.0.0").          // optional
    SetInputView(&visvise.View{MainView: "path/to/ref.png"}) // optional, input view
    SetSplitType(visvise.SegmentSplitFrontView).         // optional, split type
    SetGranularity(visvise.SegmentGranularityMedium).    // optional, granularity
    SetPrompt("segment by body parts").                 // optional, natural language prompt
    SetOnThinking(onThinking).                           // optional, thinking callback
    SetReadTimeout(120)                                  // optional, SSE read timeout (seconds)

segModelID, err := client.GenSegment2D("Model2026...", opts)
// Use the result as segmentModelID for GenMidModel
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

**Errors:**

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
algModels, err := api.ListAlgorithmModel(int(visvise.NodeTypeAnimation), nil)

// Get download URL
url, err := api.DownloadModel("Model2026...")

// Delete a single model
err = api.DeleteModel("Model2026...")

// Batch delete
err = api.BatchDeleteModel([]string{"Model2026...", "Model2026..."})

// Remove background
outURL, err := api.RemoveBackground("https://cos.../image.png")

// Text-to-motion prompt suggestions
prompts, err := api.GetText2MotionPromptList("en")
```

---

## Errors

All SDK errors inherit from `WeaverError`; you can catch the base class or any subclass.

| Error | Code | Description |
|---|---|---|
| `WeaverError` | any | Base error |
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

client := visvise.NewClient("...", "...", "...", nil)

modelID, err := client.Gen360("image.png", visvise.NewGen360Options())
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
    client := visvise.NewClient("...", "...", "...", nil)

    // Step 1: Image-to-360
    fmt.Println("Step 1: generating multi-views...")
    opts := visvise.NewGen360Options()
    mvID, _ := client.Gen360("character.png", opts)
    mv, _ := client.WaitModel(mvID, &visvise.WaitOptions{Interval: 3, Timeout: 300})
    views := mv.ImageGen360Output.OutputView

    // Step 2: High-poly model
    fmt.Println("Step 2: generating high-poly model...")
    opts = visvise.NewGenHighModelOptions().
        SetBackView(views.BackView).
        SetLeftView(views.LeftView).
        SetRightView(views.RightView)
    highID, _ := client.GenHighModel(views.MainView, opts)
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
    client := visvise.NewClient("...", "...", "...", nil)

    // Step 1: Rigging
    opts := visvise.NewGenRiggingOptions()
    rigID, _ := client.GenRigging("character.fbx", opts)
    rig, _ := client.WaitModel(rigID, &visvise.WaitOptions{Timeout: 600})
    fmt.Println("Rigged model:", rig.OutputModel)

    // Step 2: Skinning
    opts = visvise.NewGenSkinningOptions([]string{"Body_Mesh"}, []string{"Bip001", "Bip001 Pelvis"})
    skinID, _ := client.GenSkinning("rigged_character.fbx", opts)
    client.WaitModel(skinID, &visvise.WaitOptions{Timeout: 600})

    // Step 3: Video-to-animation
    opts = visvise.NewGenVideoMotionOptions()
    animID, _ := client.GenVideoMotion("skinned_model.fbx", "dance.mp4", opts)
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
    client := visvise.NewClient("...", "...", "...", nil)

    reduceFaces := []visvise.ReduceFace{
        {ReduceLevel: 1, ReducePercent: 50, FaceType: visvise.FaceTypeQuad},
        {ReduceLevel: 2, ReducePercent: 25, FaceType: visvise.FaceTypeQuad},
    }

    opts := visvise.NewGenLODOptions()
    modelIDs, _ := client.GenLOD("high_model.fbx", reduceFaces, opts)

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
