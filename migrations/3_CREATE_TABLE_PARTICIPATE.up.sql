CREATE TABLE participate (
    id              BIGSERIAL PRIMARY KEY,
    account         VARCHAR(50) NOT NULL,
    mail            VARCHAR(50) NOT NULL,
    token           UUID NOT NULL DEFAULT uuid_generate_v4(),
    poll            INT NOT NULL,
    participate     BOOLEAN NOT NULL
);
