package limiter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// Estrutura que contem o ponteiro para o client do Redis
type RedisStore struct {
	client *redis.Client
}

/*
==========================================================
  - Função: NewRedisSotre
  - Descrição : Função que cria/popula a estrutura RedisStore
  - Parametros :
  - client - cliente do redis - tipo redis.Client
  - limitIP - Limite de requisiçoes para IP - tipo: int
  - Retorno: Ponteiro de RedisStore

==========================================================
*/
func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client: client}
}

/*
==========================================================
  - Função: Incr
  - Descrição : Função que chama o metodo Incr do Redis
  - Parametros :
  - ctx - contexto onde está rodando - tipo context
  - key - chave da requisição no Redis - tipo: string
  - Retorno: valor atualizado e erro

==========================================================
*/
func (r *RedisStore) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

/*
==========================================================
  - Função: Expire
  - Descrição : Função que chama o metodo Expire do Redis
  - Parametros :
  - ctx - contexto onde está rodando - tipo context
  - key - chave da requisição no Redis - tipo: string
  - duration - tempo de bloqueio - tipo: time.Duration
  - Retorno: erro

==========================================================
*/
func (r *RedisStore) Expire(ctx context.Context, key string, duration time.Duration) error {
	return r.client.Expire(ctx, key, duration).Err()
}
