CREATE EXTENSION "pgcrypto";

--Creation of user table
CREATE TABLE IF NOT EXISTS users (
    userId UUID NOT NULL,
    email varchar(150),
    emailVerified BOOLEAN,
    createdAt TIMESTAMP,
    updatedAt TIMESTAMP,
    passwordHash varchar(150) DEFAULT ''::character varying,
    OTP INTEGER DEFAULT 0,
    PRIMARY KEY (userId)
);

--userInfo table
CREATE TABLE IF NOT EXISTS userInfo (
    userId UUID NOT NULL,
    displayName varchar(100),
    firstName varchar(100),
    lastName varchar(100),
    photoUrl varchar(500),
    gender varchar(150),
    address varchar(150),
    age integer,
    experience varchar(150),
    playingTime integer,
    preferredPlatforms varchar(150),
     CONSTRAINT fk_userId
        FOREIGN KEY(userId)
        REFERENCES users(userId)
);

--linkedAccounts table
CREATE TABLE IF NOT EXISTS linkedAccounts (
    userId UUID NOT NULL,
    federatedId varchar NOT NULL,
    providerId varchar(100) NOT NULL,
    email varchar(150),
    linkedUserId UUID NOT NULL,
    PRIMARY KEY (federatedId, providerId),
     CONSTRAINT fk_userId
        FOREIGN KEY(userId)
        REFERENCES users(userId)
);

--credentials table
CREATE TABLE IF NOT EXISTS clients (
    id UUID NOT NULL,
    clientId text,
    clientName varchar(200),
    clientSecret text,
    projectId varchar(200),
    userId UUID,
    redirectUrl varchar(500),
    createdAt TIMESTAMP,
    updatedAt TIMESTAMP,
    active BOOLEAN
);
