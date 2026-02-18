.PHONY: clear-swag swag-init

clear-swag:
	@echo "Clearing docs directory"
	del /Q docs\* 2>null || exit 0


swag-init: clear-swag
	@echo "Generate swagger docs"
	swag init -g main.go
