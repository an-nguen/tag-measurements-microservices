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


ALTER TABLE public.users OWNER TO an;
ALTER TABLE public.roles OWNER TO an;
ALTER TABLE public.privileges OWNER TO an;

CREATE TABLE IF NOT EXISTS public.users_roles (
                                                  user_id SERIAL,
                                                  role_id SERIAL,
                                                  constraint users_roles_pkey primary key (user_id, role_id),
                                                  foreign key (user_id) references users(id),
                                                  foreign key (role_id) references roles(id)
);

create table if not exists public.roles_privileges (
                                                        role_id serial,
                                                        privilege_id serial,
                                                        constraint roles_privileges_pkey primary key (role_id, privilege_id),
                                                        foreign key (role_id) references roles(id) on delete set default on update cascade,
                                                        foreign key (privilege_id) references privileges(id) on delete set default on update cascade
)
