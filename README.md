Lnaza un servicio de [OAuth Proxy](/oauth2_proxy/) en `localhost:4180` que reenviará el tráfico a un [Frontend de prueba ](/frontend_mocked/)en `localhost:8080`.

# Ejecución
Para la ejecución únicamente situese en el directorio principal y ejecute

```shell
make
```
Automáticamente creará los dos contenedores Docker necesarios y eliminará cualquier rastro suyo previo.

Puede ejecutar también

```shell
make remove
```
para únicamente eliminar todo rastro de los contenedores creados.

# Utilización
Para usar el servicio, dirigase a `localhost:4180` e inicie sesión con Google. Automáticamente será redirigido a `localhost:8080`.

