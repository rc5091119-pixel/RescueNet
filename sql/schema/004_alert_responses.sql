-- +goose Up
CREATE TABLE alert_responses(
    alert_id UUID REFERENCES alerts(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    status TEXT NOT NULL,
    responded_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (alert_id,user_id)
);

-- +goose Down
DROP TABLE alert_responses;