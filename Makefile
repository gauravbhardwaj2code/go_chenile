FRAMEWORK_MODULES=./base/... ./owiz/... ./core/... ./http/... ./test/... ./servicegen/... ./packager/...
EXAMPLE_MODULES=../chenile-examples/customer-service/... ../chenile-examples/order-service/... ../chenile-examples/mainweb-app/...
COVER_MODULES=./base/... ./owiz/... ./core/... ./http/... ./test/... ./servicegen/... ./packager/... ../chenile-examples/customer-service/customer/... ../chenile-examples/customer-service/test/... ../chenile-examples/order-service/order/... ../chenile-examples/order-service/test/... ../chenile-examples/mainweb-app/...

.PHONY: test
test:
	@echo "Testing framework modules..."
	cd chenile-framework && go test $(FRAMEWORK_MODULES)
	@echo "Testing example services..."
	cd chenile-examples && go test ./customer-service/... ./order-service/... ./mainweb-app/...

.PHONY: coverage
coverage:
	cd chenile-framework && go test -cover $(COVER_MODULES)

.PHONY: generate-example
generate-example:
	go run ./chenile-framework/servicegen/cmd/chenile-servicegen new --name customer --out ./chenile-examples --framework-root ../..
