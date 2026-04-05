create table if not exists users (
	id UUID primary key default gen_random_uuid(), 
	name VARCHAR(255) not null, 
	email VARCHAR(255) unique not null, 
	password_hash text not null, 
	created_at timestamp default current_timestamp
);
