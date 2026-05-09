module customer-service

go 1.22

require (
	core v0.0.0
	packager v0.0.0
	test v0.0.0
)

require (
	base v0.0.0 // indirect
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

replace core => ../../chenile-framework/core

replace packager => ../../chenile-framework/packager

replace test => ../../chenile-framework/test

replace http => ../../chenile-framework/http

replace base => ../../chenile-framework/base

replace owiz => ../../chenile-framework/owiz
