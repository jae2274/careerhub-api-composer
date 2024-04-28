package httputils

import (
	"context"
	"net/http"

	"github.com/jae2274/goutils/llog"
	"github.com/jae2274/goutils/terr"
)

func IsNotLoggedIn(ctx context.Context, w http.ResponseWriter, ok bool) bool {
	if !ok {
		llog.LogErr(ctx, terr.New("No claims in context")) //이미 middleware에서 처리되기 때문에 혹여 이곳에 도달한다면 예외 상황이다.
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return true
	}

	return false
}

func IsError(ctx context.Context, w http.ResponseWriter, err error) bool {
	if err != nil {
		llog.LogErr(ctx, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return true
	}

	return false
}
