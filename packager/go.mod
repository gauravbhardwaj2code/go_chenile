module github.com/ajapro/chenile-go/packager

go 1.22

require (
	github.com/ajapro/chenile-go/core v0.0.0
	github.com/ajapro/chenile-go/http v0.0.0
)

require (
	github.com/ajapro/chenile-go/base v0.0.0 // indirect
	github.com/ajapro/chenile-go/owiz v0.0.0 // indirect
)

replace github.com/ajapro/chenile-go/base => ../base

replace github.com/ajapro/chenile-go/core => ../core

replace github.com/ajapro/chenile-go/http => ../http

replace github.com/ajapro/chenile-go/owiz => ../owiz
