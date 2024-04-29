#!/bin/bash

VERSION=$1
OS=("linux" "darwin" "freebsd" "windows")
ARCH=("amd64" "386" "arm64" "arm")
PROVIDER="vcd-object-storage-ext"

package="github.com/josajunior81/terraform-provider-vcd-object-storage-ext"
output_name='terraform-provider-'${PROVIDER}'_v'${VERSION}

rm ./build/*

for i in ${OS[@]}
do
  GOOS=$i
  if [[ $GOOS == "windows" ]]; then
    output_name+='.exe'
  fi
  for j in ${ARCH[@]}
  do
    if [[ $i == "darwin"  && ($j == "386" || $j == "arm")  ]]; then
      continue
    fi
    echo "Building $i: $j"
    GOARCH=$j 
    env GOOS=$GOOS GOARCH=$GOARCH go build -o ./build/${output_name} $package

    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
      exit 1
    fi
    cd build && zip terraform-provider-${PROVIDER}_${VERSION}_${GOOS}_${GOARCH}.zip terraform-provider-${PROVIDER}_v${VERSION}*
    rm terraform-provider-${PROVIDER}_v* 
    cd ..
  done
done

cd build
shasum -a 256 *.zip > terraform-provider-${PROVIDER}_${VERSION}_SHA256SUMS
# gpg -u keybase.io/$1 --detach-sign terraform-provider-${PROVIDER}_${VERSION}_SHA256SUMS
keybase pgp sign -b -d -i terraform-provider-${PROVIDER}_${VERSION}_SHA256SUMS -o terraform-provider-${PROVIDER}_${VERSION}_SHA256SUMS.sig

