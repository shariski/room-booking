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

