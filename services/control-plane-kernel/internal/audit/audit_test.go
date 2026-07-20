package audit

import "testing"

func TestRedactMetadataRemovesSecretLikeKeys(t *testing.T) {
	metadata := RedactMetadata(map[string]string{
		"connectorType":   "postgres",
		"password":        "nope",
		"api_token":       "nope",
		"secretReference": "nope",
	})

	if len(metadata) != 1 || metadata["connectorType"] != "postgres" {
		t.Fatalf("unexpected metadata after redaction: %#v", metadata)
	}
}
