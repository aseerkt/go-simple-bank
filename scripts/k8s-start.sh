#!/bin/sh

minikube start

minikube image build . -t aseerkt/simplebank:latest

helm install simplebank helm/simplebank