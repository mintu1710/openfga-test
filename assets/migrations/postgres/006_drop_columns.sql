-- +goose Up
--  This migration drops columns `_user` and `user_type` and triggers

ALTER TABLE tuple
     DROP COLUMN _user,
     DROP COLUMN user_type;
-- Index idx_reverse_lookup_user will be automatically dropped

ALTER TABLE changelog DROP COLUMN _user;

DROP TRIGGER IF EXISTS migrate_user_column ON tuple;

DROP TRIGGER IF EXISTS migrate_user_column ON changelog;

-- +goose Down

ALTER TABLE tuple
     ADD COLUMN user_type TEXT,
     ADD COLUMN _user TEXT;

-- (..., jon, ...) -> user
-- (user, *, ...) -> userset
-- (user, jon, ...) -> user
-- (group, eng, member) -> userset
UPDATE tuple SET user_type = (CASE
                                   WHEN user_relation <> '...' THEN 'userset'
                                   WHEN user_object_id = '*' THEN 'userset'
                                   ELSE 'user'
     END),
-- (..., jon, ...) becomes jon
-- (user, *, ...) becomes user:*
-- (user, jon, ...) becomes user:jon
-- (group, eng, member) becomes group:eng#member
                  _user = (CASE
                               WHEN user_object_type = '...' THEN user_object_id
                               WHEN user_relation = '...' THEN user_object_type || ':' || user_object_id
                               ELSE user_object_type || ':' || user_object_id || '#' || user_relation
                      END);

ALTER TABLE changelog ADD COLUMN _user TEXT;

UPDATE changelog SET _user = (CASE
                                   WHEN user_object_type = '...' THEN user_object_id
                                   WHEN user_relation = '...' THEN user_object_type || ':' || user_object_id
                                   ELSE user_object_type || ':' || user_object_id || '#' || user_relation
     END);