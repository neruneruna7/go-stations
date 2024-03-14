package service

import (
	"context"
	"database/sql"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	// dbへ保存する
	var result1, e1 = s.db.ExecContext(ctx, insert, subject, description)
	if e1 != nil {
		return nil, e1
	}
	var id, e2 = result1.LastInsertId()

	if e2 != nil {
		return nil, e2
	}

	var todo model.TODO
	todo.ID = int(id)

	var rows = s.db.QueryRowContext(ctx, confirm, id)
	var e3 = rows.Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)

	return &todo, e3
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	// // クエリの中の?に束縛される値をquery以後の引数で指定してるように思われる
	// s.db.QueryRowContext(ctx, read, size)

	return nil, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	// dbを更新する
	var result1, e1 = s.db.ExecContext(ctx, update, subject, description)
	if e1 != nil {
		return nil, e1
	}

	var affected_row_count, e2 = result1.RowsAffected()
	if e2 != nil {
		return nil, e2
	}

	if affected_row_count == 0 {
		// え，error側がポインタ渡すって情報はどこなの
		return nil, &model.ErrNotFound{}
	}

	// var rows = s.db.QueryRowContext(ctx, confirm, id)
	// var todo model.TODO
	// todo.ID = int(id)
	// var e3 = rows.Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	// if e3 != nil {
	// 	return nil, e3
	// }

	return nil, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	return nil
}
