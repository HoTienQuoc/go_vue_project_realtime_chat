.PHONY: clear-swag swag-init

clear-swag:
	@echo "Clearing docs directory"
	rm -rf docs/* || true

swag-init: clear-swag
	@echo "Generate swagger docs"
	swag init -g main.go
