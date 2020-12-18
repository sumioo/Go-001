package schema

import "github.com/facebook/ent"
import "github.com/facebook/ent/schema/field"

// Article holds the schema definition for the Article entity.
type Article struct {
	ent.Schema
}

// Fields of the Article.
func (Article) Fields() []ent.Field {
	return []ent.Field{
		field.String("title"),
		field.String("content"),
	}
}

// Edges of the Article.
func (Article) Edges() []ent.Edge {
	return nil
}
