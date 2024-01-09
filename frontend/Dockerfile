# Usa una imagen base de Apache
FROM golang:alpine

#Creamos un directorio de trabajo para almacenar los ficheros
WORKDIR /server/clases
WORKDIR /server

#Copia el archivo index.html al servidor
COPY main.go /server
COPY /clases /server/clases

#Instalamos las dependencias de la aplicación
RUN go mod init ejemplo && \
    go get github.com/gin-gonic/gin && \
    go get github.com/gin-contrib/cors && \
    go mod tidy

CMD ["go", "run", "main.go"]