CREATE TABLE public.users (
                              id  SERIAL PRIMARY KEY,
                              username  VARCHAR(256) UNIQUE NOT NULL,
                              password  VARCHAR(256) NOT NULL
);


CREATE TABLE public.roles (
                              id  SERIAL PRIMARY KEY,
                              name  VARCHAR(256) UNIQUE NOT NULL,
                              description  VARCHAR(1024)
);

create table public.privileges (
                                   id serial primary key ,
                                   name varchar(256) unique not null
);

insert into public.roles (name, description) values ('ADMIN', '');

ALTER TABLE public.users OWNER TO an;
ALTER TABLE public.roles OWNER TO an;
ALTER TABLE public.privileges OWNER TO an;

CREATE TABLE IF NOT EXISTS public.users_roles (
                                                  user_id SERIAL,
                                                  role_id SERIAL,
                                                  constraint users_roles_pkey primary key (user_id, role_id),
                                                  foreign key (user_id) references users(id) on delete set default on update cascade,
                                                  foreign key (role_id) references roles(id) on delete set default on update cascade
);

create table if not exists public.roles_privileges (
                                                       role_id serial,
                                                       privilege_id serial,
                                                       constraint roles_privileges_pkey primary key (role_id, privilege_id),
                                                       foreign key (role_id) references roles(id) on delete set default on update cascade,
                                                       foreign key (privilege_id) references privileges(id) on delete set default on update cascade
);

INSERT INTO privileges (name) VALUES ('ALLOW_SHOW_TAG_NUMBER');
INSERT INTO privileges (name) VALUES ('ALLOW_SHOW_UUID');
INSERT INTO privileges (name) VALUES ('ALLOW_VERIFICATION_DATE');
INSERT INTO privileges (name) VALUES ('ALLOW_SHOW_TEMPERATURE');
INSERT INTO privileges (name) VALUES ('ALLOW_SHOW_HUMIDITY');
INSERT INTO privileges (name) VALUES ('ALLOW_SHOW_VOLTAGE');
INSERT INTO privileges (name) VALUES ('ALLOW_SHOW_BATTERY_REMAINING');
INSERT INTO privileges (name) VALUES ('ALLOW_SHOW_SIGNAL');
INSERT INTO privileges (name) VALUES ('ALLOW_TAG_EDIT');
INSERT INTO privileges (name) VALUES ('CRUD_WST_ACCOUNTS');
INSERT INTO privileges (name) VALUES ('CRUD_USER');
INSERT INTO privileges (name) VALUES ('CRUD_ROLE');
INSERT INTO privileges (name) VALUES ('CRUD_PRIVILEGE');
INSERT INTO privileges (name) VALUES ('CRUD_TEMPERATURE_ZONE');

do $$
    <<first_block>>
        declare role_admin_id int := -1;
        declare pr_id int := 0;
        declare crud_privileges varchar[] := array['CRUD_WST_ACCOUNTS', 'CRUD_USER', 'CRUD_ROLE', 'CRUD_PRIVILEGE', 'CRUD_TEMPERATURE_ZONE'];
        declare privilege_name varchar;
    begin
        select id into role_admin_id from roles where name = 'ADMIN';
        if role_admin_id <> -1 then
            foreach privilege_name in array crud_privileges loop
                    select id into pr_id from privileges where name = privilege_name;
                    insert into roles_privileges (role_id, privilege_id) values (role_admin_id, pr_id);
                end loop;
        end if;
    end first_block $$;
