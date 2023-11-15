CREATE TABLE IF NOT EXISTS PreSessions
(
    Id                   serial,
    UserId               bigint UNSIGNED NOT NULL,
    TimeOfCreation       datetime        NOT NULL DEFAULT NOW(),
    RequestId            varchar(255)    NOT NULL,
    UserIPAB             binary(16)      NOT NULL,
    AuthDataBytes        varbinary(1024) NOT NULL,
    IsCaptchaRequired    boolean         NOT NULL,
    CaptchaId            varchar(255), -- Null when not needed --
    IsVerifiedByCaptcha  boolean,      -- Null when not needed --
    IsVerifiedByPassword boolean         NOT NULL DEFAULT FALSE,
    VerificationCode     varchar(255)             DEFAULT NULL,
    IsVerifiedByEmail    boolean         NOT NULL DEFAULT FALSE
);
