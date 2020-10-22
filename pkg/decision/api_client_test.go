package decision

import (
	"testing"

	"github.com/abtasty/flagship-go-sdk/v2/pkg/decisionapi"
)

var testEnvID = "env_id_test"
var realEnvID = "blvo2kijq6pg023l8edg"

func TestNewAPIClient(t *testing.T) {
	_, err := NewAPIClient(testEnvID)

	if err == nil {
		t.Error("Api client V2 without API Key should fail")
	}

	client, _ := NewAPIClient(testEnvID, decisionapi.APIKey("test_api_key"))

	if client == nil {
		t.Error("Api client V2 with API Key should not fail")
	}

	_, err = NewAPIClient(testEnvID, decisionapi.APIVersion(2))

	if err == nil {
		t.Error("Api client V2 without API Key should fail")
	}
}

func TestGetModifications(t *testing.T) {
	client, _ := NewAPIClient(testEnvID, decisionapi.APIKey("test_api_key"))

	if client == nil {
		t.Error("Api client tracking should not be nil")
	}

	_, err := client.GetModifications("testID", nil)

	if err == nil {
		t.Error("Error should be raised for empty context")
	}
}
