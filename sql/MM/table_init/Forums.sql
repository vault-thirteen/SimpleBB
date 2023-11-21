CREATE TABLE IF NOT EXISTS Forums
(
    Id            bigint UNSIGNED AUTO_INCREMENT NOT NULL,
    SectionId     bigint UNSIGNED                NOT NULL,
    Name          varchar(255)                   NOT NULL,
    Threads       json,

    -- Meta data --
    CreatorUserId bigint UNSIGNED                NOT NULL,
    CreatorTime   datetime                       NOT NULL,
    EditorUserId  bigint UNSIGNED,
    EditorTime    datetime,

    PRIMARY KEY (Id)
);
