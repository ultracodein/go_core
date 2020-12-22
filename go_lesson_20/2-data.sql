-- Интерстеллар

INSERT INTO studios(id, name) VALUES (1, 'Legendary Pictures');

INSERT INTO directors(id, full_name, birth_date) VALUES (1, 'Кристофер Нолан', '1970-07-30');

INSERT INTO actors(id, full_name, birth_date) VALUES (1, 'Мэттью МакКонахи', '1969-11-04');
INSERT INTO actors(id, full_name, birth_date) VALUES (2, 'Энн Хэтэуэй', '1970-07-30');
INSERT INTO actors(id, full_name, birth_date) VALUES (3, 'Джессика Честейн', '1982-11-12');

INSERT INTO movies(id, title, release_year, studio_id, gross, rating) VALUES (1, 'Интерстеллар', 2014, 1, 677463813, 'PG-13');

INSERT INTO movie_actors(id, movie_id, actor_id) VALUES (1, 1, 1);
INSERT INTO movie_actors(id, movie_id, actor_id) VALUES (2, 1, 2);
INSERT INTO movie_actors(id, movie_id, actor_id) VALUES (3, 1, 3);

INSERT INTO movie_directors(id, movie_id, director_id) VALUES (1, 1, 1);

-- Че Гевара: Дневники мотоциклиста

INSERT INTO studios(id, name) VALUES (2, 'FilmFour');

INSERT INTO directors(id, full_name, birth_date) VALUES (2, 'Уолтер Саллес', '1956-04-12');

INSERT INTO actors(id, full_name, birth_date) VALUES (4, 'Гаэль Гарсиа Берналь', '1978-11-30');
INSERT INTO actors(id, full_name, birth_date) VALUES (5, 'Родриго Де ла Серна', '1976-04-18');

INSERT INTO movies(id, title, release_year, studio_id, gross, rating) VALUES (2, 'Че Гевара: Дневники мотоциклиста', 2004, 2, 57641466, 'PG-18');

INSERT INTO movie_actors(id, movie_id, actor_id) VALUES (4, 2, 4);
INSERT INTO movie_actors(id, movie_id, actor_id) VALUES (5, 2, 5);

INSERT INTO movie_directors(id, movie_id, director_id) VALUES (2, 2, 2);

-- Облачный атлас

INSERT INTO studios(id, name) VALUES (3, 'Warner Bros. Entertainment, Inc.');

INSERT INTO directors(id, full_name, birth_date) VALUES (3, 'Лана Вачовски', '1965-06-21');
INSERT INTO directors(id, full_name, birth_date) VALUES (4, 'Том Тыквер', '1965-05-23');
INSERT INTO directors(id, full_name, birth_date) VALUES (5, 'Лилли Вачовски', '1967-12-29');

INSERT INTO actors(id, full_name, birth_date) VALUES (6, 'Том Хэнкс', '1956-07-09');
INSERT INTO actors(id, full_name, birth_date) VALUES (7, 'Холли Берри', '1966-08-14');

INSERT INTO movies(id, title, release_year, studio_id, gross, rating) VALUES (3, 'Облачный атлас', 2012, 3, 130482868, 'PG-18');

INSERT INTO movie_actors(id, movie_id, actor_id) VALUES (6, 3, 6);
INSERT INTO movie_actors(id, movie_id, actor_id) VALUES (7, 3, 7);

INSERT INTO movie_directors(id, movie_id, director_id) VALUES (3, 3, 3);
INSERT INTO movie_directors(id, movie_id, director_id) VALUES (4, 3, 4);
INSERT INTO movie_directors(id, movie_id, director_id) VALUES (5, 3, 5);