#!/bin/bash

cd /api
ag ./asyncapi.yaml @asyncapi/html-template -o /asyncapi-docs