package examples

import (
	"fmt"

	"github.com/visvise/visvise-sdk-go/visvise"
)

// ExampleGen360 demonstrates how to generate a 3D model from an image
func ExampleGen360() {
	// Create client
	client := visvise.NewClient(
		"your_app_id",
		"your_secret_key",
		"your_uid",
		visvise.EnvProd,
		30,
	)

	// Generate 3D model from image using options
	// name is set via Options struct (default: "gen_360")
	opts := visvise.NewGen360Options().
		SetName("my_360_model").                // optional, default is "gen_360"
		SetAlgorithmModel("hunyuan3D-MultiView-v3.0") // optional
	// Default output format is fbx, face type is triangle

	modelID, err := client.Gen360("path/to/character.png", opts)
	if err != nil {
		panic(err)
	}

	// Wait for model generation to complete
	model, err := client.WaitModel(modelID, &visvise.WaitOptions{
		Interval: 2.0,
		Timeout:  600,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Model generated successfully!\n")
	fmt.Printf("Output model URL: %s\n", model.OutputModel)
	fmt.Printf("Preview image URL: %s\n", model.PreviewImg)
}

// ExampleGenHighModel demonstrates how to generate a high-detail 3D model
func ExampleGenHighModel() {
	client := visvise.NewClient(
		"your_app_id",
		"your_secret_key",
		"your_uid",
		visvise.EnvProd,
		30,
	)

	opts := visvise.NewGenHighModelOptions().
		SetName("my_high_model").                             // optional, default is "gen_high_model"
		SetAlgorithmModel("hunyuan3D-v2.0").                  // optional
		SetFaceNum(500000).                                    // optional
		SetOutputModelFormat(visvise.OutputModelFormatFBX).    // optional, default is fbx
		SetFaceType(visvise.FaceTypeTriangle)                 // optional, default is triangle

	modelID, err := client.GenHighModel("path/to/character.png", opts)
	if err != nil {
		panic(err)
	}

	model, err := client.WaitModel(modelID, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("High model generated: %s\n", model.OutputModel)
}

// ExampleGenLOD demonstrates how to generate LOD models
func ExampleGenLOD() {
	client := visvise.NewClient(
		"your_app_id",
		"your_secret_key",
		"your_uid",
		visvise.EnvProd,
		30,
	)

	// Configure LOD levels
	reduceFaces := []visvise.ReduceFace{
		{ReduceLevel: 1, ReducePercent: 50, FaceType: visvise.FaceTypeTriangle},
		{ReduceLevel: 2, ReducePercent: 30, FaceType: visvise.FaceTypeTriangle},
		{ReduceLevel: 3, ReducePercent: 20, FaceType: visvise.FaceTypeTriangle},
	}

	opts := visvise.NewGenLODOptions().
		SetName("my_lod_model").                  // optional, default is "gen_lod"
		SetAlgorithmModel("VISVISE-LOD-V1.0.0"). // optional
		SetGenTimes(3).                           // optional, default is 3
		SetOutputModelFormat(visvise.OutputModelFormatFBX) // optional, default is fbx

	modelIDs, err := client.GenLOD("path/to/high_model.fbx", reduceFaces, opts)
	if err != nil {
		panic(err)
	}

	// Wait for all LOD models to complete
	for i, modelID := range modelIDs {
		model, err := client.WaitModel(modelID, nil)
		if err != nil {
			panic(err)
		}
		fmt.Printf("LOD %d generated: %s\n", i+1, model.OutputModel)
	}
}

// ExampleGenRigging demonstrates how to perform skeleton rigging
func ExampleGenRigging() {
	client := visvise.NewClient(
		"your_app_id",
		"your_secret_key",
		"your_uid",
		visvise.EnvProd,
		30,
	)

	opts := visvise.NewGenRiggingOptions().
		SetName("my_rigged_model").             // optional, default is "gen_rigging"
		SetAlgorithmModel("VISVISE-GoRigging-V1.0.0"). // optional
		SetMeshCategory("humanoid")             // optional, default is humanoid

	modelID, err := client.GenRigging("path/to/model.fbx", opts)
	if err != nil {
		panic(err)
	}

	model, err := client.WaitModel(modelID, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Rigged model generated: %s\n", model.OutputModel)
}

// ExampleGenTexture demonstrates how to generate textures
func ExampleGenTexture() {
	client := visvise.NewClient(
		"your_app_id",
		"your_secret_key",
		"your_uid",
		visvise.EnvProd,
		30,
	)

	opts := visvise.NewGenTextureOptions().
		SetName("my_textured_model").           // optional, default is "gen_texture"
		SetAlgorithmModel("VISVISE-Tex-V1.0.0"). // optional
		SetInputView(&visvise.View{
			MainView:  "path/to/front.png",
			BackView:  "path/to/back.png",
			LeftView:  "path/to/left.png",
			RightView: "path/to/right.png",
		}). // optional
		SetResolution(2048).   // optional
		SetPrompt("high quality, realistic") // optional, at least one of input_view.main_view or prompt is required

	modelID, err := client.GenTexture("path/to/model.fbx", opts)
	if err != nil {
		panic(err)
	}

	model, err := client.WaitModel(modelID, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Textured model generated: %s\n", model.OutputModel)
}

// ExampleGenSegment2D demonstrates how to perform 2D segmentation
func ExampleGenSegment2D() {
	client := visvise.NewClient(
		"your_app_id",
		"your_secret_key",
		"your_uid",
		visvise.EnvProd,
		30,
	)

	opts := visvise.NewGenSegment2DOptions().
		SetName("my_segment").                           // optional, default is "gen_segment_2d"
		SetAlgorithmModel("VISVISE-Segment-V1.0.0").    // optional
		SetSplitType(visvise.SegmentSplitFrontView).    // optional
		SetGranularity(visvise.SegmentGranularityMedium). // optional
		SetPrompt("segment by body parts").             // optional
		SetOnThinking(func(thought string) {
			fmt.Println("Thinking:", thought)
		}) // optional

	modelID, err := client.GenSegment2D("model_id_from_gen_360", opts)
	if err != nil {
		panic(err)
	}

	fmt.Printf("2D segmentation completed. Segment model ID: %s\n", modelID)
}

// ExampleLowLevelAPI demonstrates how to use the low-level API
func ExampleLowLevelAPI() {
	client := visvise.NewClient(
		"your_app_id",
		"your_secret_key",
		"your_uid",
		visvise.EnvProd,
		30,
	)

	// Get user quota
	quota, err := client.GetAPI().GetUserQuota()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Remaining quota: %d\n", quota.Quota)

	// List available algorithm models
	models, err := client.GetAPI().ListAlgorithmModel(int(visvise.NodeTypeImgTo360), nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Available models for IMG_TO_360: %v\n", models)

	// Get model list
	modelList, total, err := client.GetAPI().GetModelList(nil, nil, nil, "", 20, 1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Total models: %d\n", total)
	for _, m := range modelList {
		fmt.Printf("  - %s (status=%d)\n", m.Name, m.Status)
	}

	// Download model
	downloadURL, err := client.GetAPI().DownloadModel("model_id")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Download URL: %s\n", downloadURL)
}
