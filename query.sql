create table if not exists users (
	id UUID primary key default gen_random_uuid(), 
	name VARCHAR(255) not null, 
	email VARCHAR(255) unique not null, 
	password_hash text not null, 
	created_at timestamp default current_timestamp
);

create table if not exists rooms (
	id UUID primary key default gen_random_uuid(), 
	name VARCHAR(255) not null, 
	type VARCHAR(50) not null, 
	description text, 
	created_at timestamp default current_timestamp, 
	updated_at timestamp default current_timestamp
);

create extension if not exists btree_gist;

create table if not exists bookings (
	id UUID primary key default gen_random_uuid(),
	room_id UUID references rooms(id),
	user_id UUID references users(id),
	start_date date not null,
	end_date date not null,
	created_at timestamp default current_timestamp,
	updated_at timestamp default current_timestamp,
	deleted_at timestamp,
	check (end_date > start_date), 
	constraint no_overlapping_bookings 
		exclude using gist (
			room_id with =,
			daterange(start_date, end_date, '[)') with &&
		) where (deleted_at is null)
);

create index if not exists idx_rooms_type on rooms(type);

create index if not exists idx_bookings_room_id on bookings(room_id) where deleted_at is null;

insert into rooms (name, type, description) values 
('Pajajaran Suites', 'suite', 'AC 2 PK'),
('Pajajaran Single', 'single', 'AC 1 PK'),
('Pajajaran Double', 'double', 'AC 1.5 PK');

-- same date bookings
insert into bookings (room_id, user_id, start_date, end_date) values ('2903405c-aec8-4a7b-a216-e0fab564b1a4', '64fed4bf-bab8-406f-b48e-3d21ecf72b7f', '2026-04-06', '2026-04-07');insert into bookings (room_id, user_id, start_date, end_date) values ('2903405c-aec8-4a7b-a216-e0fab564b1a4', '64fed4bf-bab8-406f-b48e-3d21ecf72b7f', '2026-04-06', '2026-04-07');

-- adjacent bookings
insert into bookings (room_id, user_id, start_date, end_date) values ('2903405c-aec8-4a7b-a216-e0fab564b1a4', '64fed4bf-bab8-406f-b48e-3d21ecf72b7f', '2026-04-06', '2026-04-07');insert into bookings (room_id, user_id, start_date, end_date) values ('2903405c-aec8-4a7b-a216-e0fab564b1a4', '64fed4bf-bab8-406f-b48e-3d21ecf72b7f', '2026-04-07', '2026-04-08');

-- partial overlap bookings
insert into bookings (room_id, user_id, start_date, end_date) values ('2903405c-aec8-4a7b-a216-e0fab564b1a4', '64fed4bf-bab8-406f-b48e-3d21ecf72b7f', '2026-04-06', '2026-04-10');insert into bookings (room_id, user_id, start_date, end_date) values ('2903405c-aec8-4a7b-a216-e0fab564b1a4', '64fed4bf-bab8-406f-b48e-3d21ecf72b7f', '2026-04-09', '2026-04-11');

insert into bookings (room_id, user_id, start_date, end_date) values ('2903405c-aec8-4a7b-a216-e0fab564b1a4', '64fed4bf-bab8-406f-b48e-3d21ecf72b7f', '2026-04-09', '2026-04-11');insert into bookings (room_id, user_id, start_date, end_date) values ('2903405c-aec8-4a7b-a216-e0fab564b1a4', '64fed4bf-bab8-406f-b48e-3d21ecf72b7f', '2026-04-06', '2026-04-10');

-- same date bookings different rooms
insert into bookings (room_id, user_id, start_date, end_date) values ('2903405c-aec8-4a7b-a216-e0fab564b1a4', '64fed4bf-bab8-406f-b48e-3d21ecf72b7f', '2026-04-06', '2026-04-07');insert into bookings (room_id, user_id, start_date, end_date) values ('046425be-d45e-434e-9cde-7a5d1150f7e0', '64fed4bf-bab8-406f-b48e-3d21ecf72b7f', '2026-04-06', '2026-04-07');

