#!/bin/sh

echo 'Start cleaning up node_module'
find node_modules/ \( \
    -name "*.html" -type f -o \
    -name "*.txt" -type f -o \
    -name "LICENSE" -type f -o \
    -name "license" -type f -o \
    -name "example*" -type f -o \
    -name "*.map" -type f -o \
    -name "*.md" -type f -o \
    -name "*.yml" -type f -o \
    -name "*.js.map" -type f -o \
    -name ".npmignore" -type f -o \
    -name ".jshintrc" -type f -o \
    -name "*.umd.js" -type f -o \
    -name "*.umd.min.js" -type f \
    \) \
    -exec rm '{}' + >/dev/null 2>&1
find node_modules/ \( \
    -name "test" -type d -o \
    -name "tests" -type d -o \
    -name "__tests__" -type d -o \
    -name "example" -type d -o \
    -name "examples" -type d -o \
    -name "docs" -type d -o \
    -name ".github" -type d -o \
    \) \
    -exec rm -rf '{}' + >/dev/null 2>&1
echo 'Cleaning up node_module has finished.'
