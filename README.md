# Desafio: Clean Architecture - Listagem de Pedidos
Este projeto é a resolução do desafio proposto no módulo de Clean Architecture do [Curso GoExpert](https://goexpert.fullcycle.com.br/pos-goexpert/), onde foi solicitado a criação de um Use Case para a listagem de pedidos. O Use Case desenvolvido pode ser acionado de três formas distintas: via REST, gRPC e GraphQL.

## Funcionalidades
- REST Endpoint: Realize uma chamada GET em /order para listar os pedidos.
- gRPC: Use o serviço ListOrders para listar os pedidos.
- GraphQL: Utilize a query ListOrders para realizar a consulta de pedidos.

- ## Execução do Projeto
Para executar o projeto, basta utilizar o comando abaixo:

```bash
docker-compose up --detach
```

Este comando irá subir todas as dependências necessárias, incluindo:

- A aplicação em si, criada com Go;
- Banco de dados MySQL;
- RabbitMQ para o gerenciamento de filas.

## Acesso aos Endpoints
1. REST
   Para listar os pedidos via REST, faça uma chamada GET para:
```bash
http://localhost:8000/order
```
2. GraphQL
   Para consultar pedidos via GraphQL, acesse o Playground:
```bash
http://localhost:8080/playground
```
3. gRPC
   Para realizar chamadas via gRPC, utilize o comando docker genérico abaixo:
```bash
docker run -it --network host --rm -v "$(pwd):/mount:ro" \
    ghcr.io/ktr0731/evans:latest \
      --path /mount/internal/infra/grpc/protofiles \
      --proto order.proto \
      --host localhost \
      --port 50051 \
      repl
```
```bash
package pb
service OrderService
call ListOrders
```

## Testes
A aplicação já estará disponível e pronta para uso via REST, GraphQL, e gRPC após a execução do comando docker-compose up. <br>
Siga os passos descritos na seção [Acesso aos Endpoints](#acesso-aos-endpoints) para verificar o funcionamento de cada interface.