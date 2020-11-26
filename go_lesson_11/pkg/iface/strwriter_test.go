package iface

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestWriteString(t *testing.T) {
	var got bytes.Buffer
	writer := io.Writer(&got)
	gotErr := WriteString(&writer, 1, "one", true, "two", nil)
	want := *bytes.NewBufferString("onetwo")

	if !reflect.DeepEqual(got, want) || gotErr != nil {
		t.Errorf("got %v, want %v", got, want)
	}
}
