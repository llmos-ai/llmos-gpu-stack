#!/bin/bash
set -e

exec dumb-init -- llmos-gpu-stack device-manager "${@}"
