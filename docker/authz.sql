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

-- CREATE TABLE permissions
-- (
--     id              BIGSERIAL PRIMARY KEY NOT NULL,
--     name            VARCHAR(128)          NOT NULL,
--     action          VARCHAR(32)           NOT NULL,
--     type            VARCHAR(32)           NOT NULL,
--     organization_id INTEGER               NOT NULL,
--     created_at      TIMESTAMP             NOT NULL DEFAULT NOW(),
--     updated_at      TIMESTAMP             NULL     DEFAULT NOW()
-- );

CREATE TABLE groups
(
    id              BIGSERIAL PRIMARY KEY NOT NULL,
    name            VARCHAR(128)          NOT NULL,
    organization_id BIGINT                NOT NULL,
    created_at      TIMESTAMP             NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP
--     deleted_at      TIMESTAMP             NULL     DEFAULT NOW()
);

CREATE TABLE users_groups
(
    user_id  BIGINT NOT NULL,
    group_id BIGINT NOT NULL
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

ALTER TABLE users_groups
    ADD CONSTRAINT pk_users_groups
        PRIMARY KEY (user_id, group_id);

ALTER TABLE users_groups
    ADD CONSTRAINT fk_users_groups_user
        FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE users_groups
    ADD CONSTRAINT fk_users_groups_group
        FOREIGN KEY (group_id) REFERENCES groups (id);

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

INSERT INTO users_groups (user_id, group_id)
values (1, 1)