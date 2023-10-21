#!/usr/bin/env bash

set -ex

npm install && npm run build && xdg-open index.html
