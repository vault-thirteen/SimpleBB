CREATE TABLE IF NOT EXISTS Sections
(
    Id            bigint UNSIGNED AUTO_INCREMENT NOT NULL,
    Parent        bigint UNSIGNED,
    ChildType     tinyint UNSIGNED DEFAULT 0,
    Children      json,
    Name          varchar(255)                   NOT NULL,

    -- Meta data --
    CreatorUserId bigint UNSIGNED                NOT NULL,
    CreatorTime   datetime                       NOT NULL,
    EditorUserId  bigint UNSIGNED,
    EditorTime    datetime,

    PRIMARY KEY (Id),
    INDEX idx_Parent USING BTREE (Parent)
);
