package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type Bot struct {
	Token  string `json:"token"`
	Handle string `json:"handle"`
}

type Editable struct {
	Id     int `json:"id"`
	ChatId int `json:"chat_id"`
}

type User struct {
	ent.Schema
}

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{entsql.Annotation{Schema: "manager"}}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.
			Int("id").
			Unique().
			Immutable().
			Positive(),
		field.
			JSON("bot", Bot{}).
			Default(Bot{}),
		field.
			JSON("editable", Editable{}).
			Default(Editable{}),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
