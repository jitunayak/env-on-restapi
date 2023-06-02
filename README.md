[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/5113309-5bf23fb5-e054-4f9b-a697-2b80a861ef66?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D5113309-5bf23fb5-e054-4f9b-a697-2b80a861ef66%26entityType%3Dcollection%26workspaceId%3D8dc0e5fa-ee53-4e23-9699-34532bd6a9d7)

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
### Start the Server

```bash
eli -server
```
### Only Run A Cron Job from command line

```bash
eli --cron --interval 10 --cmd 'echo jitu'
```

>eli --help (shows all the available command line arguments)

