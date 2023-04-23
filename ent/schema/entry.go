package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Entry holds the schema definition for the Entry entity.
type Entry struct {
	ent.Schema
}

// Fields of the Entry.
func (Entry) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").MinLen(1).Unique(),
		field.String("body").MinLen(1),
		field.String("tag"),
	}
}

// Edges of the Entry.
func (Entry) Edges() []ent.Edge {
	return nil
}
