# Usa una imagen base de Apache
FROM golang:alpine

RUN mkdir /frontend /frontend/classes /frontend/natsUtils /frontend/service

#Creamos un directorio de trabajo para almacenar los ficheros
WORKDIR /frontend

#Copia el archivo index.html al servidor
COPY main.go ./
COPY classes/* classes
COPY service/* service

#Instalamos las dependencias de la aplicación
RUN go mod init frontend && \
    go get github.com/gin-gonic/gin && \
    go get github.com/gin-contrib/cors && \
    go mod tidy

CMD go run /frontend/main.go