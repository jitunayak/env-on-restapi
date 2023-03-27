# Begin
echo STEP:1 docker build started..

# Build
docker build . --tag jitu715/env-on-restapi

echo docker build completed

# Publish

docker push jitu715/env-on-restapi:latest
echo STEP:2 docker build published to docker hub
echo STEP:3 you may use [docker pull jitu715/env-on-restapi:latest]

echo ------------------------------------
echo
echo STEP:4 build windows binary
export GOOS=windows
export GOARCH=386

go build -o .build/windows/env-on-restapi-386.exe
echo STEP:5 windows binary is located at .build/windows/env-on-restapi-386.exe

echo ------------------------------------
echo
echo STEP:6 build darwin binary
export GOOS=darwin
export GOARCH=arm64
go build -o .build/mac/env-on-restapi-arm64
echo STEP:7 darwin binary is located at .build/mac/env-on-restapi-arm64

echo ------------------------------------
echo
echo 🌟 latest code has been published to Docker 🐳
echo 🌟 binaries artifacts build sucessfully 💻
echo
echo 🌟 DONE 🙌