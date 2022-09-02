package db

import (
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
)

type Retrier struct {
}

func NewRetrier() *Retrier {

	r := Retrier{}

	return &r
}

func (r *Retrier) Insert(query squirrel.InsertBuilder) error {

	_, err := query.Exec()
	for err != nil && strings.Contains(err.Error(), "pq: deadlock detected") {
		fmt.Println("retrying insert")
		_, err = query.Exec()
	}

	return err
}

func (r *Retrier) Update(query squirrel.UpdateBuilder) error {

	_, err := query.Exec()
	for err != nil && strings.Contains(err.Error(), "pq: deadlock detected") {
		fmt.Println("retrying update")
		_, err = query.Exec()
	}

	return err
}

func (r *Retrier) Delete(query squirrel.DeleteBuilder) error {

	_, err := query.Exec()
	for err != nil && strings.Contains(err.Error(), "pq: deadlock detected") {
		fmt.Println("retrying delete")
		_, err = query.Exec()
	}

	return err
}
