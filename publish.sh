# Begin
echo STEP:1 docker build started..

# Build
docker build . --tag jitu715/env-on-restapi

echo docker build completed

# Publish

# docker push jitu715/env-on-restapi:latest
echo STEP:2 docker build published to docker hub
echo STEP:3 you may use [docker pull jitu715/env-on-restapi:latest]