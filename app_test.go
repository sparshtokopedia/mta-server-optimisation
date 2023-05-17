package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

//In the TestGetInefficientInstance function, we directly test the getInefficientInstance function
//to ensure it returns the correct inefficient instances.
//In the TestGetInstanceName function, we simulate an HTTP request using httptest package and
//verify the response. We set the X environment variable, call the handler function, and
//compare the response with the expected result from getInefficientInstance function.

func TestMain(m *testing.M) {
	if err := loadMockData(); err != nil {
		log.Fatal("Failed to load mock data: ", err)
	}
	os.Exit(m.Run())
}

func TestGetInefficientInstance(t *testing.T) {
	threshold := 1
	expectedResult := []string{"mta-prod-1", "mta-prod-3"}

	result := getInefficientInstance(threshold)

	if len(result) != len(expectedResult) {
		t.Errorf("Expected %d inefficient instances, but got %d", len(expectedResult), len(result))
	}

	for i, instance := range result {
		if instance != expectedResult[i] {
			t.Errorf("Expected inefficient instance '%s', but got '%s'", expectedResult[i], instance)
		}
	}
}

func TestGetInstanceName(t *testing.T) {
	testCases := []struct {
		Name           string
		Request        *http.Request
		ExpectedResult []string
	}{
		{
			Name: "Valid Request",
			Request: httptest.NewRequest(http.MethodGet, "/mta-hosting-optimizer", nil).
				WithContext(setEnvContext("X", "2")),
			ExpectedResult: []string{"mta-prod-3"},
		},
		{
			Name: "Invalid Threshold",
			Request: httptest.NewRequest(http.MethodGet, "/mta-hosting-optimizer", nil).
				WithContext(setEnvContext("X", "invalid")),
			ExpectedResult: nil,
		},
		{
			Name: "Missing Threshold",
			Request: httptest.NewRequest(http.MethodGet, "/mta-hosting-optimizer", nil).
				WithContext(setEnvContext("", "")),
			ExpectedResult: []string{"mta-prod-1", "mta-prod-3"},
		},
		{
			Name: "Non-Get Request",
			Request: httptest.NewRequest(http.MethodPost, "/mta-hosting-optimizer", nil).
				WithContext(setEnvContext("X", "2")),
			ExpectedResult: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			w := httptest.NewRecorder()
			getInstanceName(w, tc.Request)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
			}

			var response GetHostNameResponse
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Errorf("Error decoding response body: %s", err)
			}

			if len(response.Result) != len(tc.ExpectedResult) {
				t.Errorf("Expected %d result(s), but got %d", len(tc.ExpectedResult), len(response.Result))
			}

			for i, instance := range response.Result {
				if instance != tc.ExpectedResult[i] {
					t.Errorf("Expected instance '%s', but got '%s'", tc.ExpectedResult[i], instance)
				}
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	testCases := []struct {
		Name           string
		Key            string
		DefaultValue   string
		ExpectedResult string
	}{
		{
			Name:           "Existing Key",
			Key:            "X",
			DefaultValue:   "1",
			ExpectedResult: "1",
		},
		{
			Name:           "Non-Existing Key",
			Key:            "Y",
			DefaultValue:   "2",
			ExpectedResult: "2",
		},
		{
			Name:           "Empty Key",
			Key:            "",
			DefaultValue:   "3",
			ExpectedResult: "3",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result := getEnv(tc.Key, tc.DefaultValue)
			if result != tc.ExpectedResult {
				t.Errorf("getEnv returned wrong result for key '%s': got '%s', want '%s'", tc.Key, result, tc.ExpectedResult)
			}
		})
	}
}

func TestSetEnvContext(t *testing.T) {
	testCases := []struct {
		Name           string
		Key            string
		Value          string
		ExpectedResult context.Context
	}{
		{
			Name:           "Valid Key and Value",
			Key:            "X",
			Value:          "1",
			ExpectedResult: context.WithValue(context.Background(), "env", map[string]string{"X": "1"}),
		},
		{
			Name:           "Empty Key and Value",
			Key:            "",
			Value:          "",
			ExpectedResult: context.WithValue(context.Background(), "env", map[string]string{}),
		},
		{
			Name:           "Multiple Key-Value Pairs",
			Key:            "Y",
			Value:          "2",
			ExpectedResult: context.WithValue(context.WithValue(context.Background(), "env", map[string]string{"X": "1"}), "env", map[string]string{"X": "1", "Y": "2"}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result := setEnvContext(tc.Key, tc.Value)
			if !reflect.DeepEqual(result, tc.ExpectedResult) {
				t.Errorf("setEnvContext returned wrong result for key '%s' and value '%s'", tc.Key, tc.Value)
			}
		})
	}
}

func setEnvContext(key, value string) context.Context {
	env := make(map[string]string)
	env[key] = value
	return context.WithValue(context.Background(), "env", env)
}
