Lanza un servicio de [OAuth Proxy](https://hub.docker.com/r/bitnami/oauth2-proxy) en `localhost:4180` que reenviará el tráfico a un [Frontend de prueba ](/frontend_mocked/) en `localhost:80` no accesible directamente desde internet.

# Ejecución
Para la ejecución únicamente situese en el directorio principal y ejecute

```shell
docker-compose up
```
Automáticamente creará los dos contenedores Docker necesarios.

Puede ejecutar también

```shell
docker-compose down
```
para únicamente eliminar todo rastro de los contenedores creados.

# Utilización
Para usar el servicio, dirigase a `localhost:4180` e inicie sesión con Google. Automáticamente será redirigido al [Frontend de prueba ](/frontend_mocked/).

# Referencias
Visto en https://dev.to/lazypro/make-any-website-authenticated-a52.