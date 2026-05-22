#!/bin/bash
# Bash completion for run-tests.sh
#
# To enable this, add to your ~/.bashrc or ~/.zshrc:
#   source /path/to/e2eTests/.run-tests-completion.bash

_run_tests_completion() {
    local cur prev opts services mapexos_modules assets_modules router_modules gateway_modules

    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    # Top-level commands
    opts="all check list help -q --quiet -p --parallel -t --timeout -h --help"

    # Services
    services="mapexos assets router http_gateway"

    # Modules per service
    mapexos_modules="organizations roles groups users memberships"
    assets_modules="assets assettemplates"
    router_modules="routegroups"
    gateway_modules="datasources"

    # If we're completing after a service name, suggest modules
    case "${prev}" in
        mapexos)
            COMPREPLY=( $(compgen -W "${mapexos_modules} -q --quiet -p -t" -- ${cur}) )
            return 0
            ;;
        assets)
            COMPREPLY=( $(compgen -W "${assets_modules} -q --quiet -p -t" -- ${cur}) )
            return 0
            ;;
        router)
            COMPREPLY=( $(compgen -W "${router_modules} -q --quiet -p -t" -- ${cur}) )
            return 0
            ;;
        http_gateway)
            COMPREPLY=( $(compgen -W "${gateway_modules} -q --quiet -p -t" -- ${cur}) )
            return 0
            ;;
        -p|--parallel)
            # Suggest common parallel values
            COMPREPLY=( $(compgen -W "1 2 4 8" -- ${cur}) )
            return 0
            ;;
        -t|--timeout)
            # Suggest common timeout values
            COMPREPLY=( $(compgen -W "1m 5m 10m 30m" -- ${cur}) )
            return 0
            ;;
    esac

    # Default completion - suggest services and commands
    COMPREPLY=( $(compgen -W "${opts} ${services}" -- ${cur}) )
    return 0
}

# Register the completion function
complete -F _run_tests_completion run-tests.sh
complete -F _run_tests_completion ./run-tests.sh
