CREATE TABLE users(
    id INTEGER NOT NULL AUTO_INCREMENT,
    userid VARCHAR(100) NOT NULL,
    username VARCHAR(50) NOT NULL,
    user_home_directory VARCHAR(100) NOT NULL,
    PRIMARY KEY(id)
);