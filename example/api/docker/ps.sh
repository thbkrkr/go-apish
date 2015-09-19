#!/bin/bash
set -eu

docker ps \
        --format '{"id":"{{.ID}}", "name": "{{.Names}}", 
            "created_at": "{{.CreatedAt}}", "status": "{{.Status}}", 
            "ports": "{{.Ports}}", "size":"{{.Size}}"}' \
        | jq -s .