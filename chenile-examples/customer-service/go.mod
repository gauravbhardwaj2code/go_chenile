module customer-service

go 1.26

toolchain go1.26.3

require (
	base v0.0.0
	bdd-utils v0.0.0
	chenile v0.0.0
	config v0.0.0
	packager v0.0.0
)

require (
	core v0.0.0 // indirect
	github.com/cucumber/gherkin/go/v26 v26.2.0 // indirect
	github.com/cucumber/godog v0.15.1 // indirect
	github.com/cucumber/messages/go/v21 v21.0.1 // indirect
	github.com/gofrs/uuid v4.3.1+incompatible // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-memdb v1.3.4 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/spf13/pflag v1.0.7 // indirect
	http v0.0.0 // indirect
	owiz v0.0.0 // indirect
)

replace bdd-utils => ../../chenile-framework/bdd-utils

replace base => ../../chenile-framework/base

replace chenile => ../../chenile-framework/chenile

replace config => ../../chenile-framework/config

replace core => ../../chenile-framework/core

replace http => ../../chenile-framework/http

replace owiz => ../../chenile-framework/owiz

replace packager => ../../chenile-framework/packager
