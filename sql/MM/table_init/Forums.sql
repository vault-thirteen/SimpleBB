CREATE TABLE IF NOT EXISTS Forums
(
    Id            serial,
    Parent        bigint UNSIGNED,
    Children      json,
    Name          varchar(255)    NOT NULL,
    Threads       json,
    CreatorUserId bigint UNSIGNED NOT NULL,
    CreatorTime   datetime        NOT NULL,
    EditorUserId  bigint UNSIGNED,
    EditorTime    datetime
);
