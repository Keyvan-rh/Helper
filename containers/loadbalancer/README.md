# How to test this image

Build it

```shell
podman build -t test-lb .
```

Run it

```shell
podman run \ 
-e=HELPERPOD_CONFIG_YAML='aGVscGVycG9kOgogIGlzOiBjb29sCmhlbHBlcm5vZGU6CiAgdXNlaGVscDoKICAtIG5hbWU6IHJvYmVydAogICAgbGFzdDogc2FuZG92b2wK' \
--name=test-lb -id test-lb
```

Check the yaml

```shell
podman exec -it test-lb cat /usr/local/src/helperpod.yaml
```
