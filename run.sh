#!/usr/bin/env bash

mkdir -p albums
docker -v ./albums:/albums -it run hsinhoyeh/xuitecrawler /out/xuitecrawler --album="<album-url>" --password=<password> --destination=/albums
