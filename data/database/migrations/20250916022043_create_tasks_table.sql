-- +goose Up
-- +goose StatementBegin
create table if not exists tasks (
	id serial not null primary key,
	title varchar(100) not null,
	description varchar,
	status varchar,
	created_at timestamp not null default current_timestamp,
	updated_at timestamp not null default current_timestamp,
	deleted_at timestamp
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists tasks;
-- +goose StatementEnd
