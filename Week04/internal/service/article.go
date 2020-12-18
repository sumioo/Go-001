package service

import (
	"context"

	pb "hello.com/api"
	"hello.com/internal/biz"
)

type articleServer struct {
	pb.UnimplementedArticleServer
	repo biz.ArticleRepo
}

func NewArticleServer(repo biz.ArticleRepo) *articleServer {
	return &articleServer{repo: repo}
}

func (s *articleServer) CreateArticle(ctx context.Context, in *pb.CreateArticleRequest) (*pb.ArticleReply, error) {
	article := biz.Article{Title: in.Title, Content: in.Content}
	article.Save(ctx, s.repo)
	return &pb.ArticleReply{Id: int32(article.ID), Title: article.Title, Content: article.Content}, nil
}
