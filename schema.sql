CREATE TABLE alert (
                       id SERIAL PRIMARY KEY,
                       event_time TIMESTAMP NOT NULL,
                       received_at TIMESTAMP NOT NULL,
                       notification_sent BOOLEAN NOT NULL DEFAULT FALSE,
                       device_id VARCHAR NOT NULL REFERENCES device(id),
                       face_model_id VARCHAR NOT NULL REFERENCES face_detection_model(id),
                       mask_model_id VARCHAR NOT NULL REFERENCES mask_classification_model(id),
                       image_id INT NOT NULL REFERENCES image(id)
);

CREATE TABLE device (
                        id VARCHAR PRIMARY KEY,
                        type VARCHAR NOT NULL,
                        enrolled_on TIMESTAMP NOT NULL,
                        latitude REAL NOT NULL,
                        longitude REAL NOT NULL
);

CREATE TABLE image (
                       id SERIAL PRIMARY KEY,
                       format VARCHAR,
                       width INT,
                       height INT,
                       data bytea
);

CREATE TABLE face_detection_model (
                                      id VARCHAR PRIMARY KEY,
                                      created_at TIMESTAMP NOT NULL,
                                      name VARCHAR NOT NULL,
                                      threshold REAL NOT NULL
);

CREATE TABLE mask_classification_model (
                                           id VARCHAR PRIMARY KEY,
                                           created_at TIMESTAMP NOT NULL,
                                           name VARCHAR NOT NULL,
                                           threshold REAL NOT NULL
);