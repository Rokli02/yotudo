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

case $MODE in
    "dev")
        create_data_dir_if_not_found .

        echo Launching DEV application
        sudo wails dev -tags webkit2_41
        ;;
    "build")
        echo Building production application

        sudo chmod -R 777 frontend/dist
        sudo chmod -R 777 frontend/wailsjs/runtime

        wails build -tags webkit2_41

        echo Copying dynamic assets to build directory

        sudo chmod 777 ./build/bin

        create_data_dir_if_not_found ./build/bin
        ;;
    "test")
        echo Running tests

        output_of_test=$(go test -v ./test/**/*)
        echo "$output_of_test" | GREP_COLORS='mt=1;31' grep --color=always -e "--- FAIL"
        echo "$output_of_test" | GREP_COLORS='mt=1;32' grep --color=always -e "--- PASS"
        ;;
    *)
        echo Unknown command, use \"dev\", \"build\" or \"test\" instead!
        ;;
esac