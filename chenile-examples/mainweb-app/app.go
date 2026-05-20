package mainweb

import (
	"github.com/gauravbhardwaj2code/go_chenile/chenile-framework/packager"

	customermodule "customer-service/customer/module"
	ordermodule "order-service/order/module"
)

func NewApp() (*packager.ChenileApp, error) {
	return packager.NewChenileWebApp(
		customermodule.New(),
		ordermodule.New(),
	)
}
