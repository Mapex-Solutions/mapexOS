# List of services to run
SERVICES := mapexIam assets router http_gateway triggers events workflow mapexVault

# Default target
.DEFAULT_GOAL := help

# Help command
help:
	@echo "make run                          - Run all services in parallel"
	@echo "make run-service SERVICE=name     - Run a specific service (e.g., assets, mapexIam, workflow)"
	@echo "make stop                         - Kill all running Go processes (if running in background)"
	@echo "make build                        - Build all services"
	@echo "make build-service SERVICE=name   - Build a specific service"
	@echo "make test                         - Run all tests in packages/"
	@echo "make test-package PACKAGE=name    - Run tests for a specific package (e.g., contracts)"
	@echo "make test-service SERVICE=name    - Run unit tests for a specific service"
	@echo ""
	@echo "E2E Tests:"
	@echo "  make e2e all                    - Run all services E2E tests"
	@echo "  make e2e {service}              - Run all modules of a service (e.g., make e2e assets)"
	@echo "  make e2e {service} {module}     - Run specific module (e.g., make e2e assets assettemplates)"
	@echo "  make e2e journey                - Run all journey tests"
	@echo "  make e2e journey {name}         - Run specific journey (e.g., make e2e journey organization_isolation)"

# Run all services in parallel
run:
	@pids=""; \
	for service in $(SERVICES); do \
		echo "🚀 Starting $$service..."; \
		go run ./services/$$service/src/main.go & \
		pids="$$pids $$!"; \
		sleep 2; \
	done; \
	trap 'echo "🛑 Stopping..."; kill $$pids; wait $$pids; exit 0' INT TERM; \
	wait $$pids

# Run a specific service by name
run-service:
	@if [ -z "$(SERVICE)" ]; then \
		echo "❌ Please provide a SERVICE name: make run-service SERVICE=assets"; \
	else \
		echo "🚀 Running $(SERVICE)..."; \
		go run ./services/$(SERVICE)/src/main.go; \
	fi

# Stop all Go services
stop:
	@echo "🛑 Stopping all Go services..."
	lsof -tiTCP:5000-5010 | xargs -r kill -9

# Build all services
build:
	@echo "🔨 Building all services..."
	@all_pass=true; \
	for service in $(SERVICES); do \
		echo "  Building $$service..."; \
		go build ./services/$$service/... 2>&1 || all_pass=false; \
	done; \
	if [ "$$all_pass" = true ]; then \
		echo "✅ All services built successfully!"; \
	else \
		echo "❌ Some services failed to build!"; \
		exit 1; \
	fi

# Build a specific service
build-service:
	@if [ -z "$(SERVICE)" ]; then \
		echo "❌ Please provide a SERVICE name: make build-service SERVICE=workflow"; \
	else \
		echo "🔨 Building $(SERVICE)..."; \
		go build ./services/$(SERVICE)/...; \
	fi

# Run unit tests for a specific service
test-service:
	@if [ -z "$(SERVICE)" ]; then \
		echo "❌ Please provide a SERVICE name: make test-service SERVICE=workflow"; \
	else \
		echo "🧪 Running unit tests for $(SERVICE)..."; \
		go test ./services/$(SERVICE)/... -count=1; \
	fi

# Run all tests in packages/
test:
	@echo "🧪 Running all tests..."
	go test ./packages/...

# Run tests for a specific package
test-package:
	@if [ -z "$(PACKAGE)" ]; then \
		echo "❌ Please provide a PACKAGE name: make test-package PACKAGE=infrastructure"; \
	else \
		echo "🧪 Running tests for $(PACKAGE)..."; \
		go test ./packages/$(PACKAGE)/...; \
	fi

# ========================================
# E2E Tests
# ========================================

# E2E command router
e2e:
	@$(MAKE) -s e2e-router ARGS="$(filter-out $@,$(MAKECMDGOALS))"

# Suppress "Nothing to be done" messages for args
%:
	@:

# E2E router logic
e2e-router:
	@if [ -z "$(ARGS)" ]; then \
		echo "❌ Usage:"; \
		echo "  make e2e all                    - Run all services E2E tests"; \
		echo "  make e2e {service}              - Run all modules of a service"; \
		echo "  make e2e {service} {module}     - Run specific module"; \
		echo "  make e2e journey                - Run all journeys"; \
		echo "  make e2e journey {name}         - Run specific journey"; \
		echo ""; \
		echo "Examples:"; \
		echo "  make e2e all                    - Run all services tests"; \
		echo "  make e2e assets                 - Run all asset modules"; \
		echo "  make e2e assets assettemplates  - Run only assettemplates module"; \
		echo "  make e2e mapexIam users         - Run only users module"; \
		echo "  make e2e journey                - Run all journeys"; \
		echo "  make e2e journey organization_isolation"; \
		exit 1; \
	fi; \
	\
	arg1=$$(echo "$(ARGS)" | awk '{print $$1}'); \
	arg2=$$(echo "$(ARGS)" | awk '{print $$2}'); \
	\
	if [ "$$arg1" = "all" ]; then \
		$(MAKE) -s e2e-all; \
	elif [ "$$arg1" = "journey" ]; then \
		if [ -z "$$arg2" ]; then \
			$(MAKE) -s e2e-journey-all; \
		else \
			$(MAKE) -s e2e-journey-single JOURNEY=$$arg2; \
		fi; \
	else \
		if [ -z "$$arg2" ]; then \
			$(MAKE) -s e2e-service SERVICE=$$arg1; \
		else \
			$(MAKE) -s e2e-module SERVICE=$$arg1 MODULE=$$arg2; \
		fi; \
	fi

# Run all services E2E tests
e2e-all:
	@echo "========================================="
	@echo "  Running All Services E2E Tests"
	@echo "========================================="
	@echo ""
	@all_services_pass=true; \
	for service in mapexIam assets router http_gateway triggers; do \
		echo ""; \
		echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"; \
		$(MAKE) -s e2e-service SERVICE=$$service || all_services_pass=false; \
	done; \
	echo ""; \
	echo "========================================="; \
	if [ "$$all_services_pass" = true ]; then \
		echo "✅ All services E2E tests passed!"; \
	else \
		echo "❌ Some services E2E tests failed!"; \
		exit 1; \
	fi

# Run all modules of a service
e2e-service:
	@echo "🧪 Running all E2E tests for service: $(SERVICE)"
	@echo ""
	@if [ "$(SERVICE)" = "mapexIam" ]; then \
		modules="auth groups lists memberships organizations roles users"; \
		e2e_dir="mapexos"; \
	elif [ "$(SERVICE)" = "assets" ]; then \
		modules="assets assettemplates"; \
		e2e_dir="assets"; \
	elif [ "$(SERVICE)" = "router" ]; then \
		modules="routegroups"; \
		e2e_dir="router"; \
	elif [ "$(SERVICE)" = "http_gateway" ]; then \
		modules="datasources"; \
		e2e_dir="http_gateway"; \
	elif [ "$(SERVICE)" = "triggers" ]; then \
		modules="triggers"; \
		e2e_dir="triggers"; \
	else \
		echo "❌ Unknown service: $(SERVICE)"; \
		echo "Available services: mapexIam, assets, router, http_gateway, triggers"; \
		exit 1; \
	fi; \
	\
	total_tests=0; \
	total_time=0; \
	all_pass=true; \
	echo "| Módulo            | Testes | Status      | Tempo  |"; \
	echo "|-------------------|--------|-------------|--------|"; \
	for module in $$modules; do \
		cd packages/e2eTests && \
		output=$$(go test ./services/$$e2e_dir/$$module/... -v -count=1 2>&1); \
		cd ../..; \
		pass_count=$$(echo "$$output" | grep -c "^--- PASS" || true); \
		fail_count=$$(echo "$$output" | grep -c "^--- FAIL" || true); \
		total_count=$$((pass_count + fail_count)); \
		time=$$(echo "$$output" | grep "^ok" | awk '{print $$NF}' | sed 's/s//'); \
		if [ -z "$$time" ]; then \
			time=$$(echo "$$output" | grep "^FAIL" | awk '{print $$NF}' | sed 's/s//'); \
		fi; \
		if [ -z "$$time" ]; then time="0.000"; fi; \
		total_tests=$$((total_tests + total_count)); \
		total_time=$$(echo "$$total_time + $$time" | bc); \
		if [ $$fail_count -gt 0 ]; then \
			status="❌ FAIL"; \
			all_pass=false; \
		else \
			status="✅ PASS"; \
		fi; \
		module_name=$$(echo $$module | sed 's/.*/\u&/' | head -c 17); \
		printf "| %-17s | %-6s | %-11s | %6ss |\n" "$$module_name" "$$pass_count/$$total_count" "$$status" "$$time"; \
	done; \
	echo "|-------------------|--------|-------------|--------|"; \
	if [ "$$all_pass" = true ]; then \
		printf "| %-17s | %-6s | %-11s | %6ss |\n" "TOTAL" "$$total_tests/$$total_tests" "✅ PASS" "$$total_time"; \
		echo ""; \
		echo "✅ All $(SERVICE) tests passed!"; \
	else \
		printf "| %-17s | %-6s | %-11s | %6ss |\n" "TOTAL" "$$total_tests/$$total_tests" "❌ FAIL" "$$total_time"; \
		echo ""; \
		echo "❌ Some $(SERVICE) tests failed!"; \
		exit 1; \
	fi

