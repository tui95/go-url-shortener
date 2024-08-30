run:
	@go run main.go

watch:
	@go run github.com/air-verse/air; \
	echo "Watching...";

.PHONY: run watch
