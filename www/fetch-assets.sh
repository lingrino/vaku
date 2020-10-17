#!/bin/bash

# Pulls the latest css from remote sources
curl https://raw.githubusercontent.com/kognise/water.css/HEAD/dist/dark.min.css -o ./www/assets/css/water/dark.min.css
curl https://raw.githubusercontent.com/kognise/water.css/HEAD/dist/dark.min.css.map -o ./www/assets/css/water/dark.min.css.map
