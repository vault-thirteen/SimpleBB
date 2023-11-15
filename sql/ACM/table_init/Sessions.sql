CREATE TABLE IF NOT EXISTS Sessions
(
    Id        serial,
    UserId    bigint UNSIGNED NOT NULL,
    StartTime datetime        NOT NULL DEFAULT NOW(),
    UserIPAB  binary(16)      NOT NULL
);
