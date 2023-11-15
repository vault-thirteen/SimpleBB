CREATE TABLE IF NOT EXISTS Threads
(
    Id            serial,
    ForumId       bigint UNSIGNED NOT NULL,
    Name          varchar(255)    NOT NULL,
    Messages      json,
    CreatorUserId bigint UNSIGNED NOT NULL,
    CreatorTime   datetime        NOT NULL,
    EditorUserId  bigint UNSIGNED,
    EditorTime    datetime
);
