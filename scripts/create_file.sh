#!/usr/bin/env sh
curl -X POST http://localhost:6000/c.txt --data '{"type": "file", "data": "c.txt data"}'
