CREATE TABLE logs (
   id BIGSERIAL primary key,
   avg_value float8 not null,
   symbol varchar not null,
   created_at TIMESTAMP default now()
);