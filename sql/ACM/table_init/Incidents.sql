CREATE TABLE IF NOT EXISTS Incidents
(
    Id       bigint UNSIGNED AUTO_INCREMENT NOT NULL,
    Time     datetime                       NOT NULL DEFAULT NOW(),
    Type     tinyint                        NOT NULL,
    Email    varchar(255)                   NOT NULL,
    UserIPAB binary(16)                              DEFAULT NULL,
    PRIMARY KEY (Id)
    -- TODO
    -- INDEX idx_Time USING BTREE (Time),
    -- INDEX idx_Type USING BTREE (Type),
    -- INDEX idx_Email USING BTREE (Email),
    -- INDEX idx_UserIPAB USING BTREE (UserIPAB)
);
