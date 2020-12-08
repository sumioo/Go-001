package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/sumioo/week02/dao"
	"github.com/sumioo/week02/service"
)

func main() {
	articleID := 1
	article, err := service.GetArticle(context.Background(), articleID)
	if errors.Is(err, dao.ObjectNotFound) {
		log.Printf("404: %v", err)
		return
	}
	if err != nil {
		log.Printf("500: %+v", err)
		return
	}
	resp, err := json.Marshal(article)
	if err != nil {
		log.Printf("500: %v", err)
		return
	}
	fmt.Println(string(resp))
}
