CREATE TABLE IF NOT EXISTS ThreadSubscriptions
(
    Id       bigint UNSIGNED AUTO_INCREMENT NOT NULL,
    ThreadId bigint UNSIGNED                NOT NULL,
    Users    json,

    PRIMARY KEY (Id),
    INDEX idx_ThreadId USING BTREE (ThreadId)
);
