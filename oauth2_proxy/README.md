Crea un contenedor Docker a partir de un [Dockerfile](Dockerfile) que desplegará un servidor Debian con Oauth2-proxy y un [archivo de configuración](oauth2-proxy.cfg), para el inicio de sesión mediante Google a un [Frontend](../frontend_mocked/).

# Explicación Archivo de Configuración.
* `client-id`. ID de cliente proporcionado por Google.
* `client_secret`. Secreto del cliente proporcionado por Google.
* `cookie_secret`. Clave secreta utilizada para firmar las cookies. Generado con `openssl rand -base64 32 | tr -- '+/' '-_'`.
* `redirect_url`. URL a la que se redirige una vez autenticado. En este caso la dirección del [Frontend](../frontend_mocked/).
* `email_domains`. Lista de correos permitidos. En este caso todos.
* `http_address`. Dirección IP y puerto en el que se ejecuta Oauth2-proxy. En este caso `localhost:4180`.