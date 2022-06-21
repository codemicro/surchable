#!/bin/bash

set -ex

cog -rU -I . --verbosity=1 internal/config/config.go internal/urls/urls.go crawler/urls.py crawler/config.py

black crawler/urls.py crawler/config.py
go fmt ./internal/config/ ./internal/urls/
