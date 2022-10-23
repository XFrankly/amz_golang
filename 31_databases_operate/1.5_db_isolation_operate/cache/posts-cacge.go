package cahce

import (
	"github.com/pragmaticrveivews/golang-mux-api/entity"
)

/// 缓存接口
type PostCahce interface {
	Set(key string, value *entity.Post)
	Get(key string) *entity.Post
}
