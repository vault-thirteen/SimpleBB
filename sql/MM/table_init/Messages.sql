CREATE TABLE IF NOT EXISTS Messages
(
    Id            serial,                   -- 8B --
    ThreadId      bigint UNSIGNED NOT NULL, -- 8B --
    Text          varchar(16371)  NOT NULL,

    -- Meta data --
    CreatorUserId bigint UNSIGNED NOT NULL, -- 8B --
    CreatorTime   datetime        NOT NULL, -- 8B --
    EditorUserId  bigint UNSIGNED,          -- 8B --
    EditorTime    datetime                  -- 8B --
);
