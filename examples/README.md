# Examples

**[English](README_EN.md)** | 中文

本目录包含 VISVISE Go SDK 的完整示例代码，每个子目录对应一个功能节点。

## 目录结构

```
examples/
├── gen_360/           # 图生360 - 从单张图片生成360°多视图
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
├── gen_segment_2d/    # 2D拆分 - 对多视图进行组件分割
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
