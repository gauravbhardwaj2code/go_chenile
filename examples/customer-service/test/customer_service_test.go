package test

import (
	"io"
	"testing"

	"github.com/ajapro/chenile-go/packager"
	godogtest "github.com/ajapro/chenile-go/test/godog"

	"customer-service/customer"
)

func TestCreateCustomer(t *testing.T) {
	app, err := packager.NewWebApp(packager.Module{Name: "customer", Register: customer.Register})
	if err != nil {
		t.Fatal(err)
	}

	status := godogtest.Suite{
		Name:         "customer-service",
		Router:       app.Router,
		FeaturePaths: []string{"features/customer.feature"},
		TestingT:     t,
		Output:       io.Discard,
	}.Run()
	if status != 0 {
		t.Fatalf("godog suite failed with status %d", status)
	}
}
