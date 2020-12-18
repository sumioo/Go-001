package biz

import (
	"context"
	"fmt"
)

type Article struct {
	ID      int
	Title   string
	Content string
}

type ArticleRepo interface {
	Create(context.Context, Article) (int, error)
}

func (article *Article) Save(ctx context.Context, repo ArticleRepo) error {
	id, err := repo.Create(ctx, *article)
	if err != nil {
		fmt.Println(err)
		return err
	}
	article.ID = id
	return nil
}
