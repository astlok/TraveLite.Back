ALTER
    USER postgres WITH ENCRYPTED PASSWORD 'admin';

DROP SCHEMA IF EXISTS travelite CASCADE;

CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE SCHEMA travelite;

CREATE TABLE travelite.users
(
    id       SERIAL PRIMARY KEY NOT NULL,
    email    CITEXT UNIQUE      NOT NULL,
    nickname CITEXT UNIQUE      NOT NULL,
    password TEXT               NOT NULL,
    img      TEXT DEFAULT ''
);

INSERT INTO travelite.users(email, nickname, password)
VALUES ('kek@mem.ru',
        'mem',
        '228');

CREATE TABLE travelite.sessions
(
    user_id INT  NOT NULL,
    session TEXT NOT NULL,
    FOREIGN KEY (user_id)
        REFERENCES travelite.users (id)
);

CREATE TABLE travelite.region
(
    id   SERIAL PRIMARY KEY NOT NULL,
    name CITEXT UNIQUE      NOT NULL
);


CREATE TYPE travelite.seasons AS ENUM ('Зима', 'Весна', 'Лето', 'Осень', '');
CREATE TYPE travelite.hike_type AS ENUM ('Пеший', 'Горный', 'Водный', 'Альпинизм', 'Велотуризм', 'Бег', 'Мото', 'Авто', 'Скитур', 'Лыжный', 'Горный велотуризм', 'Бездорожье', 'Ски-альпинизм', 'Снегоступы');
CREATE TABLE travelite.trek
(
    id              SERIAL PRIMARY KEY          NOT NULL,
    name            CITEXT                      NOT NULL,
    difficult       INT                         NOT NULL,
    days            INT                         NOT NULL,
    description     CITEXT            DEFAULT '',
    best_time_to_go TRAVELITE.SEASONS DEFAULT '',
    type            TRAVELITE.HIKE_TYPE         NOT NULL,
    climb           INT                         NOT NULL,
    region          TEXT                        NOT NULL,
    creator_id      INT                         NOT NULL,
    is_moderate     BOOL              DEFAULT false,
    route           GEOMETRY(LINESTRINGZ, 4326) NOT NULL,
    start           GEOMETRY(POINTZ, 4326)      NOT NULL,
    FOREIGN KEY (creator_id)
        REFERENCES travelite.users (id)
);

-- SELECT enum_range(NULL::travelite.hike_type);

-- INSERT INTO travelite.trek (name,
--                             difficult,
--                             days,
--                             description,
--                             best_time_to_go,
--                             type,
--                             climb,
--                             region,
--                             creator_id,
--                             is_moderate,
--                             route,
--                             start)
-- VALUES (':name',
--         '8',
--         '2',
--         'asdkf',
--         'Зима',
--         'Пеший',
--         '288',
--         'Москва',
--         '1',
--         'true',
--         'LINESTRING Z(0 0 2, 1 1 2, 2 1 2, 2 2 3)',
--         'POINT Z (0 0 2)')
-- RETURNING id;

-- INSERT INTO travelite.marks(trek_id, point, title, description, image)
-- VALUES (':trek_id',
--         ':point',
--         ':title',
--         ':description',
--         ':image');

select id, name, difficult, days, description, best_time_to_go, type, climb, region, creator_id, is_moderate, ST_AsText(route) as route, ST_AsText(start) as start from travelite.trek;

CREATE TABLE travelite.marks
(
    trek_id     INT                    NOT NULL,
    point       GEOMETRY(POINTZ, 4326) NOT NULL,
    title       TEXT                   NOT NULL,
    description TEXT DEFAULT '',
    image       TEXT DEFAULT '',
    FOREIGN KEY (trek_id)
        REFERENCES travelite.trek (id) ON DELETE CASCADE
);

-- CREATE TABLE travelite.trek_rating
-- (
--     user_id INT    NOT NULL,
--     trek_id INT    NOT NULL,
--     rating  FLOAT8 NOT NULL,
--     FOREIGN KEY (user_id)
--         REFERENCES travelite.users (id),
--     FOREIGN KEY (trek_id)
--         REFERENCES travelite.trek (id),
--     UNIQUE (user_id, trek_id)
-- );

