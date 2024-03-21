package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

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
	log.Println("CreateTODO started")

	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	// dbへ保存する
	result, err := s.db.ExecContext(ctx, insert, subject, description)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()

	if err != nil {
		return nil, err
	}

	var todo = model.TODO{}
	todo.ID = id

	log.Println("CreateTodo started")

	var rows = s.db.QueryRowContext(ctx, confirm, id)
	var err2 = rows.Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err2 != nil {
		return nil, err2
	}

	log.Println("createTODO finished")
	return &todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	log.Println("ReadTODO started")

	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)
	// prev_idが指定されているかいないかは，0かそれ以外という定義だろうか

	var rows *sql.Rows = &sql.Rows{}
	if prevID == 0 {
		// クエリの中の?に束縛される値をquery以後の引数で指定してるように思われる
		inner_rows, err := s.db.QueryContext(ctx, read, size)
		if err != nil {
			return nil, err
		}
		rows = inner_rows
	} else {
		// クエリの中の?に束縛される値をquery以後の引数で指定してるように思われる
		inner_rows, err := s.db.QueryContext(ctx, readWithID, prevID, size)
		if err != nil {
			return nil, err
		}
		rows = inner_rows
	}

	// var todos []*model.TODO
	// 上記だとテストが通らない
	// 下記だとテストが通る
	// nilのままなのか？
	// 初期化されてないからの可能性が高い
	// 上記の書き方をよく見かけるから，勝手に初期化してるのかと思ったら
	// Goでも変数宣言時に必ず初期化するようにしよう

	var todos = []*model.TODO{}
	for rows.Next() {
		var todo = model.TODO{}
		var err = rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}

		todos = append(todos, &todo)
	}

	// // 条件式の前に単純な文を置けるようだ
	// // この変数はif文のスコープになるとのこと
	// if err := rows.Err(); err != nil {
	// 	return nil, err
	// }

	log.Println("ReadTODO finished")
	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	log.Println("UpdateTODO started")

	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	// dbを更新する
	// えぇ，varだとエラーになって，:=だとエラーにならない...
	// :=だとどちらか片方の値が違う変数名なら大丈夫っぽい
	// varだとどちらも違う変数名じゃないとダメっぽい
	result, err := s.db.ExecContext(ctx, update, subject, description, id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		// え，error側がポインタ渡すって情報はどこなの
		return nil, model.ErrNotFound{}
	}

	log.Println("Affected row count: ", rowsAffected)

	rows := s.db.QueryRowContext(ctx, confirm, id)
	todo := model.TODO{}
	todo.ID = id
	// err := rows.Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	// エラーになる．えぇ...
	err = rows.Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}

	log.Println("UpdateTODO finished")
	return &todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	log.Println("DeleteTODO started")

	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`
	// 空またはnilの場合は何もしない
	if len(ids) == 0 {
		return nil
	}

	// lenの数分だけ?を作る
	var deleteQuery = fmt.Sprintf(deleteFmt, strings.Repeat(", ?", len(ids)-1))

	// []int64を[]interface{}に変換
	var args = make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	// args... はスライスを展開しているようだ
	result, err := s.db.ExecContext(ctx, deleteQuery, args...)
	if err != nil {
		return err
	}

	// 影響を受けた行数を取得
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	log.Println("Affected row count: ", rowsAffected)

	// 影響を受けた行数が0の場合はNotFoundエラーを返す
	if rowsAffected == 0 {
		return &model.ErrNotFound{}
	}

	log.Println("DeleteTODO started")

	return nil
}
