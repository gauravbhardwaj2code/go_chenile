module inventory-service

go 1.26

toolchain go1.26.3

require (
	github.com/gauravbhardwaj2code/go_chenile/chenile-framework/base v0.1.1
	github.com/gauravbhardwaj2code/go_chenile/chenile-framework/bdd-utils v0.1.0
	github.com/gauravbhardwaj2code/go_chenile/chenile-framework/chenile v0.1.1
	github.com/gauravbhardwaj2code/go_chenile/chenile-framework/config v0.1.0
	github.com/gauravbhardwaj2code/go_chenile/chenile-framework/packager v0.1.0
)

require (
	github.com/cucumber/gherkin/go/v26 v26.2.0 // indirect
	github.com/cucumber/godog v0.15.1 // indirect
	github.com/cucumber/messages/go/v21 v21.0.1 // indirect
	github.com/gauravbhardwaj2code/go_chenile/chenile-framework/core v0.1.1 // indirect
	github.com/gauravbhardwaj2code/go_chenile/chenile-framework/http v0.1.1 // indirect
	github.com/gauravbhardwaj2code/go_chenile/chenile-framework/owiz v0.1.1 // indirect
	github.com/gofrs/uuid v4.3.1+incompatible // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-memdb v1.3.4 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/spf13/pflag v1.0.7 // indirect
)

replace github.com/gauravbhardwaj2code/go_chenile/chenile-framework/bdd-utils => ../../chenile-framework/bdd-utils

replace github.com/gauravbhardwaj2code/go_chenile/chenile-framework/base => ../../chenile-framework/base

replace github.com/gauravbhardwaj2code/go_chenile/chenile-framework/chenile => ../../chenile-framework/chenile

replace github.com/gauravbhardwaj2code/go_chenile/chenile-framework/config => ../../chenile-framework/config

replace github.com/gauravbhardwaj2code/go_chenile/chenile-framework/core => ../../chenile-framework/core

replace github.com/gauravbhardwaj2code/go_chenile/chenile-framework/http => ../../chenile-framework/http

replace github.com/gauravbhardwaj2code/go_chenile/chenile-framework/owiz => ../../chenile-framework/owiz

replace github.com/gauravbhardwaj2code/go_chenile/chenile-framework/packager => ../../chenile-framework/packager
