-- tables
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
    nickname citext PRIMARY KEY,
    fullname varchar(128) NOT NULL,
    about citext NOT NULL,
    email citext UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS forums (
    title varchar(128) NOT NULL,
    user_nickname citext NOT NULL REFERENCES users(nickname),
    slug citext PRIMARY KEY,
    posts int DEFAULT 0,
    threads int DEFAULT 0
);

CREATE TABLE IF NOT EXISTS threads (
    id serial PRIMARY KEY,
    title varchar(128) NOT NULL,
    author citext NOT NULL REFERENCES users(nickname),
    forum citext NOT NULL REFERENCES forums(slug) ON DELETE CASCADE,
    message text NOT NULL,
    votes int DEFAULT 0,
    slug citext UNIQUE,
    created timestamptz DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS posts (
    id serial PRIMARY KEY,
    parent int,
    author citext NOT NULL REFERENCES users(nickname),
    message TEXT NOT NULL,
    is_edited BOOLEAN NOT NULL,
    forum citext REFERENCES forums(slug) ON DELETE CASCADE,
    thread int REFERENCES threads(id) ON DELETE CASCADE,
    created timestamptz,
    post_tree int[] DEFAULT ARRAY []::integer[]
);

CREATE TABLE IF NOT EXISTS votes (
    thread int NOT NULL REFERENCES threads(id) ON DELETE CASCADE,
    nickname citext NOT NULL REFERENCES users(nickname),
    voice int NOT NULL,
    PRIMARY KEY (thread, nickname)
);

CREATE TABLE IF NOT EXISTS forum_user (
    user_nickname citext NOT NULL REFERENCES users(nickname) ON DELETE CASCADE,
    forum citext NOT NULL REFERENCES forums(slug) ON DELETE CASCADE,
    PRIMARY KEY (user_nickname, forum)
);

-- triggers
CREATE OR REPLACE FUNCTION update_thread_votes_after_insert()
RETURNS TRIGGER AS $$
    BEGIN
    UPDATE threads
    SET votes = votes + NEW.voice
    WHERE id = NEW.thread;
    RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insert_vote_trigger
AFTER INSERT ON votes
FOR EACH ROW
EXECUTE PROCEDURE update_thread_votes_after_insert();


CREATE OR REPLACE FUNCTION update_thread_votes_after_update()
RETURNS TRIGGER AS $$
    BEGIN
    UPDATE threads
    SET votes = votes + NEW.voice - OLD.voice
    WHERE id = NEW.thread;
    RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_vote_trigger
AFTER UPDATE ON votes
FOR EACH ROW
EXECUTE PROCEDURE update_thread_votes_after_update();


CREATE OR REPLACE FUNCTION update_post_tree()
RETURNS TRIGGER AS $$
    BEGIN
        NEW.post_tree =
            (SELECT post_tree
            FROM posts
            WHERE id = NEW.parent)
        || NEW.id;
    RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insert_post_trigger_post
BEFORE INSERT ON posts
FOR EACH ROW
EXECUTE PROCEDURE update_post_tree();


CREATE OR REPLACE FUNCTION update_count_threads()
RETURNS TRIGGER AS $$
    BEGIN
    UPDATE forums SET threads = threads + 1 WHERE slug = NEW.forum;
    RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insert_thread_trigger
AFTER INSERT ON threads
FOR EACH ROW
EXECUTE PROCEDURE update_count_threads();


CREATE OR REPLACE FUNCTION update_count_posts()
RETURNS TRIGGER AS $$
    BEGIN
    UPDATE forums SET posts = posts + 1 WHERE slug = NEW.forum;
    RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insert_post_trigger_forum
AFTER INSERT ON posts
FOR EACH ROW
EXECUTE PROCEDURE update_count_posts();

-- indexes
CREATE INDEX IF NOT EXISTS forums_user_nickname ON forums (user_nickname);

CREATE INDEX IF NOT EXISTS threads_author ON threads (author);
CREATE INDEX IF NOT EXISTS threads_forum ON threads (forum);

CREATE INDEX IF NOT EXISTS forum_user_forum_user_nickname ON forum_user (forum, user_nickname);

CREATE INDEX IF NOT EXISTS posts_thread_id on posts (thread, id);
CREATE INDEX IF NOT EXISTS posts_thread_post_tree on posts (thread, post_tree);
CREATE INDEX IF NOT EXISTS posts_parent_thread_id on posts (parent, thread, id);
CREATE INDEX IF NOT EXISTS posts_post_tree_one_post_tree on posts ((post_tree[1]), post_tree);

CREATE INDEX IF NOT EXISTS users_email ON users (email);
CREATE INDEX IF NOT EXISTS users_email_nickname ON users (email, nickname);

CREATE INDEX IF NOT EXISTS index_threads_slug on threads (slug);
CREATE INDEX IF NOT EXISTS index_threads_forum_created ON threads (forum, created);

VACUUM ANALYZE;