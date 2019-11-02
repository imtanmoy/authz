package models

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
