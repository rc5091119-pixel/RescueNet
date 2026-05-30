-- +goose Up
CREATE TABLE room_members (
    room_id UUID NOT NULL REFERENCES rooms(id),
    user_id UUID NOT NULL REFERENCES users(id),

    PRIMARY KEY(room_id, user_id)
);

-- +goose Down
DROP TABLE room_members;