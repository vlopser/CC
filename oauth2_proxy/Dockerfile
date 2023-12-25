# Usa una imagen base de Debian
FROM debian

# Actualiza el sistema e instala Go
RUN apt-get update && \
    apt-get install -y ca-certificates && \
    update-ca-certificates && \
    apt-get install -y golang

# Instala OAuth2 Proxy usando Go
RUN go install github.com/oauth2-proxy/oauth2-proxy/v7@latest

# Establece la variable de entorno GOPATH
ENV GOPATH=~/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

# COpia los ficheros necesarios para la puesta en marcha del Oauth2 Proxy
WORKDIR /oauth

COPY oauth2-proxy.cfg oauth2-proxy.cfg
COPY entrypoint.sh /entrypoint.sh

#EXPOSE por si solo no hace nada hasta que no se exponen los puertos con docker run -p. SOlo es util para documentar
EXPOSE 4180

ENTRYPOINT ["/bin/bash", "/entrypoint.sh"]
