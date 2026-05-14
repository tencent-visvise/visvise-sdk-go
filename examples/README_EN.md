# Examples

**[中文](README.md)** | English

This directory contains complete example code for the VISVISE Go SDK, with each subdirectory corresponding to a specific functionality node.

## Directory Structure

```
examples/
├── gen_360/           # Image to 360 - Generate 360° multi-view from a single image
├── gen_high_model/    # Image to High-poly Model - Generate high-polygon 3D models
├── gen_mid_model/     # Image to Mid-poly Model - Generate medium-polygon 3D models
├── gen_low_model/     # Image to Low-poly Model - Generate low-polygon 3D models
├── gen_retopology/    # Retopology - Optimize model topology structure
├── gen_lod/           # LOD Reduction - Generate multi-level-of-detail models
├── gen_mesh_refine/   # Mesh Refine - Optimize model mesh topology
├── gen_uv/            # UV Unwrapping - Generate UV mapping coordinates for models
├── gen_texture/       # Texture Mapping - Generate texture maps for models
├── gen_rigging/       # Rigging - Generate skeletal structure for models
├── gen_skinning/      # Skinning - Generate skin weights for models with skeleton
├── gen_video_motion/ # Video to Animation - Extract motion from video to drive models
├── gen_text_motion/  # Text to Animation - Generate animations from text descriptions
├── gen_pose/          # Batch Image to Pose - Extract poses from images to drive models
├── gen_segment_2d/    # 2D Segmentation - Segment components from multi-view images
└── assets/           # Example resource files directory
```

## Environment Configuration

Before running the examples, you need to set the following environment variables:

```bash
export VISVISE_APP_ID="your_app_id"
export VISVISE_SECRET_KEY="your_secret_key"
export VISVISE_UID="your_uid"
export VISVISE_ENV="prod"  # Optional: prod, test, dev
```

## Running Examples

Each example is an independent executable program. Enter the corresponding directory and run:

```bash
# Enter directory
cd examples/gen_360

# Run
go run main.go
```

## Example Descriptions

### gen_360 (Image to 360)
Generate 360° multi-view from a single image. The output can be used as input for subsequent image-to-high/mid/low model generation.

### gen_high_model (Image to High-poly Model)
Uses Tripo-v3.1-ultra algorithm, supports single image input (no need for four views).

### gen_mid_model (Image to Mid-poly Model)
Requires all four views. You can set the `MV_360_MODEL_ID` environment variable to automatically extract four views from gen_360 output.

### gen_low_model (Image to Low-poly Model)
Only requires main view, uses Tripo's fast generation algorithm.

### gen_retopology (Retopology)
- Hunyuan model: pass `DetailLevel`
- VISVISE proprietary model: pass `FaceNum`

### gen_lod (LOD Reduction)
Generate models with multiple levels of detail, supports quad and triangle faces.

### gen_texture (Texture Mapping)
At least one of `InputView.MainView` or `Prompt` must be provided.

### gen_rigging (Rigging)
- `MeshCategory`: `humanoid` or `tetrapod`

### gen_skinning (Skinning)
Requires providing `MeshNames` and `JointNames` lists.

### gen_video_motion (Video to Animation)
Extract motion data from video to drive 3D models.

### gen_text_motion (Text to Animation)
Generate animations from text descriptions, returns multiple versions for selection.

### gen_pose (Batch Image to Pose)
Extract poses from reference images.

### gen_segment_2d (2D Segmentation)
Requires running `gen_360` first to get `MV_360_MODEL_ID`.
The segmentation result can be used as the `SegmentModelID` parameter for `gen_mid_model` or `gen_low_model`.
