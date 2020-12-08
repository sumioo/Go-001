package service

import (
	"context"

	"github.com/sumioo/week02/dao"
	"github.com/sumioo/week02/model"
)

func GetArticle(ctx context.Context, id int) (*model.Article, error) {
	return dao.GetArticle(ctx, id)
}
