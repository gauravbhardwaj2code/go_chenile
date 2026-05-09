module mainweb-app

go 1.22

require (
	customer-service v0.0.0
	github.com/ajapro/chenile-go/packager v0.0.0
	order-service v0.0.0
)

require (
	github.com/ajapro/chenile-go/base v0.0.0 // indirect
	github.com/ajapro/chenile-go/core v0.0.0 // indirect
	github.com/ajapro/chenile-go/http v0.0.0 // indirect
	github.com/ajapro/chenile-go/owiz v0.0.0 // indirect
)

replace customer-service => ../customer-service

replace order-service => ../order-service

replace github.com/ajapro/chenile-go/base => ../../base

replace github.com/ajapro/chenile-go/core => ../../core

replace github.com/ajapro/chenile-go/http => ../../http

replace github.com/ajapro/chenile-go/owiz => ../../owiz

replace github.com/ajapro/chenile-go/packager => ../../packager

replace github.com/ajapro/chenile-go/test => ../../test
