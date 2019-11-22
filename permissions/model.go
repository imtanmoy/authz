package permissions

import (
	"context"
	"time"

	"github.com/go-pg/pg/v9/orm"
)

// Permission represent permissions table
type Permission struct {
	ID             int32     `pg:"id,notnull"`
	Name           string    `json:"name"`
	OrganizationID int32     `json:"org_id" gorm:"not null"`
	Action         string    `json:"action"`
	Type           string    `sql:"type:ENUM('feature', 'resource')" gorm:"default:'feature'"`
	CreatedAt      time.Time `pg:"created_at,notnull,default:now()"`
	UpdatedAt      time.Time `pg:"updated_at"`
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
