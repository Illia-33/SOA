#!/bin/bash
set -e

docker buildx build \
	--load \
	--file $PWD/deploy/go-builder.Dockerfile \
	--build-context go-mods=$PWD \
	--tag soa-go-builder \
	$PWD
