CREATE TABLE IF NOT EXISTS Messages
(
    Id            bigint UNSIGNED AUTO_INCREMENT NOT NULL, -- 8B --
    ThreadId      bigint UNSIGNED                NOT NULL, -- 8B --
    Text          varchar(16368)                 NOT NULL,
    TextChecksum  int UNSIGNED                   NOT NULL, -- 4B --

    -- Meta data --
    CreatorUserId bigint UNSIGNED                NOT NULL, -- 8B --
    CreatorTime   datetime                       NOT NULL, -- 8B --
    EditorUserId  bigint UNSIGNED,                         -- 8B --
    EditorTime    datetime,                                -- 8B --

    PRIMARY KEY (Id)
);
