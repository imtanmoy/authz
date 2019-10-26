package models

type Organization struct {
	ID    int32   `pg:"id"`
	Name  string  `pg:"name"`
	Users []*User `pg:"fk:organization_id"`
}

type User struct {
	ID             int32  `pg:"id"`
	Email          string `pg:"email"`
	OrganizationID int32  `pg:"organization_id"`
	Organization   *Organization
}
