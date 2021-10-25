#!/usr/bin/env sh
curl -X PUT http://localhost:6000/c.txt --data '{"data": "c.txt data new replaced content"}'
