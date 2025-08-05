#!/bin/bash

MODE=$1
WAILS_PATH=$HOME/go/bin

create_data_dir_if_not_found() {
    local base_dir=$1

    if [ ! -d "$base_dir/data" ]; then
        echo Making data folder in \"$base_dir\" directory

        mkdir -p $base_dir/data/{imgs,mscs,tmp}
    fi
}

if [ "$MODE" = "dev" ]; then
    create_data_dir_if_not_found .

    echo Launching DEV application
    sudo $WAILS_PATH/wails dev -tags webkit2_41

elif [ "$MODE" = "build" ]; then
    local build_dir="./build/bin"

    echo Building production application
    sudo $WAILS_PATH/wails build -tags webkit2_41

    echo Copying dynamic assets to build directory
    cp -r ./assets $build_dir/assets

    create_data_dir_if_not_found $build_dir

elif [ "$MODE" = "test" ]; then
    echo Running tests

    go test -v ./test/**/*
else
    echo Unknown command, use \"dev\", \"build\" or \"test\" instead!
fi