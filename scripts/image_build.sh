#!/bin/bash

label=$(git rev-parse --short HEAD)

cd ../App
for tag in $(podman images --format "{{.Tag}}" dservice); do
    if [[ "$tag" != "latest" ]]; then
        echo "Deleting old tag for dservice: $tag"
        podman rmi "dservice:$tag"
    fi
done
podman build -t dservice:$label .
podman tag dservice:$label dservice:latest
podman image prune -f

cd ../Monitor
for tag in $(podman images --format "{{.Tag}}" monitor); do
    if [[ "$tag" != "latest" ]]; then
        echo "Deleting old tag for monitor: $tag"
        podman rmi "monitor:$tag"
    fi
done
podman build -t monitor:$label .
podman tag monitor:$label monitor:latest
podman image prune -f

cd ../podman
podman-compose up -d