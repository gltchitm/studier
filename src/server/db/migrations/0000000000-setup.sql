ALTER DATABASE studier SET timezone = 'utc';

CREATE TYPE access AS ENUM ('everyone', 'friends', 'password', 'me');

CREATE TABLE users (
    user_id CHAR(16) PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    username VARCHAR(32) NOT NULL UNIQUE,
    password CHAR(60) NOT NULL,
    verified BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE auth_tokens (
    token CHAR(256) PRIMARY KEY,
    user_id CHAR(16) NOT NULL REFERENCES users(user_id),
    expires TIMESTAMP NOT NULL DEFAULT NOW() + INTERVAL '7 DAYS'
);

CREATE TABLE verification_tokens (
    token CHAR(8) PRIMARY KEY,
    user_id CHAR(16) NOT NULL UNIQUE REFERENCES users(user_id)
);

CREATE TABLE forgot_password_tokens (
    token CHAR(24) PRIMARY KEY,
    user_id CHAR(16) NOT NULL REFERENCES users(user_id),
    expires TIMESTAMP NOT NULL DEFAULT NOW() + INTERVAL '15 MINUTES'
);

CREATE TABLE forgot_password_tickets (
    ticket CHAR(256) PRIMARY KEY,
    user_id CHAR(16) NOT NULL REFERENCES users(user_id),
    expires TIMESTAMP NOT NULL DEFAULT NOW() + INTERVAL '15 MINUTES'
);

CREATE TABLE decks (
    deck_id CHAR(16) PRIMARY KEY,
    author_id CHAR(16) NOT NULL REFERENCES users(user_id),
    name VARCHAR(32) NOT NULL,
    description VARCHAR(64) NOT NULL,
    access ACCESS NOT NULL,
    password CHAR(60)
);

CREATE TABLE flashcards (
    flashcard_id CHAR(16) PRIMARY KEY,
    deck_id CHAR(16) NOT NULL REFERENCES decks(deck_id),
    index INTEGER NOT NULL,
    term TEXT NOT NULL,
    definition TEXT NOT NULL
);

CREATE TABLE deck_tokens (
    token CHAR(256) PRIMARY KEY,
    user_id CHAR(16) NOT NULL REFERENCES users(user_id),
    deck_id CHAR(16) NOT NULL REFERENCES decks(deck_id),
    expires TIMESTAMP NOT NULL DEFAULT NOW() + INTERVAL '1 HOUR'
);

CREATE TABLE pinned_decks (
    pinned_deck_id CHAR(16) PRIMARY KEY,
    created TIMESTAMP NOT NULL DEFAULT NOW(),
    deck_id CHAR(16) NOT NULL REFERENCES decks(deck_id),
    user_id CHAR(16) NOT NULL REFERENCES users(user_id)
);

CREATE TABLE friends (
    friend_id CHAR(16) PRIMARY KEY,
    timestamp TIMESTAMP NOT NULL DEFAULT NOW(),
    from_id CHAR(16) REFERENCES users(user_id),
    to_id CHAR(16) NOT NULL REFERENCES users(user_id),
    accepted BOOL NOT NULL DEFAULT false
);
