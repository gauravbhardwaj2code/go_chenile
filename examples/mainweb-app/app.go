package mainweb

import (
	"github.com/ajapro/chenile-go/packager"

	"customer-service/customer"
	"order-service/order"
)

func NewApp() (*packager.App, error) {
	return packager.NewWebApp(
		packager.Module{Name: "customer", Register: customer.Register},
		packager.Module{Name: "order", Register: order.Register},
	)
}
