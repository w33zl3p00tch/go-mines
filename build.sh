#!/usr/bin/env bash

# This builds a binary release for various architectures.
# Most users of other systems probably want to compile this themselves.
#
# run this from the repo's root directory like so: 
# $ ./build.sh

progName="go-mines_"
version="v0.0.1_pre-alpha_"
prefix="$progName$version"

declare -a arch=(
	"Linux_amd64 linux amd64"
	"Linux_x86 linux 386"
	"Windows_amd64 windows amd64"
	#"macOS darwin amd64"
	#"FreeBSD_amd64 freebsd amd64"
	#"FreeBSD_x86 freebsd 386"
	"Linux_arm32 linux arm"
	"Linux_arm64 linux arm64"
)

for i in "${arch[@]}"
do
	name=`echo $i | cut -d " " -f 1`
	goos=`echo $i | cut -d " " -f 2`
	goarch=`echo $i | cut -d " " -f 3`

	mkdir -p build/$prefix$name

	pushd build/$prefix$name
		echo "Building $prefix$name"
		GOOS=$goos GOARCH=$goarch go build -ldflags "-s -w" ../../
		cp ../../LICENSE .
		cp ../../README.md .
		cd ..
		tar -zcvf $prefix$name.tar.gz $prefix$name
		shaSum=`sha256sum $prefix$name.tar.gz`
		echo $shaSum >> sha256sums.txt
	popd
	echo -e "\n\n"
done

