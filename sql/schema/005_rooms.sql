-- +goose Up
CREATE TABLE rooms (
    id UUID PRIMARY KEY,
    alert_id UUID UNIQUE NOT NULL REFERENCES alerts(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE rooms;