# Examples

**[English](README_EN.md)** | 中文

本目录包含 VISVISE Go SDK 的完整示例代码，每个子目录对应一个功能节点。

## 目录结构

```
examples/
├── gen_360/           # 图生360 - 从单张图片生成360°多视图
├── gen_preprocess/    # 2D预处理 - 风格化或去花纹并同步保存资产
├── gen_high_model/    # 图生高模 - 生成高面数3D模型
├── gen_mid_model/     # 图生中模 - 生成中等面数3D模型
├── gen_low_model/     # 图生低模 - 生成低面数3D模型
├── gen_retopology/    # 重拓扑 - 优化模型拓扑结构
├── gen_lod/           # LOD减面 - 生成多级细节模型
├── gen_mesh_refine/   # 重布线 - 优化模型网格布线
├── gen_uv/            # UV展开 - 生成模型的UV贴图坐标
├── gen_texture/       # 贴图纹理 - 为模型生成纹理贴图
├── gen_rigging/       # 骨骼架设 - 为模型生成骨骼结构
├── gen_skinning/      # 蒙皮生成 - 为带骨骼的模型生成蒙皮权重
├── gen_video_motion/  # 视频生动画 - 从视频提取动作驱动模型
├── gen_text_motion/   # 文本生动画 - 通过文字描述生成动画
├── gen_pose/          # 批量图生Pose - 从图片提取姿态驱动模型
├── gen_segment_2d/           # 2D拆分 - 对图生360的多视图进行组件分割
├── gen_segment_2d_skip360/   # 2D拆分（跳过360）- 上传本地六视图原画，跳过360前置与去背景
├── gen_segment_2d_zip/       # 2D拆分 - 获取拆分产物 zip_path 的两种方式（详情 / SSE reply）
└── assets/            # 示例资源文件目录
```

## 环境配置

运行示例前，需要设置以下环境变量：

```bash
export VISVISE_APP_ID="your_app_id"
export VISVISE_SECRET_KEY="your_secret_key"
export VISVISE_RTX="your_rtx"
export VISVISE_ENV="prod"  # 可选: prod, test, dev
```

## 运行示例

每个示例都是独立的可执行程序，进入对应目录运行：

```bash
# 进入目录
cd examples/gen_360

# 运行
go run main.go
```

## 示例说明

### gen_360 (图生360)
从单张图片生成 360° 多视图，输出可作为后续图生高模/中模/低模的输入。

### gen_preprocess (2D预处理)
对图片进行风格化或去花纹处理，并保存为 2D 预处理资产。

### gen_high_model (图生高模)
使用 Tripo-v3.1-ultra 算法，支持单图输入（无需四视图）。

### gen_mid_model (图生中模)
需要四视图全部传入。可以通过设置 `MV_360_MODEL_ID` 环境变量从 gen_360 输出自动提取四视图。

### gen_low_model (图生低模)
仅需主视图，使用 Tripo 快速生成算法。

### gen_retopology (重拓扑)
- 混元模型：传 `DetailLevel`
- VISVISE 自研模型：传 `FaceNum`

### gen_lod (LOD减面)
生成多个细节级别的模型，支持 quad 面和 triangle 面。

### gen_texture (贴图纹理)
`InputView.MainView` 和 `Prompt` 至少需要传一个。

### gen_rigging (骨骼架设)
- `MeshCategory`: `visvise.MeshCategoryHumanoid`(人形，默认) 或 `visvise.MeshCategoryTetrapod`(四足)

### gen_skinning (蒙皮生成)
需要提供 `MeshNames` 和 `JointNames` 列表。

### gen_video_motion (视频生动画)
从视频提取动作数据驱动 3D 模型。

### gen_text_motion (文本生动画)
通过文字描述生成动画，一次返回多个版本供抽卡。

### gen_pose (批量图生Pose)
从参考图片提取姿态。

### gen_segment_2d (2D拆分)
需要先运行 `gen_360` 获取 `MV_360_MODEL_ID`。
分割结果可作为 `gen_mid_model` 的 `SegmentModelID` 参数。

### gen_segment_2d_skip360 (2D拆分 - 跳过360 + 跳过去背景 + 上传六视图)
不需要 `MV_360_MODEL_ID`，直接上传本地原画（`SEG_ASSETS_DIR`，默认 `./tests/assets`），
并跳过 360 前置与去背景预处理，原始图原样送拆分算法。
`split_type=2` 支持 1~6 个视图，素材缺失的视图会自动跳过。

### gen_segment_2d_zip (2D拆分 - 获取拆分产物 zip_path)
不需要 `MV_360_MODEL_ID`，直接上传本地六视图原画（`SEG_ASSETS_DIR`，默认 `./tests/assets`），
跳过 360 前置与去背景预处理；方式二还需账号下有可用的 2D 拆分算法模型。
演示两种取 zip_path 的方式：
1. `GenSegment2D` + `WaitModel` → `ModelInfo.SegmentOutput.ZipPath`（推荐，可同时拿到各视图部件表 / npz / 彩色预览图）
2. 底层 `GetAPI().InitSegment` 自己遍历 SSE 帧，从 `reply` 帧直接取 `zip_path`（生成瞬间拿到，不用再查一次）

两种方式拿到的是同一个带签名的 COS URL，24h 内直接 GET 即可下载（无需鉴权头）。
