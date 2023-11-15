CREATE TABLE IF NOT EXISTS PreRegisteredUsers
(
    Id                 serial,
    PreRegTime         datetime     NOT NULL DEFAULT NOW(),
    Email              varchar(255) NOT NULL,
    VerificationCode   varchar(255),
    IsEmailApproved    boolean      NOT NULL DEFAULT FALSE,
    Name               varchar(255),
    Password           varbinary(255),
    IsReadyForApproval boolean      NOT NULL DEFAULT FALSE,
    IsApproved         boolean      NOT NULL DEFAULT FALSE,
    ApprovalTime       datetime
);
