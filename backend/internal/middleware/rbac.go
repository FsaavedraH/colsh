package middleware

import (
	"encoding/json"
	"net/http"
)

// RequireRole crea un middleware que solo permite continuar si el rol
// que llega en el header X-User-Role esta entre los roles permitidos.
// RF-27, RNF-04, Escenario Arquitectural 5.
func RequireRole(rolesPermitidos ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rol := r.Header.Get("X-User-Role")

			if rol == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Falta el header X-User-Role. Debe identificar el rol del usuario.",
				})
				return
			}

			permitido := false
			for _, p := range rolesPermitidos {
				if rol == p {
					permitido = true
					break
				}
			}

			if !permitido {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Rol no autorizado para esta operacion.",
					"rol_recibido": rol,
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}