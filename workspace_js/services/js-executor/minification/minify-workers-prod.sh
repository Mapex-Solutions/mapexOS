#!/bin/sh

echo 'Start cleaning up node_module and dist'
find dist/ \( \
    -name "*.html" -o \
    -name "*.txt" -o \
    -name "license" -o \
    -name "tsconfig.*" -o \
    -name "*.ts" -o \
    -name "*.md" -o \
    -name "*.yml" -o \
    -name "*.js.map" \
    \) \
    -exec rm '{}' + >/dev/null 2>&1

find node_modules/ \( \
    -name "*.html" -type f -o \
    -name "*.txt" -type f -o \
    -name "LICENSE" -type f -o \
    -name "license" -type f -o \
    -name "tsconfig.*" -type f -o \
    -name "example*" -type f -o \
    -name "*.ts" -type f -o \
    -name "*.map" -type f -o \
    -name "*.md" -type f -o \
    -name "*.yml" -type f -o \
    -name "*.js.map" -type f -o \
    -name ".npmignore" -type f -o \
    -name ".jshintrc" -type f \
    \) \
    -exec rm -rf '{}' \; >/dev/null 2>&1

find node_modules/ \( \
    -name "test" -type d -o \
    -name "tests" -type d -o \
    -name "__tests__" -type d -o \
    -name "example" -type d -o \
    -name "examples" -type d -o \
    -name "docs" -type d -o \
    -name "docs" -type d -o \
    -name ".github" -type d -o \
    -name "esm" -type d -o \
    -name "esm5" -type d \
    \) \
    -exec rm -rf '{}' + >/dev/null 2>&1

rm -rf \
    node_modules/@swc \
    node_modules/yaml \
    node_modules/y18n \
    node_modules/yargs \
    node_modules/typescript \
    node_modules/tsconfig* \
    node_modules/rxjs/dist/bundles \
    node_modules/rxjs/dist/esm \
    node_modules/rxjs/dist/esm5 \
    node_modules/rxjs/dist/types \
    node_modules/@nestjs/cli \
    node_modules/http-request-mock \
    node_modules/.bin/tsc \
    node_modules/.bin/http-request-mock \
    node_modules/.bin/http-request-mock-cli \
    node_modules/.bin/tsserver \
    node_modules/.bin/nest >/dev/null 2>&1
echo 'Cleaning up node_module has finished.'
