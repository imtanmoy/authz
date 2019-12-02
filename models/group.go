package models

import (
	"context"
	"github.com/go-pg/pg/v9/orm"
	"time"
)

// Group represent groups table
type Group struct {
	ID             int32         `pg:"id,notnull"`
	Name           string        `pg:"name,notnull,unique:uk_groups_name_org"`
	OrganizationID int32         `pg:"organization_id,notnull,unique:uk_groups_name_org"`
	CreatedAt      time.Time     `pg:"created_at,notnull,default:now()"`
	UpdatedAt      time.Time     `pg:"updated_at,default:now()"`
	Users          []*User       `pg:"-"`
	Permissions    []*Permission `pg:"-"`
	Organization   *Organization
}

var _ orm.BeforeInsertHook = (*Group)(nil)
var _ orm.BeforeUpdateHook = (*Group)(nil)

//BeforeInsert group hooks
func (g *Group) BeforeInsert(ctx context.Context) (context.Context, error) {
	if g.CreatedAt.IsZero() {
		g.CreatedAt = time.Now()
	}
	if g.UpdatedAt.IsZero() {
		g.UpdatedAt = time.Now()
	}
	return ctx, nil
}

func (g *Group) BeforeUpdate(ctx context.Context) (context.Context, error) {
	g.UpdatedAt = time.Now()
	return ctx, nil
}
