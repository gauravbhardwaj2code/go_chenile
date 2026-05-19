module mainweb-app

go 1.26

toolchain go1.26.3

require (
	customer-service v0.0.0
	order-service v0.0.0
	packager v0.0.0
)

require (
	base v0.0.0 // indirect
	core v0.0.0 // indirect
	http v0.0.0 // indirect
	owiz v0.0.0 // indirect
)

replace customer-service => ../customer-service

replace order-service => ../order-service

replace core => ../../chenile-framework/core

replace packager => ../../chenile-framework/packager

replace http => ../../chenile-framework/http

replace base => ../../chenile-framework/base

replace owiz => ../../chenile-framework/owiz
