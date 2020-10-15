# How to test this image

Build it

```shell
podman build -t test-lb .
```

Run it

```shell
podman run \ 
-e=HELPERPOD_CONFIG_YAML='aGVscGVybm9kZToKICB1c2VoZWxwOgogIC0gbmFtZTogcm9iZXJ0CiAgICBsYXN0OiBzYW5kb3ZvbAo=' \
--name=test-lb -id test-lb
```

Check the yaml

```shell
podman exec -it test-lb cat /usr/local/src/helperpod.yaml
```
