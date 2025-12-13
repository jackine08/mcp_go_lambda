package main

import (
	"context"
	"testing"
)

func TestHandleInitialize(t *testing.T) {
	server := NewServer()
	ctx := context.Background()

	request := MCPRequest{
		JsonRPC: "2.0",
		Method:  "initialize",
		Params:  map[string]interface{}{},
		ID:      1,
	}

	response := server.Handle(ctx, request)

	if response.Error != nil {
		t.Errorf("Expected no error, got: %v", response.Error)
	}

	if response.Result == nil {
		t.Errorf("Expected result, got nil")
	}

	if response.ID != 1 {
		t.Errorf("Expected ID 1, got: %v", response.ID)
	}
}

func TestHandleResourcesList(t *testing.T) {
	server := NewServer()
	ctx := context.Background()

	request := MCPRequest{
		JsonRPC: "2.0",
		Method:  "resources/list",
		Params:  map[string]interface{}{},
		ID:      2,
	}

	response := server.Handle(ctx, request)

	if response.Error != nil {
		t.Errorf("Expected no error, got: %v", response.Error)
	}

	if response.Result == nil {
		t.Errorf("Expected result, got nil")
	}
}

func TestHandleToolsList(t *testing.T) {
	server := NewServer()
	ctx := context.Background()

	request := MCPRequest{
		JsonRPC: "2.0",
		Method:  "tools/list",
		Params:  map[string]interface{}{},
		ID:      3,
	}

	response := server.Handle(ctx, request)

	if response.Error != nil {
		t.Errorf("Expected no error, got: %v", response.Error)
	}

	if response.Result == nil {
		t.Errorf("Expected result, got nil")
	}
}

func TestHandleMethodNotFound(t *testing.T) {
	server := NewServer()
	ctx := context.Background()

	request := MCPRequest{
		JsonRPC: "2.0",
		Method:  "unknown_method",
		Params:  map[string]interface{}{},
		ID:      4,
	}

	response := server.Handle(ctx, request)

	if response.Error == nil {
		t.Errorf("Expected error for unknown method")
	}

	if response.Error.Code != -32601 {
		t.Errorf("Expected error code -32601, got: %d", response.Error.Code)
	}
}
