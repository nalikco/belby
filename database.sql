CREATE TABLE users(
    id BIGSERIAL PRIMARY KEY,
    vk_id BIGINT NULL,
    created_at TIMESTAMP NOT NUll
);

CREATE TABLE messages(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP NOT NUll,

    CONSTRAINT FK_message_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);