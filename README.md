![](https://github.com/jitunayak/env-on-restapi/releases/download/latest/snap1.png)

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

> sample response

```json
{
    "accessKeyId":"HHDUUIEKKED",
    "secretKey":"Ng4W//WAttR33ugTroNSBQrbsdsdd7PR7QH7O",
    "sessionToken":"//sddd"
}
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

### Run a cron task example

```bash
curl --location 'http://localhost:8088/aws?reAuthenticate=true&interval=5&command=mkdir%newFolder'
```

### Only Run A Cron Job from command line

```bash
env-server --server --cron --interval 10 --cmd 'echo jitu'
```
