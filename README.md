# Goexperts Rate Limiter

Este projeto é uma implementação de um limitador de taxa em Go, utilizando Redis como armazenamento. 
O limitador permite limitar requisições por IP e por token.


## Instalação

1. Clone o repositório:
```bash
   git clone https://github.com/emebit/goexperts-rate-limiter.git
   cd goexperts-rate-limiter/ 
```

2. Build o projeto:
```bash
      docker compose up --build
```

## Executando

1. ### Fazendo chamadas HTTP:

Dentro da pasta http há um arquivo test_limit.http que pode ser executado para fazer as chamadas por ip ou com o
token no header.

2. ### Testes com CURL:

   Para testar o limite de IP:
```bash
   for i in {1..10}; do curl -i http://localhost:8080/; done
```

   Para testar com um token:
```bash
   for i in {1..10}; do curl -i -H "API_KEY: teste_123" http://localhost:8080/; done  
```
   Para testar com o Apache Bench:
   Instalar o Apache Bench se necessário, e executar o codigo:
```bash
   ab -n 10 -c 1 http://localhost:8080/
```
