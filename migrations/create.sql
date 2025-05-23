
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL CHECK(LENGTH(TRIM(username)) >= 1),
    password_hash VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user' CHECK(role IN ('user', 'admin')) 
);

CREATE TABLE IF NOT EXISTS actors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL CHECK(LENGTH(TRIM(name)) >= 1),
    gender VARCHAR(20) NOT NULL CHECK(gender IN ('male', 'female', 'other')),
    birth_date DATE NOT NULL
);

CREATE TABLE IF NOT EXISTS movies (
    id SERIAL PRIMARY KEY,
    title VARCHAR(150) NOT NULL CHECK(LENGTH(TRIM(title)) BETWEEN 1 AND 150),
    description VARCHAR(1000),
    release_date DATE NOT NULL,
    rating DECIMAL(3, 1) NOT NULL DEFAULT 0 CHECK(rating >= 0 AND rating <= 10) 
);

CREATE TABLE IF NOT EXISTS movie_actors (
    movie_id INT REFERENCES movies(id) ON DELETE CASCADE,
    actor_id INT REFERENCES actors(id) ON DELETE CASCADE,
    PRIMARY KEY(movie_id, actor_id)
);

CREATE INDEX IF NOT EXISTS idx_movies_title ON movies(title);
CREATE INDEX IF NOT EXISTS idx_movies_rating ON movies(rating DESC);
CREATE INDEX IF NOT EXISTS idx_release_date ON movies(release_date DESC);
CREATE INDEX IF NOT EXISTS idx_actors_name ON actors(name);