package base58

import "testing"

func TestEncoding_EncodeToString(t *testing.T) {
	str := "hello world"

	enc := StdEncoding.EncodeToString([]byte(str))

	if enc != "StV1DL6CwTryKyV" {
		t.Errorf("Expected %s, got %s", "StV1DL6CwTryKyV", enc)
	}
}

func TestEncoding_Encode(t *testing.T) {
	str := "hello world"

	enc := StdEncoding.Encode([]byte(str))

	if string(enc) != "StV1DL6CwTryKyV" {
		t.Errorf("Expected %s, got %s", "StV1DL6CwTryKyV", string(enc))
	}
}

func TestEncoding_Decode(t *testing.T) {
	str := "StV1DL6CwTryKyV"

	dec, err := StdEncoding.Decode(str)
	if err != nil {
		t.Errorf("Error decoding string: %s", err.Error())
	}

	if string(dec) != "hello world" {
		t.Errorf("Expected %s, got %s", "hello world", string(dec))
	}
}

func TestEncoding_DecodeToString(t *testing.T) {
	str := "StV1DL6CwTryKyV"

	dec, err := StdEncoding.DecodeString(str)
	if err != nil {
		t.Errorf("Error decoding string: %s", err.Error())
	}

	if string(dec) != "hello world" {
		t.Errorf("Expected %s, got %s", "hello world", dec)
	}
}
