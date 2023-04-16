DROP TABLE IF EXISTS tasks;
CREATE TABLE tasks (
    status INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL
);
