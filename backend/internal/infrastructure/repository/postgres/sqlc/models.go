package sqlc

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Queries struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Queries {
	return &Queries{db: db}
}

type User struct {
	ID             uuid.UUID
	Email          string
	HashedPassword string
	Name           string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

type CreateUserParams struct {
	ID             uuid.UUID
	Email          string
	HashedPassword string
	Name           string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, `
		INSERT INTO users (id, email, hashed_password, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, email, hashed_password, name, created_at, updated_at, deleted_at
	`, arg.ID, arg.Email, arg.HashedPassword, arg.Name, arg.CreatedAt, arg.UpdatedAt)
	var u User
	err := row.Scan(&u.ID, &u.Email, &u.HashedPassword, &u.Name, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
	return u, err
}

func (q *Queries) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRow(ctx, `
		SELECT id, email, hashed_password, name, created_at, updated_at, deleted_at
		FROM users WHERE id = $1 AND deleted_at IS NULL
	`, id)
	var u User
	err := row.Scan(&u.ID, &u.Email, &u.HashedPassword, &u.Name, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
	return u, err
}

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, `
		SELECT id, email, hashed_password, name, created_at, updated_at, deleted_at
		FROM users WHERE email = $1 AND deleted_at IS NULL
	`, email)
	var u User
	err := row.Scan(&u.ID, &u.Email, &u.HashedPassword, &u.Name, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
	return u, err
}

func (q *Queries) SoftDeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, `UPDATE users SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`, id)
	return err
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.Exec(ctx, `
		UPDATE users SET email = $2, hashed_password = $3, name = $4, updated_at = $5
		WHERE id = $1 AND deleted_at IS NULL
	`, arg.ID, arg.Email, arg.HashedPassword, arg.Name, arg.UpdatedAt)
	return err
}

type UpdateUserParams struct {
	ID             uuid.UUID
	Email          string
	HashedPassword string
	Name           string
	UpdatedAt      time.Time
}

func (q *Queries) ListUserRoles(ctx context.Context, userID uuid.UUID) ([]string, error) {
	rows, err := q.db.Query(ctx, `
		SELECT r.name FROM user_roles ur
		JOIN roles r ON r.id = ur.role_id
		WHERE ur.user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var roles []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		roles = append(roles, name)
	}
	return roles, rows.Err()
}

func (q *Queries) ListUserPermissions(ctx context.Context, userID uuid.UUID) ([]string, error) {
	rows, err := q.db.Query(ctx, `
		SELECT DISTINCT p.code
		FROM user_roles ur
		JOIN roles r ON r.id = ur.role_id
		JOIN role_permissions rp ON rp.role_id = r.id
		JOIN permissions p ON p.id = rp.permission_id
		WHERE ur.user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var perms []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		perms = append(perms, code)
	}
	return perms, rows.Err()
}

