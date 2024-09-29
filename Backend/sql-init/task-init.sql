CREATE TABLE tasks
(
    id             UUID PRIMARY  KEY DEFAULT gen_random_uuid(),
    user_email     VARCHAR(255)  NOT NULL,
    name           VARCHAR(255)  NOT NULL UNIQUE,
    description    TEXT          NOT NULL,
    is_done        BOOLEAN      NOT NULL DEFAULT FALSE
)

