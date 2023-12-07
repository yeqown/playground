#!/bin/sh

BASE=$PWD/base
mkdir -p $BASE
OVERLAY=$PWD/overlays
mkdir -p $OVERLAY

# fetch base resources
for name in {configMap,deployment,kustomization,service}; do
  curl -sSL -o $BASE/$name.yaml https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/examples/helloWorld/$name.yaml
done

overlays="staging prod"
for overlay in $overlays; do
  mkdir -p $OVERLAY/$overlay
done