.PHONY: test test-integration test-unit

test: test-unit test-integration

#test-unit:
	#go test ./test/test_unit/...

test-integration:
	go test ./test/test_integration/...
