-- удаляем таблицы

DROP TABLE IF EXISTS movie_actors;
DROP TABLE IF EXISTS movie_directors;
DROP TABLE IF EXISTS movies;
DROP TABLE IF EXISTS studios;
DROP TABLE IF EXISTS directors;
DROP TABLE IF EXISTS actors;

-- актеры (с индексом по имени)

CREATE TABLE actors (
    id BIGSERIAL PRIMARY KEY,
    first_name TEXT NOT NULL,
	last_name TEXT NOT NULL,
	birth_date DATE NOT NULL
);

CREATE INDEX IF NOT EXISTS actors_full_name_idx ON actors USING btree (lower(last_name));

-- режиссеры (с индексом по имени)

CREATE TABLE directors (
    id BIGSERIAL PRIMARY KEY,
    first_name TEXT NOT NULL,
	last_name TEXT NOT NULL,
	birth_date DATE NOT NULL
);

CREATE INDEX IF NOT EXISTS directors_full_name_idx ON directors USING btree (lower(last_name));

-- киностудии

CREATE TABLE studios (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

-- фильмы (с индексом по названию)

CREATE TABLE movies (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
	release_year INTEGER NOT NULL CHECK (release_year >= 1800),
	studio_id BIGINT NOT NULL REFERENCES studios(id) ON DELETE CASCADE ON UPDATE CASCADE,
	gross BIGINT NOT NULL DEFAULT 0,
	rating TEXT NOT NULL CHECK (rating in ('PG-10', 'PG-13', 'PG-18'))
);

CREATE INDEX IF NOT EXISTS movies_title_idx ON movies USING btree (lower(title));

-- актеры в фильмах

CREATE TABLE movie_actors (
	id BIGSERIAL PRIMARY KEY,
	movie_id BIGINT NOT NULL REFERENCES movies(id) ON DELETE CASCADE ON UPDATE CASCADE,
	actor_id BIGINT NOT NULL REFERENCES actors(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- режиссеры в фильмах

CREATE TABLE movie_directors (
	id BIGSERIAL PRIMARY KEY,
	movie_id BIGINT NOT NULL REFERENCES movies(id) ON DELETE CASCADE ON UPDATE CASCADE,
	director_id BIGINT NOT NULL REFERENCES directors(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- функция-триггер для проверки уникальности названия фильма в рамках одного года

CREATE OR REPLACE FUNCTION check_unique_title_in_a_year()
	RETURNS TRIGGER AS $$
BEGIN
    IF (
		SELECT COUNT(*) FROM movies
		WHERE NEW.title = title AND
		NEW.release_year = release_year
	) = 0
	THEN RETURN NEW;
	ELSE RAISE EXCEPTION 'Обнаружен фильм с таким же названием и годом выпуска!';
    END IF;
END;
$$ LANGUAGE plpgsql;

-- регистрация тригера для таблицы

CREATE TRIGGER check_unique_title_in_a_year BEFORE INSERT OR UPDATE ON movies 
FOR EACH ROW EXECUTE PROCEDURE check_unique_title_in_a_year();

-- функция для возвращения списка фильмов по указанной студии

CREATE OR REPLACE FUNCTION select_movies_of_studio(st_id BIGINT)
RETURNS TABLE (id BIGINT, title TEXT, release_year INTEGER, studio_id BIGINT, gross BIGINT, rating TEXT) AS
$func$
BEGIN
   IF st_id = 0
   THEN RETURN QUERY SELECT * FROM movies;
   ELSE RETURN QUERY SELECT * FROM movies WHERE movies.studio_id = st_id;
   END IF;
END
$func$ LANGUAGE plpgsql;