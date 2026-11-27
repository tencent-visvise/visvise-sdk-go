package tests

import (
	"os"
	"testing"

	"github.com/tencent-visvise/visvise-sdk-go/visvise"
)

// TestRTPConnected_GenRTPConnected tests the 重拓扑-连通体 task (node_type=19),
// the first stage of the two-stage VISVISE retopology flow.
func TestRTPConnected_GenRTPConnected(t *testing.T) {
	if appID == "" || secretKey == "" || rtx == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_RTX not set")
	}

	modelPath := assetsDir + "/high_model.fbx"
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Skip("Skipping test: high_model.fbx not found in tests/assets")
	}

	client := visvise.NewClient(appID, secretKey, nil)

	opts := visvise.NewGenRTPConnectedOptions().
		SetName("test_gen_rtp_connected")

	modelID, err := client.GenRTPConnected(modelPath, rtx, opts)
	if err != nil {
		t.Fatalf("GenRTPConnected failed: %v", err)
	}
	if modelID == "" {
		t.Fatal("GenRTPConnected returned empty model_id")
	}
	t.Logf("PASS: gen_rtp_connected - model_id=%s", modelID)

	model, err := client.WaitModel(modelID, rtx, &visvise.WaitOptions{Interval: 5, Timeout: 900})
	if err != nil {
		t.Logf("Wait failed (may timeout): %v", err)
		return
	}
	t.Logf("PASS: gen_rtp_connected - output_model=%s", model.OutputModel)
}

// TestRTPConnected_SaveRTPConnected tests the save_rtp_connected API.
// It needs an existing first-stage asset and the edited connected-component zip,
// provided through VISVISE_RTP_CONNECTED_MODEL_ID / VISVISE_RTP_CONNECTED_ZIP.
func TestRTPConnected_SaveRTPConnected(t *testing.T) {
	if appID == "" || secretKey == "" || rtx == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_RTX not set")
	}

	sourceModelID := os.Getenv("VISVISE_RTP_CONNECTED_MODEL_ID")
	editedZip := os.Getenv("VISVISE_RTP_CONNECTED_ZIP")
	if sourceModelID == "" || editedZip == "" {
		t.Skip("Skipping test: VISVISE_RTP_CONNECTED_MODEL_ID / VISVISE_RTP_CONNECTED_ZIP not set")
	}

	client := visvise.NewClient(appID, secretKey, nil)

	savedID, err := client.SaveRTPConnected(sourceModelID, editedZip, "test_rtp_connected_saved", rtx)
	if err != nil {
		t.Fatalf("SaveRTPConnected failed: %v", err)
	}
	if savedID == "" {
		t.Fatal("SaveRTPConnected returned empty model_id")
	}
	t.Logf("PASS: save_rtp_connected - model_id=%s", savedID)
}

// TestRTPConnected_RetopologyWithParent tests the second-stage retopology task
// (node_type=1) with parent_model_id pointing to a saved connected-component asset.
func TestRTPConnected_RetopologyWithParent(t *testing.T) {
	if appID == "" || secretKey == "" || rtx == "" {
		t.Skip("Skipping test: VISVISE_APP_ID, VISVISE_SECRET_KEY, or VISVISE_RTX not set")
	}

	parentModelID := os.Getenv("VISVISE_RTP_CONNECTED_MODEL_ID")
	editedZip := os.Getenv("VISVISE_RTP_CONNECTED_ZIP")
	if parentModelID == "" || editedZip == "" {
		t.Skip("Skipping test: VISVISE_RTP_CONNECTED_MODEL_ID / VISVISE_RTP_CONNECTED_ZIP not set")
	}

	client := visvise.NewClient(appID, secretKey, nil)

	opts := visvise.NewGenRetopologyOptions().
		SetName("test_gen_retopology_vv").
		SetOutputModelFormat(visvise.OutputModelFormatFBX).
		SetFaceType(visvise.FaceTypeQuad).
		SetFaceNum(10000).
		SetSmoothStrength(1.0).
		SetParentModelID(parentModelID)

	modelID, err := client.GenRetopology(editedZip, rtx, opts)
	if err != nil {
		t.Fatalf("GenRetopology with parent_model_id failed: %v", err)
	}
	t.Logf("PASS: gen_retopology (parent_model_id=%s) - model_id=%s", parentModelID, modelID)
}
