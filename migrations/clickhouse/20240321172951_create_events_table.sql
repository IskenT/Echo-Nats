
The error message you're encountering indicates that the way you've defined the ENGINE for the events table is outdated. Let's fix the code to address the deprecation warning and create the table successfully.

Here's the modified code:

SQL
-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events (
  id        Int,
  project_id Int,
  name      String,
  description String,
  priority  Int,
  removed    Bool,
  EventTime TIMESTAMP
) ENGINE = MergeTree ORDER BY (EventTime);

CREATE INDEX index_id_events ON events(id) TYPE minmax GRANULARITY 3;
CREATE INDEX index_project_id_events ON events(project_id) TYPE minmax GRANULARITY 3;
CREATE INDEX index_name_events ON events(name) TYPE ngrambf_v1(4, 1024, 1, 0) GRANULARITY 1;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd