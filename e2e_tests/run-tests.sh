#!/bin/bash

# E2E Tests Runner Script
# Usage: ./run-tests.sh [SERVICE] [MODULE] [options]

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Default values
VERBOSE="-v"
PARALLEL=1
TIMEOUT="5m"

# Function to print colored output
print_header() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_info() {
    echo -e "${CYAN}ℹ️  $1${NC}"
}

# Check if a port is listening
check_port() {
    local port=$1
    lsof -i:$port -sTCP:LISTEN > /dev/null 2>&1
    return $?
}

# Check if services are running
check_services() {
    print_header "Checking Services"

    local all_ok=true

    # Check mapexos (port 5000) - includes auth endpoints
    if check_port 5000; then
        print_success "mapexos (5000) is running"
    else
        print_error "mapexos (5000) is NOT running"
        all_ok=false
    fi

    # Check router (port 5003) - optional
    if check_port 5003; then
        print_success "router (5003) is running"
    else
        print_warning "router (5003) is NOT running (optional for some tests)"
    fi

    # Check assets (port 5002) - optional
    if check_port 5002; then
        print_success "assets (5002) is running"
    else
        print_warning "assets (5002) is NOT running (optional for some tests)"
    fi

    # Check http_gateway (port 5001) - optional
    if check_port 5001; then
        print_success "http_gateway (5001) is running"
    else
        print_warning "http_gateway (5001) is NOT running (optional for some tests)"
    fi

    echo ""

    if [ "$all_ok" = false ]; then
        print_error "Required services are not running!"
        echo ""
        echo "Start services with:"
        echo "  cd /home/thiago/Documents/Projects/MAPEX/mapexOS/workspace_go"
        echo "  make run"
        echo ""
        exit 1
    fi
}

# Run specific module tests
run_tests() {
    local test_path=$1
    local description=$2

    print_header "Running tests: $description"
    go test "./$test_path" $VERBOSE -timeout $TIMEOUT -parallel $PARALLEL
}

# List available services and modules
list_available() {
    print_header "Available Services and Modules"

    echo -e "${CYAN}Services:${NC}"
    echo "  • mapexos        - Main service (organizations, roles, groups, users, memberships)"
    echo "  • assets         - Asset management (assets, assettemplates)"
    echo "  • router         - Route management (routegroups)"
    echo "  • http_gateway   - HTTP Gateway (datasources)"
    echo ""

    echo -e "${CYAN}Mapexos Modules:${NC}"
    echo "  • organizations  - Organization CRUD and hierarchy tests"
    echo "  • roles          - Role CRUD and permission tests"
    echo "  • groups         - Group CRUD tests"
    echo "  • users          - User CRUD tests"
    echo "  • memberships    - Membership CRUD tests"
    echo ""

    echo -e "${CYAN}Assets Modules:${NC}"
    echo "  • assets         - Asset CRUD tests"
    echo "  • assettemplates - Asset template CRUD tests"
    echo ""

    echo -e "${CYAN}Router Modules:${NC}"
    echo "  • routegroups    - Route group CRUD tests"
    echo ""

    echo -e "${CYAN}HTTP Gateway Modules:${NC}"
    echo "  • datasources    - Data source CRUD tests"
    echo ""

    echo -e "${CYAN}Saga Trigger Journeys (./run-tests.sh saga trigger-<type>):${NC}"
    echo "  • trigger-http       - HTTP trigger (in-process HTTP sink)"
    echo "  • trigger-email      - Email trigger (in-process SMTP sink)"
    echo "  • trigger-slack      - Slack trigger (HTTP webhook → shared sink)"
    echo "  • trigger-teams      - Teams trigger (HTTP webhook → shared sink)"
    echo "  • trigger-websocket  - WebSocket trigger (in-process WS sink)"
    echo "  • trigger-mqtt       - MQTT trigger (in-process mochi-mqtt broker)"
    echo "  • trigger-nats       - NATS trigger (embedded nats-server)"
    echo "  • trigger-rabbitmq   - RabbitMQ trigger (testcontainers ephemeral)"
    echo ""
}

# Show help
show_help() {
    cat << EOF
${CYAN}E2E Tests Runner${NC}

${YELLOW}Usage:${NC}
    ./run-tests.sh [SERVICE] [MODULE] [options]
    ./run-tests.sh [COMMAND] [options]

${YELLOW}Format with SERVICE and MODULE:${NC}
    ./run-tests.sh SERVICE MODULE       Run specific module tests

    ${GREEN}Examples:${NC}
    ./run-tests.sh mapexos organizations    # Run organizations tests
    ./run-tests.sh mapexos roles            # Run roles tests
    ./run-tests.sh mapexos users            # Run users tests
    ./run-tests.sh assets assets            # Run assets tests
    ./run-tests.sh router routegroups       # Run routegroups tests

    ${GREEN}With options:${NC}
    ./run-tests.sh mapexos users -q         # Run quietly (no verbose)
    ./run-tests.sh mapexos users -p 4       # Run with 4 parallel workers

${YELLOW}Shortcut Commands:${NC}
    all                 Run ALL tests from all services
    check               Check if services are running
    list                List available services and modules
    saga                Run saga trigger journeys (see below)
    help, -h, --help    Show this help message

${YELLOW}Saga Journeys:${NC}
    ./run-tests.sh saga                       Run ALL trigger journeys (phase1 + phase2)
    ./run-tests.sh saga triggers              Same as above (alias)
    ./run-tests.sh saga trigger-http          Run HTTP trigger journeys only
    ./run-tests.sh saga trigger-email         Run Email trigger journeys only
    ./run-tests.sh saga trigger-slack         Run Slack trigger journeys only
    ./run-tests.sh saga trigger-teams         Run Teams trigger journeys only
    ./run-tests.sh saga trigger-mqtt          Run MQTT trigger journeys only
    ./run-tests.sh saga trigger-nats          Run NATS trigger journeys only
    ./run-tests.sh saga trigger-rabbitmq      Run RabbitMQ trigger journeys only (testcontainers pulls image first run)
    ./run-tests.sh saga trigger-websocket     Run WebSocket trigger journeys only

${YELLOW}Run All Tests from a Service:${NC}
    ./run-tests.sh mapexos              # Run all mapexos tests
    ./run-tests.sh assets               # Run all assets tests
    ./run-tests.sh router               # Run all router tests
    ./run-tests.sh http_gateway         # Run all http_gateway tests

${YELLOW}Options:${NC}
    -q, --quiet        Run without verbose output
    -p N               Run tests in parallel (N = number of parallel tests)
    -t TIME            Set timeout (default: 5m)

${YELLOW}Environment Variables:${NC}
    MAPEXOS_URL        mapexos service URL (default: http://localhost:5000)
    ROUTER_URL         router service URL (default: http://localhost:5003)
    ASSETS_URL         assets service URL (default: http://localhost:5002)
    GATEWAY_URL        gateway service URL (default: http://localhost:5001)

${YELLOW}Examples:${NC}
    ./run-tests.sh all                          # Run ALL tests
    ./run-tests.sh mapexos                      # Run all mapexos tests
    ./run-tests.sh mapexos organizations        # Run organizations tests
    ./run-tests.sh mapexos users -q             # Run users tests quietly
    ./run-tests.sh assets assets -p 4           # Run assets tests with 4 workers
    ./run-tests.sh check                        # Check if services are running
    ./run-tests.sh list                         # List available services/modules

EOF
}

# Parse arguments
COMMAND=""
SERVICE=""
MODULE=""

# First, check for options and commands
while [[ $# -gt 0 ]]; do
    case $1 in
        -q|--quiet)
            VERBOSE=""
            shift
            ;;
        -p|--parallel)
            PARALLEL="$2"
            shift 2
            ;;
        -t|--timeout)
            TIMEOUT="$2"
            shift 2
            ;;
        help|--help|-h)
            show_help
            exit 0
            ;;
        list)
            list_available
            exit 0
            ;;
        check)
            check_services
            exit 0
            ;;
        all)
            check_services
            print_header "Running ALL E2E Tests"
            go test ./... $VERBOSE -timeout $TIMEOUT -parallel $PARALLEL
            print_success "All tests completed!"
            exit 0
            ;;
        saga)
            shift
            # saga subcommand — runs trigger journey smokes under journey/automations/
            #   ./run-tests.sh saga                       runs all trigger_* journeys
            #   ./run-tests.sh saga triggers              same as above (alias)
            #   ./run-tests.sh saga trigger-<type>        runs one trigger type (http, email, slack, teams, mqtt, nats, rabbitmq, websocket)
            #
            # Forces -p 1 because phase1/phase2 share the in-process HTTP sink
            # port (11010); serial execution avoids "address already in use".
            saga_target="${1:-triggers}"
            if [[ $# -gt 0 ]] && [[ ! $1 =~ ^- ]]; then
                shift
            fi
            check_services
            case "$saga_target" in
                triggers)
                    print_header "Running ALL saga trigger journeys"
                    go test ./journey/automations/... $VERBOSE -count=1 -timeout $TIMEOUT -p 1
                    ;;
                trigger-*)
                    saga_type="${saga_target#trigger-}"
                    saga_pkg="./journey/automations/trigger_${saga_type}/..."
                    print_header "Running saga trigger journey: $saga_type"
                    go test "$saga_pkg" $VERBOSE -count=1 -timeout $TIMEOUT -p 1
                    ;;
                *)
                    print_error "Unknown saga target: $saga_target"
                    echo "  Expected: triggers OR trigger-{http,email,slack,teams,mqtt,nats,rabbitmq,websocket}"
                    exit 1
                    ;;
            esac
            print_success "Saga run completed!"
            exit 0
            ;;
        mapexos|assets|router|http_gateway)
            SERVICE=$1
            shift
            # Check if next argument is a module or an option
            if [[ $# -gt 0 ]] && [[ ! $1 =~ ^- ]]; then
                MODULE=$1
                shift
            fi
            ;;
        *)
            print_error "Unknown argument: $1"
            echo ""
            echo "Run './run-tests.sh help' for usage information"
            echo "Run './run-tests.sh list' to see available services and modules"
            exit 1
            ;;
    esac
done

# Execute based on SERVICE and MODULE
if [ -n "$SERVICE" ]; then
    check_services

    if [ -n "$MODULE" ]; then
        # Run specific module in service
        case $SERVICE in
            mapexos)
                run_tests "services/mapexos/$MODULE" "$SERVICE/$MODULE"
                ;;
            assets)
                run_tests "services/assets/$MODULE" "$SERVICE/$MODULE"
                ;;
            router)
                run_tests "services/router/$MODULE" "$SERVICE/$MODULE"
                ;;
            http_gateway)
                run_tests "services/http_gateway/$MODULE" "$SERVICE/$MODULE"
                ;;
        esac
    else
        # Run all tests from service
        print_info "No module specified, running all tests from $SERVICE"
        echo ""
        run_tests "services/$SERVICE/..." "$SERVICE (all modules)"
    fi

    print_success "Tests completed!"
    exit 0
fi

# Default: show help
show_help
