package limiter

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// Estrutura que armazena os dados parâmetros de limite e bloqueio
type RateLimiter struct {
	store          Store
	LimitIP        int
	LimitToken     int
	BlockTimeIP    time.Duration
	BlockTimeToken time.Duration
}

/*
==========================================================
  - Função: NewRateLimiter
  - Descrição : Função que cria/popula a estrutura RateLimiter
  - Parametros :
  - store - interface de base de dados - tipo Store
  - limitIP - Limite de requisiçoes para IP - tipo: int
  - limitToken - Limite de requisiçoes para Token - tipo: int
  - blockIP - Tempo de bloqueio de requisiçoes para IP - tipo: int
  - blockToken - Tempo de bloqueio de requisiçoes para Token - tipo: int
  - Retorno: Ponteiro de RateLimiter

==========================================================
*/
func NewRateLimiter(store Store, limitIP, limitToken, blockTimeIP int, blockTimeToken int) *RateLimiter {
	return &RateLimiter{
		store:          store,
		LimitIP:        limitIP,
		LimitToken:     limitToken,
		BlockTimeIP:    time.Duration(blockTimeIP) * time.Second,
		BlockTimeToken: time.Duration(blockTimeToken) * time.Second,
	}
}

/*
==========================================================
  - Função: Allow
  - Descrição : Função que incrementa o nr de requisições
  - e o tempo de expriação para a chave no Redis.
  - Parametros :
  - ctx - contexto onde etsa sendo executado - tipo context
  - key - Chave do Redis - tipo: String
  - limit - limite de requisições - tipo: int
  - Retorno: Se o limite de requisições está dentro do limite
  - ou false se der erro - tipo: bool

==========================================================
*/
func (rl *RateLimiter) Allow(ctx context.Context, key string, limit int) bool {
	count, err := rl.store.Incr(ctx, key)
	if err != nil {
		fmt.Println("Erro ao incrementar:", err)
		return false
	}
	//fmt.Printf("Contagem atual para %s: %d\n", key, count)
	//Se for a primeira vez, define valor de expiração
	if count == 1 {
		var blockTime time.Duration

		//Se verifica se a chave é de TOKEN ou IP e usa o valor correspondente de bloqueio
		if strings.Contains(key, "token:") {
			blockTime = rl.BlockTimeToken
		} else {
			blockTime = rl.BlockTimeIP
		}

		rl.store.Expire(ctx, key, blockTime)
	}
	return count <= int64(limit)
}
