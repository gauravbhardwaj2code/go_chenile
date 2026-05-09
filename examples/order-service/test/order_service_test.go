package test

import (
	"io"
	"testing"

	"github.com/ajapro/chenile-go/packager"
	godogtest "github.com/ajapro/chenile-go/test/godog"

	"order-service/order"
)

func TestCreateOrder(t *testing.T) {
	app, err := packager.NewWebApp(packager.Module{Name: "order", Register: order.Register})
	if err != nil {
		t.Fatal(err)
	}

	status := godogtest.Suite{
		Name:         "order-service",
		Router:       app.Router,
		FeaturePaths: []string{"features/order.feature"},
		TestingT:     t,
		Output:       io.Discard,
	}.Run()
	if status != 0 {
		t.Fatalf("godog suite failed with status %d", status)
	}
}
