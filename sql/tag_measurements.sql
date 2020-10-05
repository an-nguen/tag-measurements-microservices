-- Main tables --
CREATE TABLE public.temperature_zone (
                                        id SERIAL PRIMARY KEY,
                                        name VARCHAR(32) UNIQUE NOT NULL,
                                        description VARCHAR(128) NULL,
                                        lower_temp_limit DOUBLE PRECISION,
                                        higher_temp_limit DOUBLE PRECISION,
                                        notify_emails varchar(512)
);

CREATE TABLE public.wireless_tag_account (
                                             email varchar(128) primary key default null,
                                             password varchar(128) default null
);

CREATE TABLE public.tag_manager (
                                    mac MACADDR primary key,
                                    name VARCHAR(32) NOT NULL,
                                    email varchar(128)
);


CREATE TABLE public.tag (
                            uuid UUID PRIMARY KEY,
                            name VARCHAR(32) NOT NULL,
                            mac MACADDR NOT NULL,
                            verification_date DATE DEFAULT current_date,
                            higher_temperature_limit DOUBLE PRECISION NULL,
                            lower_temperature_limit DOUBLE PRECISION NULL
);

CREATE TABLE public.measurement (
                                 id serial primary key,
                                 date timestamp default current_timestamp,
                                 temperature DOUBLE PRECISION,
                                 humidity DOUBLE PRECISION,
                                 voltage DOUBLE PRECISION,
                                 signal DOUBLE PRECISION,
                                 tag_uuid UUID NOT NULL
);


CREATE INDEX measurement_timestamp_tag_uuid_index ON public.measurement (id, date, tag_uuid);

CREATE TABLE public.temperature_zone_tag (
                                            temperature_zone_id SERIAL not null,
                                            tag_uuid uuid,
                                            constraint temperature_zone_tag_pkey primary key (temperature_zone_id, tag_uuid),
                                            foreign key (temperature_zone_id) references temperature_zone(id) on delete set default on update set default,
                                            foreign key (tag_uuid) references tag(uuid) on delete set default on update set default
);

CREATE INDEX temperature_zone_tag_index ON public.temperature_zone_tag (temperature_zone_id, tag_uuid);

INSERT INTO public.temperature_zone (name, description) VALUES ('NONE', 'TEST');

CREATE OR REPLACE FUNCTION remove_duplicate_data ()
    RETURNS void AS $$
DECLARE
    r INTEGER;
    row RECORD;
    found INT[];
    cnt INT;
BEGIN
    FOR row IN (SELECT d.date AS date,
                       d.temperature AS temperature,
                       d.humidity AS humidity
                FROM measurement d)
        LOOP
            cnt = 0;
            SELECT ARRAY (SELECT t.date FROM measurement t
                          WHERE t.humidity = row.humidity AND
                                  t.temperature = row.temperature AND
                                  t.date = row.date) as found;
            cnt = array_length(found, 1);
            IF (cnt > 0) THEN
                FOREACH r IN ARRAY found LOOP
                        DELETE FROM measurement WHERE date = r;
                    end loop;
            end if;
        end loop;
END
$$ LANGUAGE plpgsql;
