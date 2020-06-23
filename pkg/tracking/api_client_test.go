package tracking

import (
	"testing"

	"github.com/abtasty/flagship-go-sdk/pkg/model"
)

var testVisitorID = "test_visitor_id"
var testEnvID = "test_env_id"
var realEnvID = "blvo2kijq6pg023l8edg"

func TestNewAPIClient(t *testing.T) {
	client, _ := NewAPIClient(testEnvID)

	if client == nil {
		t.Error("Api client tracking should not be nil")
	}

	if client.urlTracking != defaultAPIURLTracking {
		t.Error("Api url should be set to default")
	}
}

func TestSendInternalHit(t *testing.T) {
	client, _ := NewAPIClient(testEnvID)
	err := client.SendHit(testVisitorID, nil)

	if err == nil {
		t.Error("Empty hit should return and err")
	}

	event := &model.EventHit{}
	event.SetBaseInfos(testEnvID, testVisitorID)

	err = client.SendHit(testVisitorID, event)

	if err == nil {
		t.Error("Invalid event hit should return error")
	}

	event.Action = "test_action"
	err = client.SendHit(testVisitorID, event)

	if err != nil {
		t.Errorf("Right hit should not return and err : %v", err)
	}
}

func TestActivate(t *testing.T) {
	client, _ := NewAPIClient(testEnvID)
	err := client.ActivateCampaign(model.ActivationHit{})

	if err == nil {
		t.Errorf("Expected error for empty request")
	}

	err = client.ActivateCampaign(model.ActivationHit{
		EnvironmentID:    testEnvID,
		VisitorID:        "test_vid",
		VariationGroupID: "vgid",
		VariationID:      "vid",
	})

	if err != nil {
		t.Errorf("Did not expect error for correct activation request. Got %v", err)
	}
}

func TestSendEvent(t *testing.T) {
	client, _ := NewAPIClient(realEnvID)
	err := client.SendEvent(model.Event{})

	if err == nil {
		t.Errorf("Expected error for empty event request")
	}

	err = client.SendEvent(model.Event{
		VisitorID: "test_vid",
		Type:      model.CONTEXT,
		Data:      model.Context{},
	})

	if err != nil {
		t.Errorf("Did not expect error for correct event request. Got %v", err)
	}
}
