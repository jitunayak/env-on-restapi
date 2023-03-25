## Instruction for docker

docker pull jitu715/env-on-restapi

## Build your own

docker build . --tag <username>/<image-name>
docker push <username>/<image-name>:latest

## usage

### get aws assume role credentials

```bash
curl --location 'http://localhost:8088/aws' \
--header 'Content-Type: application/json' \
'
```

### generic api to get any ENV variable value

```bash
curl --location 'http://localhost:8088/' \
--header 'Content-Type: application/json' \
--data '{
"homeBrewPrefix":"HOMEBREW_PREFIX"
}
'
```
