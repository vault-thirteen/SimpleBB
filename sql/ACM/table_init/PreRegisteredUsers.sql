CREATE TABLE IF NOT EXISTS PreRegisteredUsers
(
    Id                 bigint UNSIGNED AUTO_INCREMENT NOT NULL,
    PreRegTime         datetime                       NOT NULL DEFAULT NOW(),
    Email              varchar(255)                   NOT NULL,
    VerificationCode   varchar(255),
    IsEmailApproved    boolean                        NOT NULL DEFAULT FALSE,
    Name               varchar(255),
    Password           varbinary(255),
    IsReadyForApproval boolean                        NOT NULL DEFAULT FALSE,
    IsApproved         boolean                        NOT NULL DEFAULT FALSE,
    ApprovalTime       datetime,
    PRIMARY KEY (Id),
    INDEX idx_PreRegTime USING BTREE (PreRegTime),
    INDEX idx_Email USING BTREE (Email),
    INDEX idx_VerificationCode USING BTREE (VerificationCode),
    INDEX idx_IsEmailApproved USING BTREE (IsEmailApproved),
    INDEX idx_Name USING BTREE (Name),
    INDEX idx_IsReadyForApproval USING BTREE (IsReadyForApproval),
    INDEX idx_IsApproved USING BTREE (IsApproved)
);
