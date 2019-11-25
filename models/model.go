package models

import (
	"context"
	"github.com/go-pg/pg/v9/orm"
	"time"
)

func init() {
	// Register many to many model so ORM can better recognize m2m relation.
	// This should be done before dependant models are used.
	//orm.RegisterTable((*UserGroup)(nil))
}

// Organization represent organizations table
type Organization struct {
	ID    int32   `pg:"id,notnull,unique"`
	Name  string  `pg:"name,notnull"`
	Users []*User `pg:"fk:organization_id"`
}

// User represent users table
type User struct {
	ID             int32         `pg:"id,notnull,unique"`
	Email          string        `pg:"email,notnull,unique"`
	OrganizationID int32         `pg:"organization_id,notnull"`
	Groups         []*Group      `pg:"many2many:users_groups,fk:user_id,joinFK:group_id"`
	Organization   *Organization `pg:"-"`
}

// Group represent groups table
type Group struct {
	ID             int32         `pg:"id,notnull"`
	Name           string        `pg:"name,notnull,unique:uk_groups_name_org"`
	OrganizationID int32         `pg:"organization_id,notnull,unique:uk_groups_name_org"`
	CreatedAt      time.Time     `pg:"created_at,notnull,default:now()"`
	UpdatedAt      time.Time     `pg:"updated_at"`
	Users          []*User       `pg:"-"`
	Organization   *Organization `pg:"-"`
}

var _ orm.BeforeInsertHook = (*Group)(nil)
var _ orm.AfterInsertHook = (*Group)(nil)

//BeforeInsert group hooks
func (g *Group) BeforeInsert(ctx context.Context) (context.Context, error) {
	if g.CreatedAt.IsZero() {
		g.CreatedAt = time.Now()
	}
	return ctx, nil
}

//AfterInsert group hooks
func (g *Group) AfterInsert(ctx context.Context) error {
	return nil // here we can update the cache
}

// Permission represent permissions table
type Permission struct {
	ID             int32         `pg:"id,notnull,unique"`
	Name           string        `pg:"name,notnull"`
	OrganizationID int32         `pg:"organization_id,notnull"`
	Action         string        `pg:"action,notnull"`
	Type           string        `pg:"type,type:permission_type,default:feature"`
	CreatedAt      time.Time     `pg:"created_at,notnull,default:now()"`
	UpdatedAt      time.Time     `pg:"updated_at"`
	Users          []*User       `pg:"-"`
	Organization   *Organization `pg:"-"`
}

var _ orm.BeforeInsertHook = (*Permission)(nil)
var _ orm.AfterInsertHook = (*Permission)(nil)

//BeforeInsert hooks
func (g *Permission) BeforeInsert(ctx context.Context) (context.Context, error) {
	if g.CreatedAt.IsZero() {
		g.CreatedAt = time.Now()
	}
	return ctx, nil
}

//AfterInsert hooks
func (g *Permission) AfterInsert(ctx context.Context) error {
	return nil // here we can update the cache
}
