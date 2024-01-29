#################INSTRUCTIONS#################
# Launch all containers using Makefile
# Launch the script passing the number of request and the oAuth token in order to congest the system

#!/bin/bash

counter=0

while [ $counter -lt $1 ]; do
    cookie="_oauth2_proxy=$2"
    json_body='{"url": "https://github.com/vlopser/test1.git", "parameters": ["1"]}'
    curl -X POST \
         -H "Content-Type: application/json" \
         -H "Cookie: $cookie" \
         -d "$json_body" \
         http://localhost:4180/createTask
    ((counter++))
done
