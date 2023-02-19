CREATE TABLE users(
    id BIGSERIAL PRIMARY KEY,
    vk_id BIGINT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE messages(
    id BIGSERIAL PRIMARY KEY,
    vk_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,

    CONSTRAINT FK_message_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE products(
    id BIGSERIAL PRIMARY KEY,
    message_id BIGINT NOT NULL,
    shop VARCHAR(255) NOT NULL,
    title TEXT NOT NULL,
    price NUMERIC(50, 2) NOT NULL,
    link TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,

    CONSTRAINT FK_product_message_id FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE
);
