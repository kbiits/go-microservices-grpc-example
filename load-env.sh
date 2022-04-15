#!/bin/bash

env_file="${1}"

function loadEnv() {
    local envFile="${1?Missing environment file}"
    local environmentAsArray variableDeclaration
    mapfile environmentAsArray < <(
        grep --invert-match '^#' "${envFile}" |
            grep --invert-match '^\s*$'
    ) # Uses grep to remove commented and blank lines
    for variableDeclaration in "${environmentAsArray[@]}"; do
        export "${variableDeclaration//[$'\n']/}" # The substitution removes the line breaks
    done
}

if [[ -z "${env_file}" ]]; then
    if [[ "${APP_MODE}" == "development" ]]; then
        env_file='.env.local'
    else
        env_file='.env'
    fi
fi

if [ -f "${env_file}" ]; then
    loadEnv "${env_file}"
else
    exit 1
fi

exit 0
