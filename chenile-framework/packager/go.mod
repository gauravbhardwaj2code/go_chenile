module packager

go 1.22

require (
	core v0.0.0
	http v0.0.0
)

require (
	base v0.0.0 // indirect
	owiz v0.0.0 // indirect
)

replace core => ../core

replace http => ../http

replace base => ../base

replace owiz => ../owiz
