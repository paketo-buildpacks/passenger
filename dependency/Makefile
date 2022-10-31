.PHONY: retrieve test

retrieve:
	@cd retrieval; \
	go run . \
		--buildpack_toml_path=$(buildpackTomlPath) \
		--output=$(output)

test:
	@cd test; \
	./test.sh \
		--tarballPath $(tarballPath) \
		--expectedVersion $(version)