package middle

import (
	"net/http"
	"torchi/internal/api/handler"
	"torchi/internal/pkg/token"
)

func AuthMiddleware(tp *token.TokenProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. 쿠키에서 토큰 추출
			cookie, err := r.Cookie(handler.AccessCookieKey)
			if err != nil {
				// 토큰이 없으면 401 에러 혹은 로그인 페이지로 리다이렉트
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// 2. 토큰 검증
			claims, err := tp.Validate(cookie.Value)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// 3. Context에 유저 정보 저장 (컨트롤러에서 꺼내 쓸 수 있도록)
			ctx := token.ContextWith(r.Context(), claims)

			// 4. 다음 핸들러 실행
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
