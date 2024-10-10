CREATE TABLE IF NOT EXISTS UserSubscriptions
(
    Id      bigint UNSIGNED AUTO_INCREMENT NOT NULL,
    UserId  bigint UNSIGNED                NOT NULL,
    Threads json,

    PRIMARY KEY (Id),
    INDEX idx_UserId USING BTREE (UserId)
);
