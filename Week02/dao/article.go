package dao

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/sumioo/week02/model"
)

var ObjectNotFound = errors.New("object not found")

func GetArticle(ctx context.Context, id int) (*model.Article, error) {
	article := &model.Article{}
	row := DB.QueryRowContext(ctx, "select id, title, content from article where id = $1", id)
	err := row.Scan(&article.ID, &article.Title, &article.Content)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ObjectNotFound
		}
		return nil, errors.Wrapf(err, "fail to query article")
	}
	return article, nil

}
