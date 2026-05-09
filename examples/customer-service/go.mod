module customer-service

go 1.22

require (
	github.com/ajapro/chenile-go/core v0.0.0
	github.com/ajapro/chenile-go/packager v0.0.0
	github.com/ajapro/chenile-go/test v0.0.0
)

require (
	github.com/ajapro/chenile-go/base v0.0.0 // indirect
	github.com/ajapro/chenile-go/http v0.0.0 // indirect
	github.com/ajapro/chenile-go/owiz v0.0.0 // indirect
	github.com/cucumber/gherkin/go/v26 v26.2.0 // indirect
	github.com/cucumber/godog v0.15.1 // indirect
	github.com/cucumber/messages/go/v21 v21.0.1 // indirect
	github.com/gofrs/uuid v4.3.1+incompatible // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-memdb v1.3.4 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/spf13/pflag v1.0.7 // indirect
)

replace github.com/ajapro/chenile-go/base => ../../base

replace github.com/ajapro/chenile-go/core => ../../core

replace github.com/ajapro/chenile-go/http => ../../http

replace github.com/ajapro/chenile-go/owiz => ../../owiz

replace github.com/ajapro/chenile-go/packager => ../../packager

replace github.com/ajapro/chenile-go/test => ../../test