# Run a specific module
e2e-module:
	@echo "🧪 Running E2E tests for: $(SERVICE)/$(MODULE)"
	@echo ""
	@cd packages/e2eTests && \
	output=$$(go test ./services/$(SERVICE)/$(MODULE)/... -v -count=1 2>&1); \
	cd ../..; \
	echo "$$output"; \
	echo ""; \
	pass_count=$$(echo "$$output" | grep -c "^--- PASS" || true); \
	fail_count=$$(echo "$$output" | grep -c "^--- FAIL" || true); \
	total_count=$$((pass_count + fail_count)); \
	time=$$(echo "$$output" | grep "^ok" | awk '{print $$NF}' | sed 's/s//'); \
	if [ -z "$$time" ]; then \
		time=$$(echo "$$output" | grep "^FAIL" | awk '{print $$NF}' | sed 's/s//'); \
	fi; \
	if [ -z "$$time" ]; then time="0.000"; fi; \
	echo "| Módulo            | Testes       | Status      | Tempo  |"; \
	echo "|-------------------|--------------|-------------|--------|"; \
	if [ $$fail_count -gt 0 ]; then \
		printf "| %-17s | %-12s | %-11s | %6ss |\n" "$(MODULE)" "$$pass_count/$$total_count" "❌ FAIL" "$$time"; \
		echo ""; \
		echo "❌ Test failed!"; \
		exit 1; \
	else \
		printf "| %-17s | %-12s | %-11s | %6ss |\n" "$(MODULE)" "$$pass_count/$$total_count" "✅ PASS" "$$time"; \
		echo ""; \
		echo "✅ Test passed!"; \
	fi

# Run all journey tests
e2e-journey-all:
	@echo "🚀 Running all Journey Tests..."
	@echo ""
	@echo "| Journey Test                | Status      | Tempo  |"
	@echo "|-----------------------------|-------------|--------|"
	@all_pass=true; \
	total_time=0; \
	journeys=$$(find packages/e2eTests/journey -mindepth 1 -maxdepth 1 -type d -exec basename {} \;); \
	for journey in $$journeys; do \
		cd packages/e2eTests && \
		output=$$(MAPEXOS_ROOT=$$(cd ../../.. && pwd) go test ./journey/$$journey/... -v -count=1 -timeout=30m 2>&1); \
		cd ../..; \
		fail_count=$$(echo "$$output" | grep -c "^FAIL" || true); \
		time=$$(echo "$$output" | grep "^ok" | awk '{print $$NF}' | sed 's/s//'); \
		if [ -z "$$time" ]; then \
			time=$$(echo "$$output" | grep "^FAIL" | awk '{print $$NF}' | sed 's/s//'); \
		fi; \
		if [ -z "$$time" ]; then time="0.000"; fi; \
		total_time=$$(echo "$$total_time + $$time" | bc); \
		if [ $$fail_count -gt 0 ]; then \
			status="❌ FAIL"; \
			all_pass=false; \
		else \
			status="✅ PASS"; \
		fi; \
		journey_name=$$(echo $$journey | tr '_' ' ' | sed 's/.*/\u&/' | head -c 27); \
		printf "| %-27s | %-11s | %6ss |\n" "$$journey_name" "$$status" "$$time"; \
	done; \
	echo "|-----------------------------|-------------|--------|"; \
	echo ""; \
	if [ "$$all_pass" = true ]; then \
		echo "✅ All journey tests passed!"; \
	else \
		echo "❌ Some journey tests failed!"; \
		exit 1; \
	fi

# Run a specific journey test
e2e-journey-single:
	@echo "🚀 Running Journey Test: $(JOURNEY)"
	@echo ""
	@cd packages/e2eTests && \
	output=$$(MAPEXOS_ROOT=$$(cd ../../.. && pwd) go test ./journey/$(JOURNEY)/... -v -count=1 -timeout=30m 2>&1); \
	cd ../..; \
	echo "$$output"; \
	echo ""; \
	fail_count=$$(echo "$$output" | grep -c "^FAIL" || true); \
	time=$$(echo "$$output" | grep "^ok" | awk '{print $$NF}' | sed 's/s//'); \
	if [ -z "$$time" ]; then \
		time=$$(echo "$$output" | grep "^FAIL" | awk '{print $$NF}' | sed 's/s//'); \
	fi; \
	if [ -z "$$time" ]; then time="0.000"; fi; \
	echo "| Journey Test                | Status      | Tempo  |"; \
	echo "|-----------------------------|-------------|--------|"; \
	if [ $$fail_count -gt 0 ]; then \
		printf "| %-27s | %-11s | %6ss |\n" "$(JOURNEY)" "❌ FAIL" "$$time"; \
		echo ""; \
		echo "❌ Journey test failed!"; \
		exit 1; \
	else \
		printf "| %-27s | %-11s | %6ss |\n" "$(JOURNEY)" "✅ PASS" "$$time"; \
		echo ""; \
		echo "✅ Journey test passed!"; \
	fi
