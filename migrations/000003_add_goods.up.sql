CREATE TABLE IF NOT EXISTS GOODS (
    id serial NOT NULL,
    project_id int NOT NULL,
    name VARCHAR(256) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    priority INT NOT NULL DEFAULT 1,
    removed bool NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY(id, project_id)
    );

CREATE INDEX ON GOODS(name);