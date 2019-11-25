CREATE TABLE organizations
(
    id   BIGINT PRIMARY KEY NOT NULL,
    name VARCHAR(255)       NOT NULL
);

CREATE TABLE users
(
    id              BIGINT PRIMARY KEY NOT NULL,
    email           VARCHAR(128)       NOT NULL,
    organization_id BIGINT             NOT NULL
);

create type permission_type as enum('feature', 'resource');

CREATE TABLE permissions
(
    id              BIGSERIAL PRIMARY KEY NOT NULL,
    name            VARCHAR(128)          NOT NULL,
    action          VARCHAR(32)           NOT NULL,
    type            permission_type       NOT NULL DEFAULT 'feature',
    organization_id INTEGER               NOT NULL,
    created_at      TIMESTAMP             NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP             NULL
);

CREATE TABLE groups
(
    id              BIGSERIAL PRIMARY KEY NOT NULL,
    name            VARCHAR(128)          NOT NULL,
    organization_id BIGINT                NOT NULL,
    created_at      TIMESTAMP             NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP
--     deleted_at      TIMESTAMP             NULL     DEFAULT NOW()
);

CREATE TABLE casbin_rules
(
    p_type           VARCHAR(10) , 
    v0              VARCHAR(256), 
    v1              VARCHAR(256), 
    v2              VARCHAR(256), 
    v3              VARCHAR(256), 
    v4              VARCHAR(256), 
    v5              VARCHAR(256)
);

ALTER TABLE users
    ADD CONSTRAINT fk_users_organization
        FOREIGN KEY (organization_id)
            REFERENCES organizations (id) ON DELETE CASCADE;

ALTER TABLE groups
    ADD CONSTRAINT fk_groups_organization
        FOREIGN KEY (organization_id)
            REFERENCES organizations (id) ON DELETE CASCADE;

ALTER TABLE groups
    ADD CONSTRAINT uk_groups_name_org UNIQUE (name, organization_id);


ALTER TABLE permissions
    ADD CONSTRAINT fk_permissions_organization
        FOREIGN KEY (organization_id)
            REFERENCES organizations (id) ON DELETE CASCADE;

ALTER TABLE permissions
    ADD CONSTRAINT uk_permissions_name_org UNIQUE (name, organization_id);

-- ALTER TABLE groups
--     ADD CONSTRAINT uk_groups_name_org_del UNIQUE (name, organization_id, deleted_at);


INSERT INTO organizations (id, name)
VALUES (1, 'Cramstack Ltd');

INSERT INTO users(id, email, organization_id)
VALUES (1, 'admin@cramstack.com', 1);


INSERT INTO organizations (id, name)
VALUES (2, 'Cramstack2 Ltd');

INSERT INTO groups (id, name, organization_id)
VALUES (1, 'ADMIN', 1);

INSERT INTO permissions (id, name, action, organization_id)
VALUES (1, 'PERMISSION_1', 'ALL', 1);

INSERT INTO permissions (id, name, action, organization_id)
VALUES (2, 'PERMISSION_2', 'ALL', 1);

