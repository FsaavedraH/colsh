package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsuarioRepository struct {
	Pool *pgxpool.Pool
}

type UsuarioResumen struct {
	IDUsuario uuid.UUID `json:"id_usuario"`
	Nombre    string    `json:"nombre"`
	Email     string    `json:"email"`
	Rol       string    `json:"rol"`
}

// RF-27: lista todos los usuarios del sistema
func (r *UsuarioRepository) ListarUsuarios(ctx context.Context) ([]UsuarioResumen, error) {
	rows, err := r.Pool.Query(ctx,
		`SELECT id_usuario, nombre, email, rol FROM usuario ORDER BY rol, nombre`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultado []UsuarioResumen
	for rows.Next() {
		var u UsuarioResumen
		if err := rows.Scan(&u.IDUsuario, &u.Nombre, &u.Email, &u.Rol); err != nil {
			return nil, err
		}
		resultado = append(resultado, u)
	}
	return resultado, nil
}