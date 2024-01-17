# Usa una imagen base de Apache
FROM debian
RUN apt-get update && \
    apt-get install -y \
    wget \
    sudo \
    vim

# Install go 1.21.6
RUN wget -O /tmp/go-1.21.6.tar.gz https://go.dev/dl/go1.21.6.linux-386.tar.gz
RUN sudo tar -C /usr/local -xzf /tmp/go-1.21.6.tar.gz
ENV PATH=$PATH:/usr/local/go/bin

SHELL ["/bin/bash", "-c"]
RUN source ~/.bashrc


# RUN mkdir /cc \
#           /cc/src \
#           /cc/src/pkg \
#           /cc/src/services

# COPY go.mod go.sum /cc/src/

# # COPY creará el directorio src y src/pkg si no existe
# COPY pkg /cc/src/pkg 