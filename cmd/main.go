/*
=====================================================================================================

  - main.go : Desenvolver um rate limiter em Go que possa ser configurado para limitar o número

  - máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.

  - Requisitos:

  - O rate limiter deve poder trabalhar como um middleware que é injetado ao servidor web

  - O rate limiter deve permitir a configuração do número máximo de requisições permitidas por segundo.

  - O rate limiter deve ter ter a opção de escolher o tempo de bloqueio do IP ou do Token caso a

  - quantidade de requisições tenha sido excedida.

  - As configurações de limite devem ser realizadas via variáveis de ambiente ou em um arquivo “.env”

  - na pasta raiz.

  - Deve ser possível configurar o rate limiter tanto para limitação por IP quanto por token de acesso.

  - O sistema deve responder adequadamente quando o limite é excedido:

  - Código HTTP: 429

  - Mensagem: you have reached the maximum number of requests or actions allowed within a certain time frame

  - Todas as informações de "limiter” devem ser armazenadas e consultadas de um banco de dados Redis.

  - Você pode utilizar docker-compose para subir o Redis.

  - Crie uma “strategy” que permita trocar facilmente o Redis por outro mecanismo de persistência.

  - A lógica do limiter deve estar separada do middleware.

=====================================================================================================
*/
package main

import (
	"Desafios/3/goexperts-desafio-rate-limiter/internal/limiter"
	"Desafios/3/goexperts-desafio-rate-limiter/internal/middleware"
	"fmt"
	"net/http"

	"log"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

/*
==========================================================
  - Função: configInit
  - Descrição : Função que carrega os parametros do
  - arquivo .ENV
  - Parametros :
  - Retorno:

==========================================================
*/
func configInit() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Erro ao ler o arquivo de configuração: %s\n", err)
		return
	}
	viper.AutomaticEnv()
}

func main() {
	configInit() //Carrega parâmetros

	//Guarda parâmetros carregados
	rateLimitIP := viper.GetInt("RATE_LIMIT_IP")
	rateLimitToken := viper.GetInt("RATE_LIMIT_TOKEN")
	blockTimeIP := viper.GetInt("BLOCK_TIME_IP")
	blockTimeToken := viper.GetInt("BLOCK_TIME_TOKEN")
	redisHost := viper.GetString("REDIS_HOST")
	redisPort := viper.GetString("REDIS_PORT")

	//fmt.Printf("Limite IP: %d, Limite Token: %d, Tempo de BloqueioIP: %d, Tempo de BloqueioToken: %d, Host: %s, Porta: %s\n", rateLimitIP, rateLimitToken, blockTimeIP, blockTimeToken, redisHost, redisPort)

	//Cria client do Redis
	redisDB := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
	})

	rediStore := limiter.NewRedisStore(redisDB)

	rateLim := limiter.NewRateLimiter(rediStore, rateLimitIP, rateLimitToken, blockTimeIP, blockTimeToken)

	//fmt.Println(rateLim)

	//Injeta middleware no servidor
	http.Handle("/", middleware.RateLimiter(rateLim, http.HandlerFunc(handler)))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %s\n", err)
	}

	fmt.Printf("Servidor rodando na porta 8080 ...\n")
}

/*
==========================================================

  - Função: handler
  - Descrição : Mostra mensagem "Hello World!"
  - Parametros :
  - w - http.ResponseWriter
  - r - http.Request
  - Retorno:

==========================================================
*/
func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello, World!</h1>\n"))
}
