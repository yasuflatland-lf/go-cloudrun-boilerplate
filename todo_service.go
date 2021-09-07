package main

import (
	"context"
	"golang.org/x/xerrors"
	"time"
)

type (
	TodoService interface {
		Create(todo *Todo) (*Todo, error)
		Delete(todo *Todo) (rowsAffected int64, err error)
		Get(id int) (*Todo, error)
	}

	todoService struct {
		Repository CloudSQL
	}

	//https://qiita.com/sky0621/items/90a8b6e7dd097cd671cd
	// https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
	Todo struct {
		//[ 0] id                                             bigint               null: false  primary: true   isArray: false  auto: true   col: bigint          len: -1      default: []
		ID int64 `gorm:"primary_key;AUTO_INCREMENT;column:id;type:bigint;"`
		//[ 1] slug                                           varchar              null: false  primary: false  isArray: false  auto: false  col: varchar         len: 0       default: []
		Slug string `gorm:"column:slug;type:varchar;"`
		//[ 2] task                                           text                 null: false  primary: false  isArray: false  auto: false  col: text            len: 0       default: []
		Task string `gorm:"column:task;type:text;"`
		//[ 3] status                                         tinyint              null: true   primary: false  isArray: false  auto: false  col: tinyint         len: -1      default: [0]
		Status bool `gorm:"column:status;type:tinyint;default:0;"`
		//[ 4] created_at                                     timestamp            null: false  primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: [CURRENT_TIMESTAMP]
		CreatedAt time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;"`
		//[ 5] updated_at                                     timestamp            null: false  primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: [CURRENT_TIMESTAMP]
		UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP;"`
	}
)

func NewTodoService(ctx context.Context) TodoService {
	t := &todoService{}
	t.Repository = NewCloudSQL(ctx)
	return t
}

func (t *todoService) List(status bool, order string, offset int64, limit int64) ([]Todo, error) {
	todos := &[]Todo{}
	tx := t.Repository.DB().
		Where(map[string]interface{}{"status": status}).
		Order("updated_at " + order).
		Offset(int(offset)).
		Limit(int(limit)).
		Find(&todos)

	if tx.Error != nil {
		return []Todo{}, xerrors.Errorf("Create : %w", tx.Error)
	}

	return *todos, nil
}

// Delete
// https://gorm.io/docs/delete.html
func (t *todoService) Delete(todo *Todo) (rowsAffected int64, err error) {
	tx := t.Repository.DB().Delete(todo)

	if tx.Error != nil {
		return -1, xerrors.Errorf("Create : %w", tx.Error)
	}

	return tx.RowsAffected, nil
}

// Create
// https://gorm.io/docs/create.html
func (t *todoService) Create(todo *Todo) (*Todo, error) {
	tx := t.Repository.DB().Save(todo)

	if tx.Error != nil {
		return nil, xerrors.Errorf("Create : %w", tx.Error)
	}

	return todo, nil
}

// Query
// https://gorm.io/docs/query.html
// https://gorm.io/docs/advanced_query.html
func (t *todoService) Get(id int) (*Todo, error) {
	todo := &Todo{}
	if err := t.Repository.DB().First(todo, id).Error; err != nil {
		return nil, xerrors.Errorf("Get : %w", err)
	}

	return todo, nil
}

// Update
// https://gorm.io/docs/update.html
func (t *todoService) Update(todo *Todo) (*Todo, error) {
	tx := t.Repository.DB().Save(todo)

	if tx.Error != nil {
		return nil, xerrors.Errorf("Update : %w", tx.Error)
	}

	return todo, nil
}
