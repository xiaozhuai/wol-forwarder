# wol-forwarder

## How it works

It listens udp port `${WOL_ADDR}:${WOL_PORT}` (`0.0.0.0:1999` by default) and receive magic packet,
then broadcast to `${WOL_BADDR}:${WOL_BPORT}` (`255.255.255.255:9` by default)

## Docker

```shell
docker run -d \
    --name wol-forwarder \
    --restart=unless-stopped \
    --network=host \
    -e WOL_BADDR=192.168.1.255 \
    -p 1999:1999/udp \
    xiaozhuai/wol-forwarder:latest
```
