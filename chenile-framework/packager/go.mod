module github.com/gauravbhardwaj2code/go_chenile/chenile-framework/packager

go 1.26

toolchain go1.26.3

require (
	github.com/gauravbhardwaj2code/go_chenile/chenile-framework/chenile v0.1.0
	github.com/gauravbhardwaj2code/go_chenile/chenile-framework/core v0.1.0
	github.com/gauravbhardwaj2code/go_chenile/chenile-framework/http v0.1.0
)

require (
	github.com/gauravbhardwaj2code/go_chenile/chenile-framework/base v0.1.0 // indirect
	github.com/gauravbhardwaj2code/go_chenile/chenile-framework/owiz v0.1.0 // indirect
)

replace github.com/gauravbhardwaj2code/go_chenile/chenile-framework/chenile => ../chenile

replace github.com/gauravbhardwaj2code/go_chenile/chenile-framework/core => ../core

replace github.com/gauravbhardwaj2code/go_chenile/chenile-framework/http => ../http

replace github.com/gauravbhardwaj2code/go_chenile/chenile-framework/base => ../base

replace github.com/gauravbhardwaj2code/go_chenile/chenile-framework/owiz => ../owiz
