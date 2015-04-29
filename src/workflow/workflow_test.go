package workflow

import (
	"path/filepath"
	"testing"

	"golang.org/x/net/context"
)

func TestWorkflow(t *testing.T) {
	testConf, err := filepath.Abs("../../test/workflow.yaml")
	if err != nil {
		t.Fatal(err)
	}

	counter := 0

	ok := func(ctx context.Context) (context.Context, string, error) {
		counter = counter + 1
		return ctx, "OK", nil
	}
	NewAgent(ok, "start", "step1", "step2", "step3", "step4")

	wf := New(testConf, context.Background())
	err = wf.Handle("start")
	if err != nil {
		t.Error(err)
	}
	if counter != 3 {
		t.Error(counter, "is not 3")
	}

	counter = 0
	err = wf.Handle("step1")
	if err != nil {
		t.Error(err)
	}
	if counter != 1 {
		t.Error(counter, "is not 1")
	}
}

func TestContextValueChange(t *testing.T) {
	testConf, err := filepath.Abs("../../test/workflow.yaml")
	if err != nil {
		t.Fatal(err)
	}

	okFn := func(ctx context.Context) (context.Context, string, error) {
		ctx = context.WithValue(ctx, "key", "value")
		return ctx, "OK", nil
	}
	NewAgent(okFn, "start", "step1", "step2", "step3", "step4")

	wf := New(testConf, context.Background())
	err = wf.Handle("start")
	if err != nil {
		t.Error(err)
	}
	value, ok := wf.Context.Value("key").(string)
	if !ok {
		t.Error("value is not string")
	}
	if value != "value" {
		t.Error("value is not value", value)
	}

}
