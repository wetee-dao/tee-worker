package graph

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vedhavyas/go-subkey/v2"
	"github.com/vedhavyas/go-subkey/v2/sr25519"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"wetee.app/worker/graph/model"
)

var (
	loginStatCtxKey = &contextKey{"loginStat"}
)

type contextKey struct {
	name string
}

func AuthCheck(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	user := ctx.Value(loginStatCtxKey).(*model.User)
	if user == nil {
		err = gqlerror.Errorf("Please log in first.")
		return
	}
	if user.Timestamp+360000 < time.Now().Unix() {
		err = gqlerror.Errorf("Login expired, please log in again.")
		return
	}
	return next(ctx)
}

// Middleware decodes the share session cookie and packs the session into context
func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := InitUser(w, r)

			// put it in context
			ctx := context.WithValue(r.Context(), loginStatCtxKey, user)

			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func InitUser(w http.ResponseWriter, r *http.Request) *model.User {
	header := r.Header
	token, ok := header["Authorization"]
	if !ok || len(token) == 0 {
		return nil
	}

	return decodeToken(token[0])
}

func decodeToken(tokenStr string) *model.User {
	token := strings.Split(tokenStr, "||")

	if len(token) != 2 {
		return nil
	} else {
		bt, terr := subkey.DecodeHex(token[0])
		if !terr {
			return nil
		}

		user := &model.User{}
		err := json.Unmarshal(bt, user)
		if err != nil {
			return nil
		}

		// 解析地址
		_, pubkeyBytes, err := subkey.SS58Decode(user.Address)
		if err != nil {
			return nil
		}

		// 解析公钥
		pubkey, err := sr25519.Scheme{}.FromPublicKey(pubkeyBytes)
		if err != nil {
			return nil
		}

		// 解析签名
		sig, chainerr := subkey.DecodeHex(token[1])
		if !chainerr {
			return nil
		}

		// 验证签名
		ok := pubkey.Verify(bt, sig)
		if !ok {
			// return nil
		}

		return user
	}
}
