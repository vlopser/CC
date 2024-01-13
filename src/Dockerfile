# Usa una imagen base de Apache
FROM debian
RUN apt-get update && \
    apt-get install -y \
    golang

RUN mkdir cc

RUN mkdir cc/src

COPY pkg /cc/src/

COPY go.mod go.sum /cc/src/

RUN mkdir /cc/src/services
COPY services/frontend /cc/src/

# RUN mkdir /frontend /frontend/classes /frontend/natsUtils /frontend/service

#Creamos un directorio de trabajo para almacenar los ficheros
WORKDIR /cc/src/services/frontend

#Copia el archivo index.html al servidor
# COPY main.go hola.txt ./
# COPY classes/* classes
# COPY service/* service

# #Instalamos las dependencias de la aplicación
# RUN go mod init frontend && \
#     go get github.com/gin-gonic/gin && \
#     go get github.com/gin-contrib/cors && \
#     go mod tidy

# CMD go run /frontend/main.go