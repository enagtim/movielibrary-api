
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL CHECK(LENGTH(TRIM(username)) >= 1),
    password_hash VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK(role IN ('user', 'admin')) DEFAULT 'user'
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
    rating DECIMAL(3, 1) CHECK(rating >= 0 AND rating <= 10)
);

CREATE TABLE IF NOT EXISTS movie_actors (
    movie_id INT REFERENCES movies(id) ON DELETE CASCADE,
    actor_id INT REFERENCES actors(id) ON DELETE CASCADE,
    PRIMARY KEY(movie_id, actor_id)
);