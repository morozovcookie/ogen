#!/bin/sh
#
# Change imports from github.com/ogen-go/ogen to github.com/morozovcookie/ogen

MODULE_NAME='github.com/ogen-go/ogen'
SED_EXPR='s/github.com\/ogen-go\/ogen/github.com\/morozovcookie\/ogen/g'

grep -lrnw '.' --exclude "./update_module_name.sh" -e ${MODULE_NAME} | \
    xargs -n 1 sed -i '' ${SED_EXPR}

make tidy
make lint
make test
