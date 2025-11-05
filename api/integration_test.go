package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/semaphoreui/semaphore/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtract_HeaderAndCaseInsensitive(t *testing.T) {
	h := http.Header{}
	h.Set("x-token", "abc123") // lower-case to verify case-insensitive get

	values := []db.IntegrationExtractValue{
		{
			Name:         "Token header",
			ValueSource:  db.IntegrationExtractHeaderValue,
			Key:          "X-Token", // different case
			Variable:     "TOKEN",
			VariableType: db.IntegrationVariableEnvironment,
		},
	}

	got := Extract(values, h, nil)

	require.Equal(t, "abc123", got["TOKEN"], "TOKEN header value should match")
}

func TestExtract_JSONBody_VariousTypesAndMissing(t *testing.T) {
	payload := []byte(`{
		"num": 42,
		"str": "hello",
		"bool": true,
		"nullv": null,
		"obj": {"k":"v"},
		"arr": [1,2,3],
		"nested": {"items":[{"c":123},{"c":"str"}]}
	}`)

	values := []db.IntegrationExtractValue{
		{ // number coerced to string via fmt.Sprintf
			ValueSource:  db.IntegrationExtractBodyValue,
			BodyDataType: db.IntegrationBodyDataJSON,
			Key:          "num",
			Variable:     "NUM",
		},
		{ // string stays same content
			ValueSource:  db.IntegrationExtractBodyValue,
			BodyDataType: db.IntegrationBodyDataJSON,
			Key:          "str",
			Variable:     "STR",
		},
		{ // boolean -> "true"
			ValueSource:  db.IntegrationExtractBodyValue,
			BodyDataType: db.IntegrationBodyDataJSON,
			Key:          "bool",
			Variable:     "BOOL",
		},
		{ // null should not be set (Find returns nil or we skip when nil)
			ValueSource:  db.IntegrationExtractBodyValue,
			BodyDataType: db.IntegrationBodyDataJSON,
			Key:          "nullv",
			Variable:     "NULLV",
		},
		{ // array will be formatted with %v, expect Go-like format
			ValueSource:  db.IntegrationExtractBodyValue,
			BodyDataType: db.IntegrationBodyDataJSON,
			Key:          "arr",
			Variable:     "ARR",
		},
		{ // object -> formatted map with %v
			ValueSource:  db.IntegrationExtractBodyValue,
			BodyDataType: db.IntegrationBodyDataJSON,
			Key:          "obj",
			Variable:     "OBJ",
		},
		{ // missing key should not create an entry
			ValueSource:  db.IntegrationExtractBodyValue,
			BodyDataType: db.IntegrationBodyDataJSON,
			Key:          "missing",
			Variable:     "MISSING",
		},
		{ // nested array index path
			ValueSource:  db.IntegrationExtractBodyValue,
			BodyDataType: db.IntegrationBodyDataJSON,
			Key:          "nested.items.[0].c",
			Variable:     "NESTED_C",
		},
		{ // first element of arr
			ValueSource:  db.IntegrationExtractBodyValue,
			BodyDataType: db.IntegrationBodyDataJSON,
			Key:          "arr.[0]",
			Variable:     "ARR0",
		},
	}

	got := Extract(values, http.Header{}, payload)

	// Basic scalar assertions
	assert.Equal(t, "42", got["NUM"], "NUM should equal stringified number")
	assert.Equal(t, "hello", got["STR"], "STR should match")
	assert.Equal(t, "true", got["BOOL"], "BOOL should be string 'true'")

	// Indexed lookups
	assert.Equal(t, "123", got["NESTED_C"], "NESTED_C should equal nested.items[0].c")
	assert.Equal(t, "1", got["ARR0"], "ARR0 should equal arr[0]")

	// Null should be absent
	assert.NotContains(t, got, "NULLV", "NULLV should not be present for null JSON value")

	// Array/object string formats: we assert non-empty presence rather than exact formatting,
	// because %v formatting of gojsonq return types may vary across versions.
	assert.Contains(t, got, "ARR", "ARR key should be present")
	assert.NotEmpty(t, got["ARR"], "ARR value should be non-empty")
	assert.Contains(t, got, "OBJ", "OBJ key should be present")
	assert.NotEmpty(t, got["OBJ"], "OBJ value should be non-empty")

	// Missing should not appear
	assert.NotContains(t, got, "MISSING", "MISSING should not be present for missing key")
}

func TestExtract_BodyString_ReturnsFullPayload(t *testing.T) {
	payload := []byte("raw body data here")
	values := []db.IntegrationExtractValue{
		{
			ValueSource:  db.IntegrationExtractBodyValue,
			BodyDataType: db.IntegrationBodyDataString,
			Variable:     "BODY",
			Key:          "ignored",
		},
	}
	got := Extract(values, http.Header{}, payload)
	if got["BODY"] != string(payload) {
		t.Fatalf("expected BODY to equal full payload; got %q", got["BODY"])
	}
}

func TestExtract_MalformedJSON_SkipsSetting(t *testing.T) {
	payload := []byte("{not: valid json}")
	values := []db.IntegrationExtractValue{
		{
			ValueSource:  db.IntegrationExtractBodyValue,
			BodyDataType: db.IntegrationBodyDataJSON,
			Variable:     "BAD",
			Key:          "a.b",
		},
	}
	got := Extract(values, http.Header{}, payload)
	if _, ok := got["BAD"]; ok {
		t.Fatalf("expected BAD to be absent for malformed JSON payload")
	}
}

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
