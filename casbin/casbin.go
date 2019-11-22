package casbin

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/go-pg/pg/v9"

	"github.com/imtanmoy/authz/casbin/adapter"
)

//Enforcer casbin Enforcer
var Enforcer *casbin.Enforcer

// InitCasbin initialze the Conf
func InitCasbin(db *pg.DB) {
	text :=
		`
		[request_definition]
		r = sub, obj, act
		
		[policy_definition]
		p = sub, obj, act
		
		[role_definition]
		g = _, _
		
		[policy_effect]
		e = some(where (p.eft == allow))
		
		[matchers]
		m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
		`

	m, err := model.NewModelFromString(text)
	if err != nil {
		panic(err)
	}
	// Load the policy rules from the .CSV file adapter.
	// Replace it with your adapter to avoid files.
	a := adapter.NewAdapter(db)

	// Create the enforcer.
	Enforcer, err = casbin.NewEnforcer(m, a)
	if err != nil {
		panic(err)
	}
}
