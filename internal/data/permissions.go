package data

import (
	"context"
	"github.com/lib/pq"
	"time"
)

type Permissions []string

func (p Permissions) Include(code string) bool {
	for i := range p {
		if code == p[i] {
			return true
		}
	}
	return false
}

//GetAllForUser retrieves the permissions as a string array for the given user e.g {movies:write} will yeild an array
// of ["movie:write"]
func (m PermissionModel) GetAllForUser(userID int64) (Permissions, error) {

	query := `
	SELECT permissions.code
	FROM permissions
	INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
	WHERE users_permissions.user_id = $1`

	// I don't think I need to actually Join the Users Table as I can get the USER ID from user_permissions table only
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var permissions Permissions

	for rows.Next() {
		var permission string
		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil

}

//AddForUser adds ...
func (m PermissionModel) AddForUser(userId int64, codes ...string) error {
	query := `
	INSERT INTO users_permissions
	SELECT $1, permissions.id FROM permissions WHERE permissions.code = ANY($2)
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, userId, pq.Array(codes))
	return err
}
