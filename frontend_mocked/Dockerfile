# Usa una imagen base de Apache
FROM httpd:latest

# Copia archivos locales al contenedor (opcional)
# COPY ./mi-sitio-web /usr/local/apache2/htdocs/

# El contenedor expone automáticamente el puerto 80 para HTTP
EXPOSE 80

CMD ["httpd-foreground"]