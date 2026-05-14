# VISVISE Weaver Go SDK

**[English](README_EN.md)** | 中文

VISVISE Weaver OpenAPI 的 Go SDK，提供：

- 全部原子 API 方法（逐一对应 OpenAPI 接口）
- 各节点类型的高阶 `Gen*()` 方法（自动上传文件 + 创建任务）
- `WaitModel()` 异步轮询方法

---

## 目录

- [安装](#安装)
- [快速开始](#快速开始)
- [客户端初始化](#客户端初始化)
- [枚举常量](#枚举常量)
- [高阶方法参考](#高阶方法参考)
  - [Gen360 — 图生360](#gen360--图生360)
  - [GenHighModel — 图生高模](#genhighmodel--图生高模)
  - [GenMidModel — 图生中模](#genmidmodel--图生中模)
  - [GenLowModel — 图生低模](#genlowmodel--图生低模)
  - [GenMeshRefine — 重布线](#genmeshrefine--重布线)
  - [GenRetopology — 重拓扑](#genretopology--重拓扑)
  - [GenLOD — LOD](#genlod--lod)
  - [GenUV — UV展开](#genuv--uv展开)
  - [GenTexture — 贴图纹理](#gentexture--贴图纹理)
  - [GenRigging — 骨骼架设](#genrigging--骨骼架设)
  - [GenSkinning — 蒙皮生成](#genskinning--蒙皮生成)
  - [GenVideoMotion — 视频生动画](#genvideomotion--视频生动画)
  - [GenTextMotion — 文本生动画](#gentextmotion--文本生动画)
  - [GenPose — 图生Pose](#genpose--图生pose)
  - [GenSegment2D — 2D 拆分](#gensegment2d--2d-拆分)
  - [WaitModel — 等待完成](#waitmodel--等待完成)
- [原子 API 方法参考](#原子-api-方法参考)
- [错误说明](#错误说明)
- [完整流程示例](#完整流程示例)

---

## 安装

使用 `go get` 安装：

```bash
go get github.com/tencent-visvise/visvise-sdk-go
```

---

## 快速开始

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

    // ① 图生360：上传本地图片，生成多视图
    // name 在 Options 中设置（可选，默认自动生成）
    opts := visvise.NewGen360Options().
        SetEnableAPose(true)  // 可选
    mvModelID, err := client.Gen360("character.png", opts)
    if err != nil {
        panic(err)
    }

    // ② 等待图生360完成，获取多视图输出
    mvInfo, err := client.WaitModel(mvModelID, &visvise.WaitOptions{Interval: 3, Timeout: 300})
    if err != nil {
        panic(err)
    }
    outputView := mvInfo.ImageGen360Output.OutputView

    // ③ 图生高模（多视图输出的 COS URL 直接传入）
    opts = visvise.NewGenHighModelOptions().
        SetBackView(outputView.BackView, "").
        SetLeftView(outputView.LeftView, "").
        SetRightView(outputView.RightView, "")
    highModelID, err := client.GenHighModel(outputView.MainView, opts)
    if err != nil {
        panic(err)
    }

    // ④ 等待高模完成
    modelInfo, err := client.WaitModel(highModelID, &visvise.WaitOptions{Timeout: 900})
    if err != nil {
        panic(err)
    }
    fmt.Println("输出模型:", modelInfo.OutputModel)
}
```

---

## 客户端初始化

```go
import "github.com/visvise/visvise-sdk-go/visvise"

client := visvise.NewClient(
    "your_app_id",       // 必填，由平台分配
    "your_secret_key",   // 必填，由平台分配
    "your_uid",          // 必填，从申请 key 的登录账号获取
    visvise.EnvProd,     // 可选，默认生产环境
    30,                  // 可选，单次请求超时（秒），默认 30
)
```

| 参数 | 必填 | 说明 |
|---|---|---|
| `appID` | ✅ | 由平台分配的客户端标识 |
| `secretKey` | ✅ | 由平台分配的签名密钥 |
| `uid` | ✅ | 用户 ID，从申请 key 的登录账号获取 |
| `env` | — | 环境：`EnvProd`（默认）/ `EnvTest` / `EnvDev` |
| `timeout` | — | 单次 HTTP 请求超时（秒），默认 30 |

---

## 枚举常量

SDK 提供以下常量，推荐使用常量替代硬编码数字/字符串：

```go
import "github.com/visvise/visvise-sdk-go/visvise"

// 面数类型
visvise.FaceTypeTriangle // 1 - 三角面
visvise.FaceTypeQuad     // 2 - 四边面

// 精细程度（重拓扑）
visvise.DetailLevelLow    // 1 - 低
visvise.DetailLevelMedium // 2 - 中
visvise.DetailLevelHigh   // 3 - 高

// 输出模型格式
visvise.OutputModelFormatFBX // "fbx"
visvise.OutputModelFormatOBJ // "obj"
visvise.OutputModelFormatGLB // "glb"

// 布线优化模式
visvise.MeshRefineModeOptimize // 1 - 布线优化
visvise.MeshRefineModeDensify  // 2 - 布线加密

// 2D 拆分方式
visvise.SegmentSplitFrontView // 1 - 生成正视图拆分（默认）
visvise.SegmentSplitFourView  // 2 - 生成四视图拆分

// 2D 拆分颗粒度
visvise.SegmentGranularityCoarse  // 1 - 粗
visvise.SegmentGranularityMedium  // 2 - 中（默认）
visvise.SegmentGranularityFine    // 3 - 细
```

---

## 高阶方法参考

高阶方法封装了「COS 文件上传 + 创建异步任务」两步，传入文件路径（本地）或 COS URL 均可，返回 `model_id`。

所有 `Gen*()` 方法采用 **Options 结构体** 模式，通过链式调用设置可选参数，使用更简洁：

> **关于 `name` 参数：** 所有 Gen* 方法的 `name` 参数均为可选，默认值在 `New*Options()` 中设置（如 `gen_360`、`gen_high_model` 等）。可通过 `SetName()` 自定义。
>
> **关于 `algorithmModel` 参数：** 所有 Gen* 方法的 `algorithmModel` 参数均为可选。若不传，SDK 将自动调用 `ListAlgorithmModel` 获取当前账号可用的第一个算法模型。

### Gen360 — 图生360

从单张图片生成 360 度多视图。

```go
// 可选参数通过 Options 结构体传递，支持链式调用
// 默认 name="gen_360"，输出格式为 fbx，面数类型为三角面
opts := visvise.NewGen360Options().
    SetName("my_360").                                    // 可选，默认 "gen_360"
    SetAlgorithmModel("hunyuan3D-MultiView-v3.0").        // 可选
    SetOutputModelFormat(visvise.OutputModelFormatFBX).   // 可选，输出格式（默认 fbx）
    SetFaceType(visvise.FaceTypeTriangle).                // 可选，面数类型（默认三角面）
    SetEnableAPose(true).                                 // 可选，是否开启 A-Pose
    SetStyle("anime").                                    // 可选，风格类型
    SetBackView("path/to/back.png", "back.png").          // 可选，背视图
    SetLeftView("path/to/left.png", "left.png").          // 可选，左视图
    SetRightView("path/to/right.png", "right.png")        // 可选，右视图

modelID, err := client.Gen360("path/to/character.png", opts)
```

---

### GenHighModel — 图生高模

从图片/多视图生成高精度 3D 模型（node_type=3）。

```go
opts := visvise.NewGenHighModelOptions().
    SetName("my_high_model").                            // 可选，默认 "gen_high_model"
    SetAlgorithmModel("hunyuan3D-v3.1").                 // 可选
    SetOutputModelFormat(visvise.OutputModelFormatFBX).  // 可选，输出格式（默认 fbx）
    SetFaceType(visvise.FaceTypeTriangle).               // 可选，面数类型（默认三角面）
    SetFaceNum(500000).                                  // 可选，目标面数（1000~1500000）
    SetBackView(outputView.BackView, "").                 // 可选，背视图
    SetLeftView(outputView.LeftView, "").                 // 可选，左视图
    SetRightView(outputView.RightView, "")                 // 可选，右视图

modelID, err := client.GenHighModel("path/to/main.png", opts)
```

---

### GenMidModel — 图生中模

中模要求四视图全部必传（node_type=11）。

```go
opts := visvise.NewGenMidModelOptions().
    SetName("my_mid_model").                             // 可选，默认 "gen_mid_model"
    SetAlgorithmModel("VISVISE-MeshGen-V1.0.0").         // 可选
    SetOutputModelFormat(visvise.OutputModelFormatFBX).  // 可选，输出格式
    SetFaceType(visvise.FaceTypeTriangle).              // 可选，面数类型
    SetSegmentModelID("Model2026...")                    // 可选，2D 分割资产 ID

// mainView, backView, leftView, rightView 四个视图必填
modelID, err := client.GenMidModel(
    "path/to/main.png",
    "path/to/back.png",
    "path/to/left.png",
    "path/to/right.png",
    opts,
)
```

---

### GenLowModel — 图生低模

低模只需主视图（node_type=13）。

```go
opts := visvise.NewGenLowModelOptions().
    SetName("my_low_model").                             // 可选，默认 "gen_low_model"
    SetAlgorithmModel("Tripo-v1.0-快速生成").             // 可选
    SetOutputModelFormat(visvise.OutputModelFormatFBX).  // 可选，输出格式
    SetFaceType(visvise.FaceTypeTriangle).               // 可选，面数类型
    SetBackView("path/to/back.png", "back.png").          // 可选，背视图
    SetLeftView("path/to/left.png", "left.png").          // 可选，左视图
    SetRightView("path/to/right.png", "right.png")        // 可选，右视图

modelID, err := client.GenLowModel("path/to/main.png", opts)
```

---

### GenMeshRefine — 重布线

对模型进行布线优化（node_type=10）。

```go
opts := visvise.NewGenMeshRefineOptions().
    SetName("my_mesh_refine").                           // 可选，默认 "gen_mesh_refine"
    SetAlgorithmModel("VISVISE-MeshRefine-V1.0.0").      // 可选
    SetInputModelFormat(visvise.OutputModelFormatFBX).   // 可选，输入模型格式（默认 fbx）
    SetMode(visvise.MeshRefineModeOptimize).             // 可选，布线优化模式
    SetColorModel("path/to/color.fbx", "color.fbx")      // 可选，色彩模型

modelID, err := client.GenMeshRefine("path/to/model.fbx", opts)
```

---

### GenRetopology — 重拓扑

对高面数模型进行拓扑优化（node_type=1）。

> 注意：混元模型传 `DetailLevel`，VISVISE 自研模型传 `FaceNum`，二选一。

```go
// 混元模型示例
opts := visvise.NewGenRetopologyOptions().
    SetName("my_retopo").                                // 可选，默认 "gen_retopology"
    SetAlgorithmModel("hunyuan3D-RTP-v1.5").             // 可选
    SetOutputModelFormat(visvise.OutputModelFormatFBX).   // 可选，输出格式
    SetFaceType(visvise.FaceTypeQuad).                   // 可选，面数类型（默认四边面）
    SetDetailLevel(visvise.DetailLevelMedium)            // 可选，混元模型

// VISVISE 自研模型示例
opts := visvise.NewGenRetopologyOptions().
    SetAlgorithmModel("VISVISE-RTP-V1.0.0").
    SetFaceNum(50000)                                    // 可选，VISVISE 自研模型

modelID, err := client.GenRetopology("path/to/model.fbx", opts)
```

---

### GenLOD — LOD

生成多级细节模型（node_type=2），支持抽卡。默认抽卡次数为 3。

```go
reduceFaces := []visvise.ReduceFace{
    {ReduceLevel: 1, ReducePercent: 50, FaceType: int(visvise.FaceTypeQuad)},
    {ReduceLevel: 2, ReducePercent: 25, FaceType: int(visvise.FaceTypeQuad)},
}

opts := visvise.NewGenLODOptions().
    SetName("my_lod").                                   // 可选，默认 "gen_lod"
    SetAlgorithmModel("VISVISE-LOD-V1.0.0").              // 可选
    SetOutputModelFormat(visvise.OutputModelFormatFBX).   // 可选，输出格式（默认 fbx）
    SetGenTimes(3)                                       // 可选，抽卡次数（默认 3）

modelIDs, err := client.GenLOD("path/to/model.fbx", reduceFaces, opts)
```

---

### GenUV — UV展开

自动 UV 展开（node_type=9）。

```go
opts := visvise.NewGenUVOptions().
    SetName("my_uv").                                     // 可选，默认 "gen_uv"
    SetAlgorithmModel("hunyuan3D-UV-v2.0").              // 可选
    SetEnableAutoSmoothing(true)                          // 可选，是否启用自动平滑

modelID, err := client.GenUV("path/to/model.fbx", opts)
```

---

### GenTexture — 贴图纹理

为模型生成贴图纹理（node_type=8）。

> `InputView.MainView` 与 `Prompt` 至少传一个，可同时传入。

```go
view := &visvise.View{MainView: "path/to/ref.png"}

opts := visvise.NewGenTextureOptions().
    SetName("my_texture").                               // 可选，默认 "gen_texture"
    SetAlgorithmModel("hunyuan3D-TEX-v2.0").              // 可选
    SetInputView(view).                                   // 可选，原画视图（与 prompt 至少传一个）
    SetResolution(2048).                                  // 可选，分辨率
    SetUnwarpUV(true).                                    // 可选，是否同时展开 UV
    SetPrompt("high quality, realistic")                   // 可选，贴图文本提示词

modelID, err := client.GenTexture("path/to/model.fbx", opts)
```

---

### GenRigging — 骨骼架设

自动为模型生成骨骼（node_type=5）。SDK 自动将模型文件与参数 JSON 打包成 zip 上传，无需手动准备 zip 包。

```go
opts := visvise.NewGenRiggingOptions().
    SetName("my_rigging").                                // 可选，默认 "gen_rigging"
    SetAlgorithmModel("VISVISE-GoRigging-V1.0.0").       // 可选
    SetMeshCategory("humanoid").                          // 可选，"humanoid"（人形，默认）或 "tetrapod"（四足）
    SetTemplateSkeleton("path/to/skeleton.fbx", "skeleton.fbx") // 可选，模板骨骼

modelID, err := client.GenRigging("path/to/model.fbx", opts)
```

---

### GenSkinning — 蒙皮生成

自动绑定蒙皮权重（node_type=6）。SDK 自动将模型文件与参数 JSON 打包成 zip 上传。

```go
meshNames := []string{"Body_Mesh", "Hair_Mesh"}
jointNames := []string{"Bip001", "Bip001 Pelvis"}

opts := visvise.NewGenSkinningOptions(meshNames, jointNames).
    SetName("my_skinning").                               // 可选，默认 "gen_skinning"
    SetAlgorithmModel("VISVISE-GoSkinning-V1.0.0")        // 可选

modelID, err := client.GenSkinning("path/to/rigged_model.fbx", opts)
```

---

### GenVideoMotion — 视频生动画

从视频中提取动作驱动 3D 模型（node_type=4）。

```go
opts := visvise.NewGenVideoMotionOptions().
    SetName("my_video_motion").                           // 可选，默认 "gen_video_motion"
    SetAlgorithmModel("VISVISE-FramingAI-Base-V1.5.0").  // 可选
    SetOutputModelFormat(visvise.OutputModelFormatFBX).   // 可选，输出格式
    SetWithHand(true).                                    // 可选，是否开启手部捕捉
    SetMultipleTrack(false).                              // 可选，是否开启多人捕捉
    SetRotateAxisAngle(0, 0, 0)                           // 可选，旋转轴角 [x, y, z]（弧度）

modelID, err := client.GenVideoMotion("path/to/model.fbx", "path/to/dance.mp4", opts)
```

---

### GenTextMotion — 文本生动画

通过提示词生成动画，一次返回 4 个模型供抽卡（node_type=4）。

```go
opts := visvise.NewGenTextMotionOptions().
    SetName("my_text_motion").                            // 可选，默认 "gen_text_motion"
    SetAlgorithmModel("VISVISE-TextMotion-V1.1.0").       // 可选
    SetOutputModelFormat(visvise.OutputModelFormatFBX).   // 可选，输出格式

modelIDs, err := client.GenTextMotion("path/to/model.fbx", "一个人在跳街舞", opts)
// modelIDs 包含 4 个 ID，等待其中你需要的那个即可
```

---

### GenPose — 图生Pose

从参考图生成 Pose 模型（最多 10 张图片）。

```go
inputImages := []visvise.FileInput{
    "path/to/pose_ref_1.png",
    "path/to/pose_ref_2.png",
}

opts := visvise.NewGenPoseOptions().
    SetName("my_pose").                                   // 可选，默认 "gen_pose"
    SetAlgorithmModel("VISVISE-PosingAI-V1.0.0").         // 可选
    SetOutputModelFormat(visvise.OutputModelFormatFBX).  // 可选，输出格式
    SetImageFilenames([]string{"pose_ref_1.png", "pose_ref_2.png"}) // 可选

modelIDs, err := client.GenPose("path/to/model.fbx", inputImages, opts)
```

---

### GenSegment2D — 2D 拆分

对图生 360 输出的多视图进行组件分割（node_type=14，SSE 协议）。生成的分割资产 `model_id` 可作为图生中模/低模的 `segmentModelID` 输入。

```go
onThinking := func(content string) {
    fmt.Println("[思考]", content)
}

opts := visvise.NewGenSegment2DOptions().
    SetName("my_segment").                                // 可选，默认 "gen_segment_2d"
    SetAlgorithmModel("VISVISE-Seg2D-V1.0.0").           // 可选
    SetSplitType(int(visvise.SegmentSplitFrontView)).    // 可选，正视图/四视图拆分
    SetGranularity(int(visvise.SegmentGranularityMedium)).// 可选，颗粒度
    SetPrompt("segment by body parts").                   // 可选，自然语言描述拆分规则
    SetOnThinking(onThinking)                             // 可选，处理 thinking 事件的回调

segModelID, err := client.GenSegment2D("Model2026...", opts)
// 后续可作为 segmentModelID 传给 GenMidModel
```

---

### WaitModel — 等待完成

轮询等待异步任务完成，返回 `ModelInfo`。

```go
modelInfo, err := client.WaitModel(
    "Model2026033100192028",
    &visvise.WaitOptions{
        Interval: 2,  // 轮询间隔（秒），默认 2
        Timeout: 600, // 超时时长（秒），默认 600
    },
)

fmt.Println(modelInfo.OutputModel) // 输出模型下载 URL
fmt.Println(modelInfo.TimeCost)   // 耗时（秒）
```

**错误：**

- `PollingTimeoutError`：超时仍未完成时抛出
- `ModelGenerationError`：模型生成失败（status=4）时抛出
- `InvalidParamsError`：轮询接口返回参数错误时立即抛出（不重试）
- 其他网络/业务错误不抛出，会打印日志并继续重试

---

## 原子 API 方法参考

通过 `client.GetAPI().xxx()` 访问底层接口：

```go
api := client.GetAPI()

// 获取临时上传凭证
cred, err := api.GetCosCred()

// 查询剩余配额
quota, err := api.GetUserQuota()
fmt.Println(quota.Quota) // 剩余次数

// 拉取模型列表
models, total, err := api.GetModelList(
    []string{"Model2026..."},
    nil, nil, "", 10, 1,
)

// 获取算法模型列表
algModels, err := api.ListAlgorithmModel(int(visvise.NodeTypeAnimation), nil)

// 获取下载链接
url, err := api.DownloadModel("Model2026...")

// 删除单个
err = api.DeleteModel("Model2026...")

// 批量删除
err = api.BatchDeleteModel([]string{"Model2026...", "Model2026..."})

// 去除背景
outURL, err := api.RemoveBackground("https://cos.../image.png")

// 文生动画提示词列表
prompts, err := api.GetText2MotionPromptList("zh")
```

---

## 错误说明

所有 SDK 错误均继承自 `WeaverError`，可以捕获基类也可以精确捕获子类。

| 错误类 | 对应错误码 | 说明 |
|---|---|---|
| `WeaverError` | 任意 | 基础错误 |
| `NetworkError` | — | 网络连接失败、超时等 |
| `SignatureError` | 400 | 签名错误 |
| `InvalidParamsError` | 120008 | 请求参数错误 |
| `UserNotFoundError` | 120017 | 用户未找到 |
| `PermissionDeniedError` | 120018 | 用户无权限 |
| `QuotaExceededError` | 120020 | 每日配额超出上限 |
| `ProjectPermissionError` | 120027 | 项目权限未授权 |
| `ServerNetworkError` | 120028 | 服务器网络错误 |
| `ServerTimeoutError` | 120032 | 服务器处理超时 |
| `RateLimitError` | 120040 | 请求过于频繁 |
| `ModelGenerationError` | — | 模型生成失败（status=4） |
| `PollingTimeoutError` | — | WaitModel 等待超时 |

```go
import (
    "fmt"
    "github.com/visvise/visvise-sdk-go/visvise"
)

client := visvise.NewClient("...", "...", "...", visvise.EnvProd, 30)

modelID, err := client.Gen360("image.png", visvise.NewGen360Options())
if err != nil {
    if _, ok := err.(*visvise.QuotaExceededError); ok {
        fmt.Println("今日配额已用完，明天再试")
    } else if _, ok := err.(*visvise.PollingTimeoutError); ok {
        fmt.Println("等待超时")
    } else {
        fmt.Printf("接口错误: %v\n", err)
    }
}
```

---

## 完整流程示例

### 示例一：图片 → 高模（图生360 + 图生高模）

```go
package main

import (
    "fmt"
    "github.com/visvise/visvise-sdk-go/visvise"
)

func main() {
    client := visvise.NewClient("...", "...", "...", visvise.EnvProd, 30)

    // Step 1: 图生360
    fmt.Println("Step 1: 生成多视图...")
    opts := visvise.NewGen360Options()
    mvID, _ := client.Gen360("character.png", opts)
    mv, _ := client.WaitModel(mvID, &visvise.WaitOptions{Interval: 3, Timeout: 300})
    views := mv.ImageGen360Output.OutputView

    // Step 2: 图生高模
    fmt.Println("Step 2: 生成高模...")
    opts = visvise.NewGenHighModelOptions().
        SetBackView(views.BackView, "").
        SetLeftView(views.LeftView, "").
        SetRightView(views.RightView, "")
    highID, _ := client.GenHighModel(views.MainView, opts)
    highModel, _ := client.WaitModel(highID, &visvise.WaitOptions{Timeout: 900})
    fmt.Println("高模下载地址:", highModel.OutputModel)
}
```

---

### 示例二：动画生成流水线（骨骼 → 蒙皮 → 动画）

```go
package main

import (
    "fmt"
    "github.com/visvise/visvise-sdk-go/visvise"
)

func main() {
    client := visvise.NewClient("...", "...", "...", visvise.EnvProd, 30)

    // Step 1: 骨骼架设
    opts := visvise.NewGenRiggingOptions()
    rigID, _ := client.GenRigging("character.fbx", opts)
    rig, _ := client.WaitModel(rigID, &visvise.WaitOptions{Timeout: 600})
    fmt.Println("骨骼模型:", rig.OutputModel)

    // Step 2: 蒙皮生成
    opts = visvise.NewGenSkinningOptions([]string{"Body_Mesh"}, []string{"Bip001", "Bip001 Pelvis"})
    skinID, _ := client.GenSkinning("rigged_character.fbx", opts)
    client.WaitModel(skinID, &visvise.WaitOptions{Timeout: 600})

    // Step 3: 视频生动画
    opts = visvise.NewGenVideoMotionOptions()
    animID, _ := client.GenVideoMotion("skinned_model.fbx", "dance.mp4", opts)
    anim, _ := client.WaitModel(animID, &visvise.WaitOptions{Timeout: 900})
    fmt.Println("动画下载地址:", anim.OutputModel)
}
```

---

### 示例三：LOD 生成（含抽卡）

```go
package main

import (
    "fmt"
    "github.com/visvise/visvise-sdk-go/visvise"
)

func main() {
    client := visvise.NewClient("...", "...", "...", visvise.EnvProd, 30)

    reduceFaces := []visvise.ReduceFace{
        {ReduceLevel: 1, ReducePercent: 50, FaceType: int(visvise.FaceTypeQuad)},
        {ReduceLevel: 2, ReducePercent: 25, FaceType: int(visvise.FaceTypeQuad)},
    }

    opts := visvise.NewGenLODOptions()
    modelIDs, _ := client.GenLOD("high_model.fbx", reduceFaces, opts)

    // 等待全部完成
    for _, mid := range modelIDs {
        r, _ := client.WaitModel(mid, &visvise.WaitOptions{Timeout: 300})
        fmt.Println(r.ModelID, r.OutputModel)
    }
}
```

---

## 许可证

MIT License
