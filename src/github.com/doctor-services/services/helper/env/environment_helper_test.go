package env

import (
	"os"
	"testing"
)

func TestGetEnvString(t *testing.T) {
	os.Setenv("test", "test")
	expectedValue := GetEnvString("test", "default")
	if expectedValue != "test" {
		t.Fatalf("Expected %v but got %v", "test", expectedValue)
	}
	os.Unsetenv("NotExist")
	defaultValue := GetEnvString("NotExist", "default")
	if defaultValue != "default" {
		t.Fatalf("Expected %v but got %v", "test", defaultValue)
	}
}
