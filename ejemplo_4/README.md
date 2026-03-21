Ejemplo 4
===

En este ejemplo integramos AWS Lambda con API Gateway, el flujo es el siguiente:
- Request llega a API gateway
- El request hace el trigger de un evento que invoca a AWS Lambda
- La función lambda lee el contenido del event y lo procesa como si fuera un request HTTP, claramente son objetos propios de AWS y no de `net/http`
