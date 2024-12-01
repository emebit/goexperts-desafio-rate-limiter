package middleware

import (
	"context"
	"log"
	"net"
	"net/http"

	"Desafios/3/goexperts-desafio-rate-limiter/internal/limiter"
)

/*
==========================================================
  - Função: RateLimiter
  - Descrição : Middlware que orquestra a verficação do limite
  - e o bloqueio das requisições.
  - Parametros :
  - rl - Estrutura que armazena os dados parâmetros de limite
  - e bloqueio - tipo: ponteiro para RateLimiter
  - next - um handle do http - tipo: http.Handler
  - Retorno:

==========================================================
*/
func RateLimiter(rl *limiter.RateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}
		//Pega o token, se houver
		token := r.Header.Get("API_KEY")
		ctx := context.Background()
		//Se token foi informado, passa o token na key
		if token != "" {
			if !rl.Allow(ctx, "token:"+token, rl.LimitToken) {
				log.Println("Token limit exceeded:", token)
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}
		} else { //Senão, passa o IP na key
			if !rl.Allow(ctx, "ip:"+ip, rl.LimitIP) {
				log.Println("IP limit exceeded:", ip)
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
