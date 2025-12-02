.PHONY: fmt

fmt:
	@echo "🎨  Formatting code with goimports..."
	@if ! command -v goimports > /dev/null; then \
		echo "🔍  goimports not found, installing..."; \
		go install golang.org/x/tools/cmd/goimports@latest; \
	fi
	@goimports -w .
	@echo "✨  Done!"
