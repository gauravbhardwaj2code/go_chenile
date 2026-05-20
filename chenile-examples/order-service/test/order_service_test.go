package test

import (
	"io"
	"testing"

	"packager"
	godogtest "bdd-utils/godog"

	"order-service/order/module"
)

func TestCreateOrder(t *testing.T) {
	app, err := packager.NewChenileWebApp(module.New())
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
