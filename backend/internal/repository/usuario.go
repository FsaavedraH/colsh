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

type UsuarioConCredenciales struct {
	IDUsuario    uuid.UUID
	Nombre       string
	Email        string
	Rol          string
	PasswordHash string
}

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

// BuscarPorEmail obtiene el usuario junto con su hash de contraseña, para login.
func (r *UsuarioRepository) BuscarPorEmail(ctx context.Context, email string) (*UsuarioConCredenciales, error) {
	var u UsuarioConCredenciales
	err := r.Pool.QueryRow(ctx,
		`SELECT id_usuario, nombre, email, rol, password_hash FROM usuario WHERE email = $1`, email,
	).Scan(&u.IDUsuario, &u.Nombre, &u.Email, &u.Rol, &u.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &u, nil
}


// CrearUsuario inserta un usuario nuevo con el rol indicado.
func (r *UsuarioRepository) CrearUsuario(ctx context.Context, idUsuario uuid.UUID, nombre, email, rol, passwordHash string) error {
	_, err := r.Pool.Exec(ctx,
		`INSERT INTO usuario (id_usuario, nombre, email, rol, password_hash) VALUES ($1, $2, $3, $4, $5)`,
		idUsuario, nombre, email, rol, passwordHash,
	)
	return err
}