type Role struct {
	ID          uuid.UUID
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type CreateRoleParams struct {
	ID          uuid.UUID
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (q *Queries) CreateRole(ctx context.Context, arg CreateRoleParams) (Role, error) {
	row := q.db.QueryRow(ctx, `
		INSERT INTO roles (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, description, created_at, updated_at, deleted_at
	`, arg.ID, arg.Name, arg.Description, arg.CreatedAt, arg.UpdatedAt)
	var r Role
	err := row.Scan(&r.ID, &r.Name, &r.Description, &r.CreatedAt, &r.UpdatedAt, &r.DeletedAt)
	return r, err
}

func (q *Queries) GetRoleByID(ctx context.Context, id uuid.UUID) (Role, error) {
	row := q.db.QueryRow(ctx, `
		SELECT id, name, description, created_at, updated_at, deleted_at
		FROM roles WHERE id = $1 AND deleted_at IS NULL
	`, id)
	var r Role
	err := row.Scan(&r.ID, &r.Name, &r.Description, &r.CreatedAt, &r.UpdatedAt, &r.DeletedAt)
	return r, err
}

func (q *Queries) GetRoleByName(ctx context.Context, name string) (Role, error) {
	row := q.db.QueryRow(ctx, `
		SELECT id, name, description, created_at, updated_at, deleted_at
		FROM roles WHERE name = $1 AND deleted_at IS NULL
	`, name)
	var r Role
	err := row.Scan(&r.ID, &r.Name, &r.Description, &r.CreatedAt, &r.UpdatedAt, &r.DeletedAt)
	return r, err
}

func (q *Queries) ListRoles(ctx context.Context) ([]Role, error) {
	rows, err := q.db.Query(ctx, `
		SELECT id, name, description, created_at, updated_at, deleted_at
		FROM roles WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var roles []Role
	for rows.Next() {
		var r Role
		if err := rows.Scan(&r.ID, &r.Name, &r.Description, &r.CreatedAt, &r.UpdatedAt, &r.DeletedAt); err != nil {
			return nil, err
		}
		roles = append(roles, r)
	}
	return roles, rows.Err()
}

func (q *Queries) SoftDeleteRole(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, `UPDATE roles SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`, id)
	return err
}

type UpdateRoleParams struct {
	ID          uuid.UUID
	Name        string
	Description string
	UpdatedAt   time.Time
}

func (q *Queries) UpdateRole(ctx context.Context, arg UpdateRoleParams) error {
	_, err := q.db.Exec(ctx, `
		UPDATE roles SET name = $2, description = $3, updated_at = $4
		WHERE id = $1 AND deleted_at IS NULL
	`, arg.ID, arg.Name, arg.Description, arg.UpdatedAt)
	return err
}

type UpsertRolePermissionParams struct {
	RoleID       uuid.UUID
	PermissionID uuid.UUID
}

func (q *Queries) UpsertRolePermission(ctx context.Context, arg UpsertRolePermissionParams) error {
	_, err := q.db.Exec(ctx, `
		INSERT INTO role_permissions (role_id, permission_id)
		VALUES ($1, $2)
		ON CONFLICT (role_id, permission_id) DO NOTHING
	`, arg.RoleID, arg.PermissionID)
	return err
}

func (q *Queries) ListRolePermissions(ctx context.Context, roleID uuid.UUID) ([]string, error) {
	rows, err := q.db.Query(ctx, `
		SELECT p.code FROM role_permissions rp
		JOIN permissions p ON p.id = rp.permission_id
		WHERE rp.role_id = $1
	`, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var perms []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		perms = append(perms, code)
	}
	return perms, rows.Err()
}

type Permission struct {
	ID          uuid.UUID
	Code        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type CreatePermissionParams struct {
	ID          uuid.UUID
	Code        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (q *Queries) CreatePermission(ctx context.Context, arg CreatePermissionParams) (Permission, error) {
	row := q.db.QueryRow(ctx, `
		INSERT INTO permissions (id, code, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, code, description, created_at, updated_at, deleted_at
	`, arg.ID, arg.Code, arg.Description, arg.CreatedAt, arg.UpdatedAt)
	var p Permission
	err := row.Scan(&p.ID, &p.Code, &p.Description, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt)
	return p, err
}

func (q *Queries) GetPermissionByCode(ctx context.Context, code string) (Permission, error) {
	row := q.db.QueryRow(ctx, `
		SELECT id, code, description, created_at, updated_at, deleted_at
		FROM permissions WHERE code = $1 AND deleted_at IS NULL
	`, code)
	var p Permission
	err := row.Scan(&p.ID, &p.Code, &p.Description, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt)
	return p, err
}

func (q *Queries) ListPermissions(ctx context.Context) ([]Permission, error) {
	rows, err := q.db.Query(ctx, `
		SELECT id, code, description, created_at, updated_at, deleted_at
		FROM permissions WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var perms []Permission
	for rows.Next() {
		var p Permission
		if err := rows.Scan(&p.ID, &p.Code, &p.Description, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt); err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}
	return perms, rows.Err()
}

type Product struct {
	ID          uuid.UUID
	SKU         string
	Name        string
	Description string
	Price       float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type CreateProductParams struct {
	ID          uuid.UUID
	SKU         string
	Name        string
	Description string
	Price       float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error) {
	row := q.db.QueryRow(ctx, `
		INSERT INTO products (id, sku, name, description, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, sku, name, description, price, created_at, updated_at, deleted_at
	`, arg.ID, arg.SKU, arg.Name, arg.Description, arg.Price, arg.CreatedAt, arg.UpdatedAt)
	var p Product
	err := row.Scan(&p.ID, &p.SKU, &p.Name, &p.Description, &p.Price, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt)
	return p, err
}

func (q *Queries) GetProductByID(ctx context.Context, id uuid.UUID) (Product, error) {
	row := q.db.QueryRow(ctx, `
		SELECT id, sku, name, description, price, created_at, updated_at, deleted_at
		FROM products WHERE id = $1 AND deleted_at IS NULL
	`, id)
	var p Product
	err := row.Scan(&p.ID, &p.SKU, &p.Name, &p.Description, &p.Price, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt)
	return p, err
}

func (q *Queries) ListProducts(ctx context.Context) ([]Product, error) {
	rows, err := q.db.Query(ctx, `
		SELECT id, sku, name, description, price, created_at, updated_at, deleted_at
		FROM products WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.SKU, &p.Name, &p.Description, &p.Price, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

type UpdateProductParams struct {
	ID          uuid.UUID
	SKU         string
	Name        string
	Description string
	Price       float64
	UpdatedAt   time.Time
}

func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) error {
	_, err := q.db.Exec(ctx, `
		UPDATE products SET sku = $2, name = $3, description = $4, price = $5, updated_at = $6
		WHERE id = $1 AND deleted_at IS NULL
	`, arg.ID, arg.SKU, arg.Name, arg.Description, arg.Price, arg.UpdatedAt)
	return err
}

func (q *Queries) SoftDeleteProduct(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, `UPDATE products SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`, id)
	return err
}

type Inventory struct {
	ID        uuid.UUID
	ProductID uuid.UUID
	Quantity  int64
	Location  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type CreateInventoryParams struct {
	ID        uuid.UUID
	ProductID uuid.UUID
	Quantity  int64
	Location  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) CreateInventory(ctx context.Context, arg CreateInventoryParams) (Inventory, error) {
	row := q.db.QueryRow(ctx, `
		INSERT INTO inventories (id, product_id, quantity, location, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, product_id, quantity, location, created_at, updated_at, deleted_at
	`, arg.ID, arg.ProductID, arg.Quantity, arg.Location, arg.CreatedAt, arg.UpdatedAt)
	var inv Inventory
	err := row.Scan(&inv.ID, &inv.ProductID, &inv.Quantity, &inv.Location, &inv.CreatedAt, &inv.UpdatedAt, &inv.DeletedAt)
	return inv, err
}

func (q *Queries) GetInventoryByID(ctx context.Context, id uuid.UUID) (Inventory, error) {
	row := q.db.QueryRow(ctx, `
		SELECT id, product_id, quantity, location, created_at, updated_at, deleted_at
		FROM inventories WHERE id = $1 AND deleted_at IS NULL
	`, id)
	var inv Inventory
	err := row.Scan(&inv.ID, &inv.ProductID, &inv.Quantity, &inv.Location, &inv.CreatedAt, &inv.UpdatedAt, &inv.DeletedAt)
	return inv, err
}

func (q *Queries) ListInventoryByProduct(ctx context.Context, productID uuid.UUID) ([]Inventory, error) {
	rows, err := q.db.Query(ctx, `
		SELECT id, product_id, quantity, location, created_at, updated_at, deleted_at
		FROM inventories WHERE product_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []Inventory
	for rows.Next() {
		var inv Inventory
		if err := rows.Scan(&inv.ID, &inv.ProductID, &inv.Quantity, &inv.Location, &inv.CreatedAt, &inv.UpdatedAt, &inv.DeletedAt); err != nil {
			return nil, err
		}
		res = append(res, inv)
	}
	return res, rows.Err()
}

type AdjustInventoryParams struct {
	ID        uuid.UUID
	Delta     int64
	UpdatedAt time.Time
}

func (q *Queries) AdjustInventory(ctx context.Context, arg AdjustInventoryParams) error {
	_, err := q.db.Exec(ctx, `
		UPDATE inventories SET quantity = quantity + $2, updated_at = $3
		WHERE id = $1 AND deleted_at IS NULL
	`, arg.ID, arg.Delta, arg.UpdatedAt)
	return err
}

func (q *Queries) SoftDeleteInventory(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, `UPDATE inventories SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`, id)
	return err
}

type AuditLog struct {
	ID        uuid.UUID
	ActorID   uuid.NullUUID
	Action    string
	Resource  string
	Metadata  map[string]any
	CreatedAt time.Time
}

type CreateAuditLogParams struct {
	ID        uuid.UUID
	ActorID   *uuid.UUID
	Action    string
	Resource  string
	Metadata  map[string]any
	CreatedAt time.Time
}

func (q *Queries) CreateAuditLog(ctx context.Context, arg CreateAuditLogParams) error {
	_, err := q.db.Exec(ctx, `
		INSERT INTO audit_logs (id, actor_id, action, resource, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, arg.ID, arg.ActorID, arg.Action, arg.Resource, arg.Metadata, arg.CreatedAt)
	return err
}

func (q *Queries) ListAuditLogs(ctx context.Context, actorID *uuid.UUID) ([]AuditLog, error) {
	var rows pgx.Rows
	var err error
	if actorID == nil {
		rows, err = q.db.Query(ctx, `
			SELECT id, actor_id, action, resource, metadata, created_at
			FROM audit_logs ORDER BY created_at DESC
		`)
	} else {
		rows, err = q.db.Query(ctx, `
			SELECT id, actor_id, action, resource, metadata, created_at
			FROM audit_logs WHERE actor_id = $1 ORDER BY created_at DESC
		`, *actorID)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []AuditLog
	for rows.Next() {
		var l AuditLog
		if err := rows.Scan(&l.ID, &l.ActorID, &l.Action, &l.Resource, &l.Metadata, &l.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, l)
	}
	return res, rows.Err()
}
