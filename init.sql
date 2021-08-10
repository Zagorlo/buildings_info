drop table if exists buildings;

create table buildings
(
    id bigserial not null,
    name varchar(50) not null,
    floors_count int2 not null,
    parking_count int2 not null,
    parking_available bool not null,
    constraint buildings_pk primary key (id),
    constraint check_building_name check (length("name") > 2),
    constraint check_building_floors_count check (floors_count > 0),
    constraint check_building_parkings check (parking_count >= 0 and (parking_count = 0 or parking_available = true))
);

drop index if exists buildings_idx;

create index buildings_idx on buildings(id);

drop table if exists buildings_updates;

create table buildings_updates
(
    id bigserial not null,
    name varchar(50) null,
    floors_count int2 null,
    parking_count int2 null,
    parking_available bool null,
    removed_at timestamp not null
);

drop index if exists buildings_updates_idx;

create index buildings_updates_idx on buildings_updates(id);

drop function if exists stash_building_info;

create function stash_building_info()
 returns trigger
 language plpgsql
as $function$
 begin
  if (old.* = new.*) then
    return new;
  end if;
  insert into buildings_updates values(
    old.id,
    case when old.name <> new.name then old.name else null end,
    case when old.floors_count <> new.floors_count then old.floors_count else null end,
    case when old.parking_count <> new.parking_count then old.parking_count else null end,
    case when old.parking_available <> new.parking_available then old.parking_available else null end,
    now()
  );
  return new;
 end;
$function$
;

drop trigger if exists stash_building_before_update on buildings;

create trigger stash_building_before_update before
update
    on
    buildings for each row execute function stash_building_info();

insert into buildings(name, floors_count, parking_count, parking_available)
values
('aaa', 10, 20, true),
('bbb', 2, 0, false),
('ccc', 15, 450, true);

insert into buildings_updates(name, floors_count, parking_count, parking_available, removed_at)
values
('ddd', 5, 35, true, '2012-12-20'::timestamp - interval '2 day'),
('eee', 1, 0, false, '2012-12-20'::timestamp);
