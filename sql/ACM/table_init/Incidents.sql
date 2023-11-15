CREATE TABLE IF NOT EXISTS Incidents
(
    Id       serial,
    Time     datetime     NOT NULL,
    Type     tinyint      NOT NULL,
    Email    varchar(255) NOT NULL,
    UserIPAB binary(16) DEFAULT NULL
);
