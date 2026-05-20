package test

import (
	"io"
	"testing"

	godogtest "github.com/gauravbhardwaj2code/go_chenile/chenile-framework/bdd-utils/godog"
	"github.com/gauravbhardwaj2code/go_chenile/chenile-framework/packager"

	"state-order-service/order"
)

func TestStateOrderWorkflow(t *testing.T) {
	app, err := packager.NewWebApp(packager.Module{Name: "state-order", Register: order.Register})
	if err != nil {
		t.Fatal(err)
	}

	status := godogtest.Suite{
		Name:         "state-order-service",
		Router:       app.Router,
		FeaturePaths: []string{"features/state_order.feature"},
		TestingT:     t,
		Output:       io.Discard,
	}.Run()
	if status != 0 {
		t.Fatalf("godog suite failed with status %d", status)
	}
}
