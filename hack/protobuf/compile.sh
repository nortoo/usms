#!/bin/bash

set -e  # Exit on error

function getFileSuffix() {
    local filename="$1"
    if [[ -n "$filename" ]]; then
        echo "${filename##*.}"
    fi
}

function compile() {
    local proto_dir="$1"
    if [[ ! -d "$proto_dir" ]]; then
        echo "Error: Directory $proto_dir does not exist"
        exit 1
    fi

    for file in $(ls "$proto_dir"); do
        if [[ -d "$proto_dir/$file" ]]; then
            if [[ $file != '.' && $file != '..' ]]; then
                compile "$proto_dir/$file"
            fi
        else
            if [[ "$(getFileSuffix "${file}")" == "proto" ]]; then
                echo "Compiling $proto_dir/$file..."
                protoc \
                    -I. \
                    -I/usr/local/include \
                    -Ipkg/proto \
                    --proto_path=pkg/proto \
                    --go_out=. \
                    --go_opt=paths=source_relative \
                    --go-grpc_out=. \
                    --go-grpc_opt=paths=source_relative \
                    "$proto_dir/$file"
                
                if [ $? -ne 0 ]; then
                    echo "Error: Failed to compile $proto_dir/$file"
                    exit 1
                fi
            fi
        fi
    done
}

# Ensure we're in the project root directory
cd "$(dirname "$0")/../.."

compile "$1"