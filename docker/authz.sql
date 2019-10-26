CREATE TABLE organizations
(
    id   INTEGER PRIMARY KEY NOT NULL,
    name VARCHAR(255)        NOT NULL
);

CREATE TABLE users
(
    id              INTEGER PRIMARY KEY NOT NULL,
    email           VARCHAR(255)        NOT NULL,
    organization_id INTEGER
);

ALTER TABLE users
    ADD CONSTRAINT fk_users_organization
        FOREIGN KEY (organization_id)
            REFERENCES organizations (id) ON DELETE CASCADE;


INSERT INTO organizations (id, name)
VALUES (1, 'Cramstack Ltd');

INSERT INTO users(id, email, organization_id)
VALUES (1, 'admin@cramstack.com', 1);


INSERT INTO organizations (id, name)
VALUES (2, 'Cramstack2 Ltd');