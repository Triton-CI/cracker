package main

import (
	"bytes"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	extism "github.com/extism/go-sdk"
)

func TestStorePlugin(t *testing.T) {
	// Clear the plugins map before testing
	plugins = make(map[string]*extism.Plugin)

	// Create a mock plugin
	plugin := &extism.Plugin{}

	StorePlugin(plugin)

	if len(plugins) != 1 {
		t.Errorf("Expected 1 plugin in storage, got %d", len(plugins))
	}

	if plugins["code"] != plugin {
		t.Error("Plugin not stored correctly")
	}
	//t.Error("Plugin not stored correctly")
}

func TestGetPlugin_Success(t *testing.T) {
	// Clear and setup
	plugins = make(map[string]*extism.Plugin)
	expectedPlugin := &extism.Plugin{}
	plugins["code"] = expectedPlugin

	plugin, err := GetPlugin()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if &plugin == nil {
		t.Error("Expected plugin to be returned")
	}
}

func TestGetPlugin_NotFound(t *testing.T) {
	// Clear the plugins map
	plugins = make(map[string]*extism.Plugin)

	_, err := GetPlugin()

	if err == nil {
		t.Error("Expected error when plugin not found")
	}

	expectedErr := errors.New("🔴 no plugin")
	if err.Error() != expectedErr.Error() {
		t.Errorf("Expected error '%v', got '%v'", expectedErr, err)
	}
}

func TestGetBytesBody(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		expected []byte
	}{
		{
			name:     "Simple JSON body",
			body:     `{"key":"value"}`,
			expected: []byte(`{"key":"value"}`),
		},
		{
			name:     "Empty body",
			body:     "",
			expected: []byte{},
		},
		{
			name:     "Complex JSON",
			body:     `{"model":"test","system":"content","user":"message"}`,
			expected: []byte(`{"model":"test","system":"content","user":"message"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyReader := bytes.NewBufferString(tt.body)
			req := httptest.NewRequest("POST", "/", bodyReader)
			req.ContentLength = int64(len(tt.body))

			result := GetBytesBody(req)

			if !bytes.Equal(result, tt.expected) {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetBytesBody_WithRealRequest(t *testing.T) {
	body := "test data"
	req := httptest.NewRequest("POST", "/", bytes.NewBufferString(body))
	req.ContentLength = int64(len(body))

	result := GetBytesBody(req)

	if string(result) != body {
		t.Errorf("Expected '%s', got '%s'", body, string(result))
	}
}

func TestGetBytesBody_ZeroLength(t *testing.T) {
	req := httptest.NewRequest("POST", "/", nil)
	req.ContentLength = 0

	result := GetBytesBody(req)

	if len(result) != 0 {
		t.Errorf("Expected empty byte slice, got %v with length %d", result, len(result))
	}
}

func TestConcurrentPluginAccess(t *testing.T) {
	// Test concurrent access to plugins map
	plugins = make(map[string]*extism.Plugin)
	plugin := &extism.Plugin{}
	StorePlugin(plugin)

	done := make(chan bool)

	// Spawn multiple goroutines to access the plugin
	for i := 0; i < 10; i++ {
		go func() {
			_, err := GetPlugin()
			if err != nil {
				t.Errorf("Concurrent access failed: %v", err)
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestHTTPHandler_MissingPlugin(t *testing.T) {
	// Clear plugins to simulate missing plugin
	plugins = make(map[string]*extism.Plugin)

	req := httptest.NewRequest("POST", "/", bytes.NewBufferString(`{"test":"data"}`))
	req.ContentLength = int64(len(`{"test":"data"}`))
	httptest.NewRecorder()

	// Note: We can't easily test the full handler without mocking extism.Plugin.Call
	// This test verifies the GetPlugin error case
	_, err := GetPlugin()
	if err == nil {
		t.Error("Expected error when plugin is missing")
	}
}

func TestGetBytesBody_LargeBody(t *testing.T) {
	// Test with a larger body
	largeBody := make([]byte, 1024*1024) // 1MB
	for i := range largeBody {
		largeBody[i] = byte(i % 256)
	}

	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(largeBody))
	req.ContentLength = int64(len(largeBody))

	result := GetBytesBody(req)

	if len(result) != len(largeBody) {
		t.Errorf("Expected %d bytes, got %d", len(largeBody), len(result))
	}

	if !bytes.Equal(result, largeBody) {
		t.Error("Large body content mismatch")
	}
}

func TestGetBytesBody_PartialRead(t *testing.T) {
	// Test when ContentLength is set but body has less data
	body := "short"
	req := httptest.NewRequest("POST", "/", bytes.NewBufferString(body))
	req.ContentLength = 100 // Claim more than actual

	result := GetBytesBody(req)

	// Should allocate 100 bytes but only read what's available
	if len(result) != 100 {
		t.Errorf("Expected allocated length of 100, got %d", len(result))
	}

	// First 5 bytes should match "short"
	if string(result[:5]) != body {
		t.Errorf("Expected '%s', got '%s'", body, string(result[:5]))
	}
}

func TestMultiplePluginStores(t *testing.T) {
	plugins = make(map[string]*extism.Plugin)

	plugin1 := &extism.Plugin{}
	plugin2 := &extism.Plugin{}

	StorePlugin(plugin1)
	StorePlugin(plugin2)

	// Should only have 1 entry since both use "code" key
	if len(plugins) != 1 {
		t.Errorf("Expected 1 plugin, got %d", len(plugins))
	}

	// Latest stored should be plugin2
	if plugins["code"] != plugin2 {
		t.Error("Expected latest stored plugin to be plugin2")
	}
}

func BenchmarkGetBytesBody(b *testing.B) {
	body := "benchmark test data"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/", bytes.NewBufferString(body))
		req.ContentLength = int64(len(body))
		GetBytesBody(req)
	}
}

func BenchmarkStoreAndGetPlugin(b *testing.B) {
	plugins = make(map[string]*extism.Plugin)
	plugin := &extism.Plugin{}
	StorePlugin(plugin)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetPlugin()
	}
}

// TestGetBytesBody_NilBody tests handling of nil body
func TestGetBytesBody_NilBody(t *testing.T) {
	req := httptest.NewRequest("POST", "/", nil)
	req.ContentLength = 0

	result := GetBytesBody(req)

	if result == nil {
		t.Error("Expected non-nil result")
	}

	if len(result) != 0 {
		t.Errorf("Expected empty slice, got length %d", len(result))
	}
}

// TestGetBytesBody_ReadError tests when Body.Read might not read all expected bytes
func TestGetBytesBody_BodyReader(t *testing.T) {
	testData := []byte("test content")
	req := httptest.NewRequest("POST", "/", io.NopCloser(bytes.NewReader(testData)))
	req.ContentLength = int64(len(testData))

	result := GetBytesBody(req)

	if !bytes.Equal(result, testData) {
		t.Errorf("Expected %s, got %s", testData, result)
	}
}
