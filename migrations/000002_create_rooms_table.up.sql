create table if not exists rooms (
	id UUID primary key default gen_random_uuid(), 
	name VARCHAR(255) not null, 
	type VARCHAR(50) not null, 
	description text, 
	created_at timestamp default current_timestamp, 
	updated_at timestamp default current_timestamp
);

