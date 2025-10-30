package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/semaphoreui/semaphore/db"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationMatch(t *testing.T) {
	body := []byte("{\"hook_id\": 4856239453}")
	var header = make(http.Header)
	matched := Match(db.IntegrationMatcher{
		ID:            0,
		Name:          "Test",
		IntegrationID: 0,
		MatchType:     db.IntegrationMatchBody,
		Method:        db.IntegrationMatchMethodEquals,
		BodyDataType:  db.IntegrationBodyDataJSON,
		Key:           "hook_id",
		Value:         "4856239453",
	}, header, body)

	assert.True(t, matched)
}

func TestGetTaskDefinitionSuccess(t *testing.T) {
	integration := db.Integration{
		ID:         11,
		ProjectID:  22,
		TemplateID: 33,
		TaskParams: &db.TaskParams{
			ProjectID:   22,
			Environment: `{"existing":"value"}`,
			Params:      db.MapStringAnyField{"original": "keep"},
		},
	}

	header := make(http.Header)
	header.Set("X-Env", "header-value")
	payload := []byte(`{"data":{"param":"payload-value"}}`)

	extractorCalled := false
	task, err := GetTaskDefinition(integration, payload, header, func(projectID, integrationID int) ([]db.IntegrationExtractValue, error) {
		extractorCalled = true

		if projectID != integration.ProjectID {
			t.Fatalf("expected projectID %d, got %d", integration.ProjectID, projectID)
		}
		if integrationID != integration.ID {
			t.Fatalf("expected integrationID %d, got %d", integration.ID, integrationID)
		}

		return []db.IntegrationExtractValue{
			{
				VariableType: db.IntegrationVariableEnvironment,
				ValueSource:  db.IntegrationExtractHeaderValue,
				Key:          "X-Env",
				Variable:     "HOOK_ENV",
			},
			{
				VariableType: db.IntegrationVariableTaskParam,
				ValueSource:  db.IntegrationExtractBodyValue,
				BodyDataType: db.IntegrationBodyDataJSON,
				Key:          "data.param",
				Variable:     "payloadParam",
			},
		}, nil
	})

	assert.NoError(t, err)
	assert.True(t, extractorCalled)

	if assert.NotNil(t, task.IntegrationID) {
		assert.Equal(t, integration.ID, *task.IntegrationID)
	}

	assert.Equal(t, integration.ProjectID, task.ProjectID)
	assert.Equal(t, integration.TemplateID, task.TemplateID)
	assert.NotEmpty(t, task.Environment)

	var env map[string]any
	if assert.NoError(t, json.Unmarshal([]byte(task.Environment), &env)) {
		assert.Equal(t, "value", env["existing"])
		assert.Equal(t, "header-value", env["HOOK_ENV"])
	}

	if assert.NotNil(t, task.Params) {
		if assert.Contains(t, task.Params, "original") {
			assert.Equal(t, "keep", task.Params["original"])
		}

		if assert.Contains(t, task.Params, "payloadParam") {
			payloadParam, ok := task.Params["payloadParam"].(string)
			assert.True(t, ok)
			assert.Equal(t, "payload-value", payloadParam)
		}
	}
}

func TestGetTaskDefinitionExtractorError(t *testing.T) {
	integration := db.Integration{
		ID:         44,
		ProjectID:  55,
		TemplateID: 66,
	}

	header := make(http.Header)
	payload := []byte(`{}`)

	expectedErr := errors.New("extractor failure")

	extractorCalled := false
	task, err := GetTaskDefinition(integration, payload, header, func(projectID, integrationID int) ([]db.IntegrationExtractValue, error) {
		extractorCalled = true
		return nil, expectedErr
	})

	assert.True(t, extractorCalled)
	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	assert.Nil(t, task.IntegrationID)
}

func TestGetTaskDefinitionInvalidEnvironmentJSON(t *testing.T) {
	integration := db.Integration{
		ID:         77,
		ProjectID:  88,
		TemplateID: 99,
		TaskParams: &db.TaskParams{
			ProjectID:   88,
			Environment: "{not-json}",
			Params:      db.MapStringAnyField{},
		},
	}

	header := make(http.Header)
	payload := []byte(`{}`)

	_, err := GetTaskDefinition(integration, payload, header, func(projectID, integrationID int) ([]db.IntegrationExtractValue, error) {
		return nil, nil
	})

	assert.Error(t, err)
}

func TestGetTaskDefinitionIntegrationWithoutTaskParams(t *testing.T) {
	integration := db.Integration{
		ID:         44,
		ProjectID:  55,
		TemplateID: 66,
	}

	header := make(http.Header)
	payload := []byte(`{}`)
	extractorCalled := false
	task, err := GetTaskDefinition(integration, payload, header, func(projectID, integrationID int) ([]db.IntegrationExtractValue, error) {
		extractorCalled = true

		if projectID != integration.ProjectID {
			t.Fatalf("expected projectID %d, got %d", integration.ProjectID, projectID)
		}
		if integrationID != integration.ID {
			t.Fatalf("expected integrationID %d, got %d", integration.ID, integrationID)
		}

		return []db.IntegrationExtractValue{
			{
				VariableType: db.IntegrationVariableEnvironment,
				ValueSource:  db.IntegrationExtractHeaderValue,
				Key:          "X-Env",
				Variable:     "HOOK_ENV",
			},
			{
				VariableType: db.IntegrationVariableTaskParam,
				ValueSource:  db.IntegrationExtractBodyValue,
				BodyDataType: db.IntegrationBodyDataJSON,
				Key:          "data.param",
				Variable:     "payloadParam",
			},
		}, nil
	})

	assert.True(t, extractorCalled)
	assert.Nil(t, err)
	assert.NotNil(t, task)
}
