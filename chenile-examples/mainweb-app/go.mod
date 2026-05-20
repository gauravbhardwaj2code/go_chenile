module mainweb-app

go 1.26

toolchain go1.26.3

require (
	customer-service v0.0.0
	github.com/gauravbhardwaj2code/go_chenile/chenile-framework/packager v0.1.0
	order-service v0.0.0
)

require (
	github.com/gauravbhardwaj2code/go_chenile/chenile-framework/base v0.1.1 // indirect
	github.com/gauravbhardwaj2code/go_chenile/chenile-framework/chenile v0.1.1 // indirect
	github.com/gauravbhardwaj2code/go_chenile/chenile-framework/core v0.1.1 // indirect
	github.com/gauravbhardwaj2code/go_chenile/chenile-framework/http v0.1.1 // indirect
	github.com/gauravbhardwaj2code/go_chenile/chenile-framework/owiz v0.1.1 // indirect
)

replace customer-service => ../customer-service

replace order-service => ../order-service

replace github.com/gauravbhardwaj2code/go_chenile/chenile-framework/chenile => ../../chenile-framework/chenile

replace github.com/gauravbhardwaj2code/go_chenile/chenile-framework/core => ../../chenile-framework/core

replace github.com/gauravbhardwaj2code/go_chenile/chenile-framework/packager => ../../chenile-framework/packager

replace github.com/gauravbhardwaj2code/go_chenile/chenile-framework/http => ../../chenile-framework/http

replace github.com/gauravbhardwaj2code/go_chenile/chenile-framework/base => ../../chenile-framework/base

replace github.com/gauravbhardwaj2code/go_chenile/chenile-framework/owiz => ../../chenile-framework/owiz
