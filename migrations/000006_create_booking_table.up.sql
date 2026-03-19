CREATE TABLE IF NOT EXISTS booking(
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    computer_id INT NOT NULL,
    time_start TIMESTAMP NOT NULL,
    time_end TIMESTAMP NOT NULL,

    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_computer FOREIGN KEY(computer_id) REFERENCES computers(id)
);

CREATE INDEX idx_booking_computer_time ON booking (computer_id, time_start, time_end);