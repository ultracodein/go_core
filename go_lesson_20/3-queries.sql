-- средний возраст актера

SELECT DATE_TRUNC('year', AVG(AGE(now(), birth_date))) FROM actors;

-- фильмы с наибольшими сборами на актера

SELECT title, gross/(SELECT COUNT(*) FROM movie_actors WHERE movie_id = movies.id) as gross_per_actor FROM movies ORDER BY gross_per_actor DESC;

-- самые кассовые режиссеры

SELECT full_name FROM directors
INNER JOIN movie_directors ON directors.id = movie_directors.director_id
INNER JOIN movies ON movie_directors.movie_id = movies.id
WHERE movies.gross > 500000000;