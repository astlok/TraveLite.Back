ALTER
    USER postgres WITH ENCRYPTED PASSWORD 'admin';

DROP SCHEMA IF EXISTS travelite CASCADE;

CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE SCHEMA IF NOT EXISTS travelite;

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
CREATE TYPE travelite.moderate_status AS ENUM ('no status', 'pending', 'failed', 'verified');
CREATE TABLE travelite.trek
(
    id              SERIAL PRIMARY KEY           NOT NULL,
    name            CITEXT                       NOT NULL,
    difficult       INT                          NOT NULL,
    days            INT                          NOT NULL,
    description     CITEXT                    DEFAULT '',
    best_time_to_go TRAVELITE.SEASONS         DEFAULT '',
    type            TRAVELITE.HIKE_TYPE          NOT NULL,
    climb           INT                          NOT NULL,
    region          TEXT                         NOT NULL,
    creator_id      INT                          NOT NULL,
    mod_status      TRAVELITE.MODERATE_STATUS DEFAULT 'no status',
    distance        INT                          NOT NULL,
    route           GEOGRAPHY(LINESTRINGZ, 4326) NOT NULL,
    start           GEOGRAPHY(POINTZ, 4326)      NOT NULL,
    rate            INT                       DEFAULT 0,
    FOREIGN KEY (creator_id)
        REFERENCES travelite.users (id)
);

CREATE INDEX start_trek_idx ON travelite.trek USING GIST (start);

CREATE OR REPLACE FUNCTION travelite.insert_random_rate() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE travelite.trek SET rate=Cast(random() * 5 as int) WHERE id = NEW.id;
    RETURN NEW;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER rand_rate
    AFTER INSERT
    ON travelite.trek
    FOR EACH ROW
EXECUTE FUNCTION travelite.insert_random_rate();

CREATE TABLE travelite.marks
(
    trek_id     INT                     NOT NULL,
    point       GEOGRAPHY(POINTZ, 4326) NOT NULL,
    title       TEXT                    NOT NULL,
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

CREATE TYPE travelite.file_owner AS ENUM ('route', 'comment', 'user');

CREATE TABLE travelite.files
(
    id       UUID PRIMARY KEY     NOT NULL,
    filename TEXT                 NOT NULL,
    owner    travelite.file_owner NOT NULL,
    owner_id BIGINT               NOT NULL,
    link     text default ''

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

-- SELECT *
-- FROM travelite.region;
--
-- SELECT
--     ST_Intersects(
--             ST_GeogFromText('SRID=4326;POLYGON Z((55.745359 37.658375 1230481234, 55.745526 37.705746 129341234, 55.724144 37.709792 7865823796, 55.723866 37.627189 12348124, 55.745359 37.658375 1238947812))'),
--             ST_GeogFromText('SRID=4326;POINT Z(55.737454 37.674165 19234129034781298743)'));
--
-- SELECT id, name, difficult, days, description, best_time_to_go, type, climb, region, creator_id, is_moderate, ST_AsText(start) AS START from travelite.trek
-- WHERE ST_Intersects(
--               ST_GeogFromText('SRID=4326;POLYGON ((55.745359 37.658375, 55.745526 37.705746, 55.724144 37.709792, 55.723866 37.627189, 55.745359 37.658375))'),
--               START);
--
-- SELECT 'SRID=4326;POINT (55.745359 37.658375)'::geometry;
