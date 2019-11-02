package models

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v9/orm"
	"time"
)

type Organization struct {
	ID    int32   `pg:"id,notnull,unique"`
	Name  string  `pg:"name,notnull"`
	Users []*User `pg:"fk:organization_id"`
}

type User struct {
	ID             int32  `pg:"id,notnull,unique"`
	Email          string `pg:"email,notnull,unique"`
	OrganizationID int32  `pg:"organization_id,notnull"`
	Organization   *Organization
}

type Group struct {
	ID             int32     `pg:"id,notnull"`
	Name           string    `pg:"name,notnull,unique:uk_groups_name_org"`
	OrganizationID int32     `pg:"organization_id,notnull,unique:uk_groups_name_org"`
	CreatedAt      time.Time `pg:"created_at,notnull,default:now()"`
	UpdatedAt      time.Time `pg:"updated_at"`
	Organization   *Organization
}

var _ orm.BeforeInsertHook = (*Group)(nil)
var _ orm.AfterInsertHook = (*Group)(nil)

func (g *Group) BeforeInsert(ctx context.Context) (context.Context, error) {
	if g.CreatedAt.IsZero() {
		g.CreatedAt = time.Now()
	}
	return ctx, nil
}

func (g *Group) AfterInsert(ctx context.Context) error {
	fmt.Println(g.ID)
	return nil // here we can update the cache
}
