package data

import (
	"context"
	"fmt"
	"log"

	"database/sql"

	"hello.com/ent"
	"hello.com/internal/biz"

	"github.com/facebook/ent/dialect"
	entsql "github.com/facebook/ent/dialect/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "zclng"
	password = ""
	dbname   = "zclng"
)

type articleRepo struct {
	client *ent.Client
}

func Open(databaseURL string) *ent.Client {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	// Create an ent.Driver from `db`.
	drv := entsql.OpenDB(dialect.Postgres, db)
	return ent.NewClient(ent.Driver(drv))
}

func NewArticleRepo() articleRepo {
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s?sslmode=disable", user, password, host)
	client := Open(databaseURL)
	return articleRepo{client: client}
}

func (repo articleRepo) Create(ctx context.Context, article biz.Article) (int, error) {
	articleEntity, err := repo.client.Article.Create().SetTitle(article.Title).SetContent(article.Content).Save(ctx)
	if err != nil {
		return 0, fmt.Errorf("fail creating article:%v", err)
	}
	log.Println("article was create: ", articleEntity)
	return articleEntity.ID, nil
}
