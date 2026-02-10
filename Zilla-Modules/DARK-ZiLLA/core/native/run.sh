#!/bin/bash

# Ensure the build directory exists and build using CMake
echo "Building CodezillA using CMake..."
cmake -B build
if [ $? -ne 0 ]; then
    echo "CMake configuration failed!"
    exit 1
fi

cmake --build build
if [ $? -ne 0 ]; then
    echo "CMake build failed!"
    exit 1
fi

# Run the executable
echo "Running CodezillA..."
./build/codezilla