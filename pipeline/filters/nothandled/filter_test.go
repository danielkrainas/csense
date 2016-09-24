package nothandled_test

import (
	"testing"

	"github.com/danielkrainas/csense/context"
	"github.com/danielkrainas/csense/pipeline"
	"github.com/danielkrainas/csense/pipeline/filters/nothandled"
)

func TestHandleMessage(t *testing.T) {
	f := nothandled.New()
	err := f.HandleMessage(context.Background(), pipeline.NewNullMessage())
	if err == nil {
		t.Error("no error returned for unhandled message")
	}
}
