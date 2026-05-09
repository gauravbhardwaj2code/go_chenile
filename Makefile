GO_MODULES=./base/... ./owiz/... ./core/... ./http/... ./test/... ./servicegen/... ./packager/... ./examples/customer-service/... ./examples/order-service/... ./examples/mainweb-app/...
COVER_MODULES=./base/... ./owiz/... ./core/... ./http/... ./test/... ./servicegen/... ./packager/... ./examples/customer-service/customer/... ./examples/customer-service/test/... ./examples/order-service/order/... ./examples/order-service/test/... ./examples/mainweb-app/...

.PHONY: test
test:
	go test $(GO_MODULES)

.PHONY: coverage
coverage:
	go test -cover $(COVER_MODULES)

.PHONY: generate-example
generate-example:
	go run ./servicegen/cmd/chenile-servicegen new --name customer --out ./examples --framework-root ../..
