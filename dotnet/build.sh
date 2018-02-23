#!/usr/bin/env bash

set -e

dotnet build Metaparticle.Package

# Build examples
cd examples
for D in *; do
    if [ -d "${D}" ]; then
        dotnet build ${D}
    fi
done
