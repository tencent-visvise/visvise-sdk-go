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

	// Generate 3D model from image
	modelID, err := client.Gen360(
		"path/to/character.png",
		"hunyuan3D-MultiView-v3.0", // algorithm model
		"my_360_model",             // name
		"",                         // mainViewFilename
		nil,                        // enableAPose
		"",                         // style
		"",                         // backView (empty string for nil)
		"",                         // backViewFilename
		"",                         // leftView (empty string for nil)
		"",                         // leftViewFilename
		"",                         // rightView (empty string for nil)
		"",                         // rightViewFilename
	)
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

	// Enable A-Pose
	enableAPose := true

	modelID, err := client.GenHighModel(
		"path/to/character.png",
		"hunyuan3D-v2.0",            // algorithm model
		string(visvise.OutputModelFormatFBX),
		1,                          // faceType: Triangle
		"my_high_model",
		"character.png",            // mainViewFilename
		nil,                        // faceNum
		"",                         // backView
		"",                         // backViewFilename
		"",                         // leftView
		"",                         // leftViewFilename
		"",                         // rightView
		"",                         // rightViewFilename
	)
	if err != nil {
		panic(err)
	}

	model, err := client.WaitModel(modelID, nil)
	if err != nil {
		panic(err)
	}

	_ = enableAPose // Use the variable
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
		{ReduceLevel: 1, ReducePercent: 50, FaceType: int(visvise.FaceTypeTriangle)},
		{ReduceLevel: 2, ReducePercent: 30, FaceType: int(visvise.FaceTypeTriangle)},
		{ReduceLevel: 3, ReducePercent: 20, FaceType: int(visvise.FaceTypeTriangle)},
	}

	modelIDs, err := client.GenLOD(
		"path/to/high_model.fbx",
		reduceFaces,
		"VISVISE-LOD-V1.0.0", // algorithm model
		string(visvise.OutputModelFormatFBX),
		"my_lod_model",
		"high_model.fbx",     // filename
		3,                    // genTimes: generate 3 results for comparison
	)
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

	modelID, err := client.GenRigging(
		"path/to/model.fbx",
		"VISVISE-GoRigging-V1.0.0", // algorithm model
		"humanoid",                  // meshCategory
		"my_rigged_model",
		"model.fbx", // filename
		nil,         // templateSkeleton
		"",          // templateSkeletonFilename
	)
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

	// Resolution pointer
	resolution := 2048

	modelID, err := client.GenTexture(
		"path/to/model.fbx",
		"VISVISE-Tex-V1.0.0", // algorithm model
		"my_textured_model",
		"model.fbx", // filename
		&visvise.View{
			MainView: "path/to/front.png",
			BackView: "path/to/back.png",
			LeftView: "path/to/left.png",
			RightView: "path/to/right.png",
		},
		&resolution,
		nil,                           // unwarpUV
		"high quality, realistic",      // prompt
	)
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

	// Split type and granularity
	splitType := 1
	granularity := 2

	modelID, err := client.GenSegment2D(
		"model_id_from_gen_360", // model_id from Gen360
		"VISVISE-Segment-V1.0.0", // algorithm model
		"my_segment",
		nil,       // inputView
		&splitType,    // splitType: Front view
		&granularity, // granularity: Medium
		"segment by body parts",
		func(thought string) {
			fmt.Println("Thinking:", thought)
		},
	)
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
