#################INSTRUCTIONS#################
# Launch all containers using Makefile
# The script launches 50 requests in order to congest the system

#!/bin/bash

counter=0

while [ $counter -lt $1 ]; do

    json_file="/tmp/request_body.json"
    cookie_file="/tmp/cookie.txt"

    echo '{"url": "https://github.com/vlopser/test1.git", "parameters": ["1"]}' > "$json_file"

    echo "_oauth2_proxy=6w3Jd174V77505d7R1gt33b95S4F-uCFWc_gz1edNDJwKT0_4IlnJoFulQ3aKTX3rd1SsZhysR1lO7xb3r1YZ1UbQIBKN6PVy1xelKIAf1RWaYPiOVfE2sdDQ6OgJP9jgvTEELtjZ4nbh7GYvyzmIPXTgbEgCG9MajgaB5obdY8zuhoBUIzhrP5j3b8DKQ0ixHdRa6Zs0U-gZIir22AqBZVsWKgqfd6Yhr5gvuzFKK280AQXv6d_hG0C9t8TDo6gaH6PcFlU4fRZdo2ms4i-DrtP9YagCG5yCJujahTLJ1gwBL1o2rPvBen_HMwBCOp5cxiVHk3egf2tu_1n3plKfDBjxZirZmDU3j84ss6m-yiu3XCl7VxNk3ZVSwWIYA5H8o9yZ39E1Q1AcTh3GF3Z0aC3uYKQjXlaotlwexdWIfwyk3qGuonVqDEKVyEmP0Nq3Z618gAa9ivNFiXzlsARSRQB3uuC4So9tuUwaEWUI74Lo76Y8Y_JiSLmFcWaK64gl9qCNiR4zdjaUOwzJkjTsrCNDoM5cs4LbOc0SsXpaJ2m0EKAJY3rXQ3qWDui8kmrSW_mq3Fq9Y3sFSArMi5njmXR9a4BzObxA6qeAfXT3Rv-JOu90C9ivv6i1SpUWuqL1eySNSsUG0-8dcHjeIFpzgEitkcxfRRcUdp64u6wdTwSII9Ae257ZHlJkKfgMRU8XdOpJOhE65Q6zJqQUu953RIe7aFHQHwT4kLSEblhGkyjAkvvsxXT__ApnvP0Zys2x04ChWubhDPzKqrkF2rTlwb-5P_MjUFWn7WJK9XaMCylzQy5XIsqVesBulCeD_-0Cs4FXgAVncny2a4AkeoWmGEi1VC1hvDSrjUg2F_V55VAiSPa6tdZ2WFG-5F8EyBYYxTEYCx9fak9Ru5gapVr17ujEzljmxOS5AFSqjjmci880owfZMT3b1Mfq2DnJR9dbBPyHvODSgBiRPiQqBe-YlXSlcNWiIq5czXlyaq-PytWLZFWAQVjiLnU44qsqwCyRLzpuVFkXW6PCPZ_06p0wdO_-lVoYBhB-ZMQWmk-YvwFwPuXsXp6SG__zZdIRBiUP0FqIU5vvYtDQN14H2wBQxfDP1exhDIcSqRRyd-JsSJALOBmF50IpIm7PWaxCRlqcOE_ZFqsWjZp4cD5ZcsRgBWTIuEyxyR9iNpgGY4CPz8pFpYMmuoJMmtV58SkK6UyUefTCwOEErDnHm1LN4phgONzdAltwFwdzZkS9lkOdZuL8tECyADk3R58iA8PUpXeA8LOPP7FQ9rXuC5yzP5FzcGFis5nhUYdyuX3KgPwNnAfwop9O-oMcKgi0o3_hd6IUy0yxUJL13wCXW8C1d2QpS-c-UJzeh_N-nZAcwIeS8RbKSredSsbo9JRuy4hPL3NuuXO5Zjjb9p4G2QQSIQpH288ijUCc0ru6ivF0Tsgf2XbtKGQqR7dQR8Po2eRbUk_bK4H3bdS8NB7Mkwg4gdKRPXMU-PGrjJ4AGgRpbXJUnWPqcr9b9xscCpb4A3vNmUMueDyyYQ-3zc3DZCshwQhFDOsjF6dAjpn16aI3uw-Lc5pXl4sRIZCNnaUqkNvVMChBltnqH9mc3Og_NUmb-QhJ1w7Y6HoLeBBn2GLv4jlhZp9TxeRaG7DSAnBNdcQuff1wwVhp7ly9-qyftxI7zcycM0HAhF8G3HhtNR5y2LOXbiuasKI4IUcdXxnpNoNlw==|1705932718|KW4UhW_8BSUmERjoL3VazcecnhSx8QIMQzVr04Y9j48=; Path=/; Expires=Mon, 29 Jan 2024 14:11:58 GMT;" > "$cookie_file"
    curl -X POST \
         -H "Content-Type: application/json" \
         -H "Cookie: $(cat $cookie_file)" \
         -d "@$json_file" \
         http://localhost:4180/createTask

    ((counter++))

done