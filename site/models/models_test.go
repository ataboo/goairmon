package models

import (
	"reflect"
	"testing"
)

func TestEncodeDecodeFlashBag(t *testing.T) {
	bag := &FlashBag{
		Success: []string{"Success one", "Success two"},
		Info:    []string{"Info one", "Info two"},
		Error:   []string{"Error one", "Error two"},
	}

	strVal := bag.Encode()

	decoded := &FlashBag{}

	err := decoded.Decode(strVal)
	if err != nil {
		t.Error("unexpected error", err)
	}

	if !reflect.DeepEqual(bag, decoded) {
		t.Errorf("should match: %+v, %+v", bag, decoded)
	}
}
