-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.users (
    username varchar(32) not null primary key
);

CREATE TABLE IF NOT EXISTS public.movies (
    title varchar(50) not null primary key
);

CREATE TABLE IF NOT EXISTS public.reviews (
    id serial primary key,
    reviewer varchar(32)
        references public.users(username)
        on delete cascade,
    movie_title varchar(50)
        references public.movies(title)
        on delete cascade,
    text varchar(200) not null,
    rating int not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.reviews;
DROP TABLE IF EXISTS public.movies;
DROP TABLE IF EXISTS public.users;
-- +goose StatementEnd
