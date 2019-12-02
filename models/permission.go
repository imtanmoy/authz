package models

import (
	"context"
	"github.com/go-pg/pg/v9/orm"
	"time"
)

// Permission represent permissions table
type Permission struct {
	ID             int32     `pg:"id,notnull,unique"`
	Name           string    `pg:"name,notnull"`
	OrganizationID int32     `pg:"organization_id,notnull"`
	Action         string    `pg:"action,notnull"`
	Type           string    `pg:"type,type:permission_type,default:feature"`
	CreatedAt      time.Time `pg:"created_at,notnull,default:now()"`
	UpdatedAt      time.Time `pg:"updated_at"`
	//Users          []*User   `pg:"-"`
	Organization *Organization
}

var _ orm.BeforeInsertHook = (*Permission)(nil)
var _ orm.BeforeUpdateHook = (*Permission)(nil)

//BeforeInsert hooks
func (p *Permission) BeforeInsert(ctx context.Context) (context.Context, error) {
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}
	return ctx, nil
}

func (p *Permission) BeforeUpdate(ctx context.Context) (context.Context, error) {
	p.UpdatedAt = time.Now()
	return ctx, nil
}
