FROM debian-base

RUN useradd -u 1001 -m -s /bin/bash client && echo "client:client" | chpasswd

RUN mkdir app/

COPY worker app/

WORKDIR /app

CMD ["./worker"]

