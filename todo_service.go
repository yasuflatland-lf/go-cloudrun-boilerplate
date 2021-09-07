package main

import (
	"context"
	"golang.org/x/xerrors"
	"time"
)

type (
	TodoService interface {
		List(status bool, page, pagesize int, order string) (todos []*Todo, totalRows int, err error)
		Create(todo *Todo) (*Todo, error)
		CreateInBatches(todos []Todo) ([]Todo, error)
		Delete(ID int64) (rowsAffected int64, err error)
		Get(id int64) (*Todo, error)
		Update(todo *Todo) (*Todo, error)
	}

	todoService struct {
		Repository CloudSQL
	}

	Todo struct {
		// field named `ID` will be used as a primary field by default
		// https://gorm.io/docs/conventions.html#ID-as-Primary-Key
		ID        int64     `gorm:"primary_key;AUTO_INCREMENT;column:id;type:bigint;" faker:"-"`
		Slug      string    `gorm:"column:slug;type:varchar;" faker:"uuid_hyphenated"`
		Task      string    `gorm:"column:task;type:text;"`
		Status    bool      `gorm:"column:status;type:tinyint;default:0;"`
		CreatedAt time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;" faker:"-"`
		UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP;" faker:"-"`
	}
)

func NewTodoService(ctx context.Context) TodoService {
	t := &todoService{}
	t.Repository = NewCloudSQL(ctx)
	return t
}

func (t *todoService) List(status bool, page, pagesize int, order string) (todos []*Todo, totalRows int, err error) {
	resultOrm := t.Repository.DB().Model(&Todo{})

	if page > 0 {
		offset := (page - 1) * pagesize
		resultOrm = resultOrm.Offset(offset).Limit(pagesize)
	} else {
		resultOrm = resultOrm.Limit(pagesize)
	}

	if order != "" {
		resultOrm = resultOrm.Order(order)
	}

	resultOrm = resultOrm.Where("status = ?", status)

	if err = resultOrm.Find(&todos).Error; err != nil {
		return nil, -1, xerrors.Errorf("List : can not find the record : %w", err)
	}

	return todos, len(todos), nil
}

// Delete
// https://gorm.io/docs/delete.html
func (t *todoService) Delete(ID int64) (rowsAffected int64, err error) {
	todo := &Todo{}
	tx := t.Repository.DB().First(todo, ID)
	if tx.Error != nil {
		return -1, xerrors.Errorf("Delete : can not find the record : %w", tx.Error)
	}

	tx = t.Repository.DB().Delete(todo)
	if tx.Error != nil {
		return -1, xerrors.Errorf("Can not Delete : %w", tx.Error)
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

// Create in Batches
// https://gorm.io/docs/create.html
func (t *todoService) CreateInBatches(todos []Todo) ([]Todo, error) {
	tx := t.Repository.DB().CreateInBatches(todos,len(todos))

	if tx.Error != nil {
		return nil, xerrors.Errorf("Create : %w", tx.Error)
	}

	return todos, nil
}


// Query
// https://gorm.io/docs/query.html
// https://gorm.io/docs/advanced_query.html
func (t *todoService) Get(id int64) (*Todo, error) {
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
