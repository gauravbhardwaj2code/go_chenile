package mainweb

import (
	"packager"

	customermodule "customer-service/customer/module"
	ordermodule "order-service/order/module"
)

func NewApp() (*packager.ChenileApp, error) {
	return packager.NewChenileWebApp(
		customermodule.New(),
		ordermodule.New(),
	)
}
