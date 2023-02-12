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

CREATE TABLE results(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    message_id BIGINT NOT NULL,
    shop VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    price DECIMAL(50,2) NOT NULL,
    link TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,

    CONSTRAINT FK_result_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT FK_result_message_id FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE
);