CREATE TABLE travelite.things
(
    id   SERIAL PRIMARY KEY NOT NULL,
    name CITEXT             NOT NULL
);

CREATE TABLE travelite.trek_things
(
    trek_id  INT NOT NULL,
    thing_id INT NOT NULL,
    FOREIGN KEY (trek_id)
        REFERENCES travelite.trek (id) ON DELETE CASCADE,
    FOREIGN KEY (thing_id)
        REFERENCES travelite.things (id) ON DELETE CASCADE
);

CREATE TABLE travelite.comment
(
    id          SERIAL PRIMARY KEY NOT NULL,
    trek_id     INT                NOT NULL,
    user_id     INT                NOT NULL,
    rating      FLOAT8             NOT NULL,
    description CITEXT,
    FOREIGN KEY (user_id)
        REFERENCES travelite.users (id),
    UNIQUE (user_id, trek_id)
);

CREATE TABLE travelite.comment_photo
(
    comment_id INT  NOT NULL,
    photo_url  TEXT NOT NULL,
    FOREIGN KEY (comment_id)
        REFERENCES travelite.comment (id)
);

INSERT INTO travelite.region (name)
VALUES ('Москва и МО'),
       ('Белгородская область'),
       ('Брянская область'),
       ('Владимирская область'),
       ('Воронежская область'),
       ('Ивановская область'),
       ('Калужская область'),
       ('Костромская область'),
       ('Курская область'),
       ('Липецкая область'),
       ('Орловская область'),
       ('Рязанская область'),
       ('Смоленская область'),
       ('Тамбовская область'),
       ('Тверская область'),
       ('Тульская область'),
       ('Ярославская область'),
       ('Республика Карелия'),
       ('Республика Коми'),
       ('Архангельская область'),
       ('Ненецкий автономный округ'),
       ('Вологодская область'),
       ('Калининградская область'),
       ('Ленинградская область'),
       ('Мурманская область'),
       ('Новгородская область'),
       ('Псковская область'),
       ('Республика Адыгея'),
       ('Республика Дагестан'),
       ('Республика Ингушетия'),
       ('Кабардино-Балкарская Республика'),
       ('Республика Калмыкия'),
       ('Карачаево-Черкесская Республика'),
       ('Республика Северная Осетия - Алания'),
       ('Чеченская Республика'),
       ('Краснодарский край'),
       ('Ставропольский край'),
       ('Астраханская область'),
       ('Волгоградская область'),
       ('Ростовская область'),
       ('Республика Башкортостан'),
       ('Республика Марий Эл'),
       ('Республика Мордовия'),
       ('Республика Татарстан'),
       ('Удмуртская Республика'),
       ('Чувашская Республика'),
       ('Кировская область'),
       ('Нижегородская область'),
       ('Оренбургская область'),
       ('Пензенская область'),
       ('Пермская область'),
       ('Коми-Пермяцкий автономный округ'),
       ('Самарская область'),
       ('Саратовская область'),
       ('Ульяновская область'),
       ('Курганская область'),
       ('Свердловская область'),
       ('Тюменская область'),
       (' Ханты-Мансийский автономный округ - Югра'),
       ('Ямало-Ненецкий автономный округ'),
       ('Челябинская область'),
       ('Республика Алтай'),
       ('Республика Бурятия'),
       ('Республика Тыва'),
       ('Республика Хакасия'),
       ('Алтайский край'),
       ('Красноярский край'),
       ('Таймырский автономный округ'),
       ('Эвенкийский автономный округ'),
       ('Иркутская область'),
       ('Кемеровская область'),
       ('Новосибирская область'),
       ('Омская область'),
       ('Томская область'),
       ('Читинская область'),
       ('Агинский Бурятский автономный округ'),
       ('Республика Саха (Якутия)'),
       ('Приморский край'),
       ('Хабаровский край'),
       ('Амурская область'),
       ('Камчатская область'),
       ('Корякский автономный округ'),
       ('Магаданская область'),
       ('Сахалинская область'),
       ('Еврейская автономная область'),
       ('Чукотский автономный округ');

SELECT *
FROM travelite.region;