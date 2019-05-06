#!/bin/bash

group="assets"
program="assest_server"
program_version="v1.0"
compiler_version=`go version`
author=`whoami`
build_time=`date`

package() {
    if [ ! -e "./${program}" ]
    then
        build
        if [ ! -e "./${program}" ]
        then
            echo "build ${program} failed!!!"
            exit
        fi
    fi
    pkg_name="${program}_${program_version}.tar.gz"
    tar -zcf $pkg_name ${program} --exclude=log/*.log log etc
}

build() {
    go build -ldflags \
        "-X 'main.PROGRAM_VERSION=${program_version}' \
         -X 'main.COMPILER_VERSION=${compiler_version}' \
         -X 'main.BUILD_TIME=${build_time}' \
         -X 'main.AUTHOR=${author}'" \
        -o ${program}
    # RETVAL=$? && [ $RETVAL -ne 0 ] && exit 1
}

if [ "$1"x == "pkg"x ]
then
    package
else
    build
fi