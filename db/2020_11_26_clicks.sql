CREATE TABLE clicks(
    id INT NOT NULL AUTO_INCREMENT,
    screenX INT NOT NULL,
    screenY INT NOT NULL,
    captureCoordinateX INT NOT NULL,
    captureCoordinateY INT NOT NULL,
    capturedTime VARCHAR(255) NOT NULL,
    capturedDay VARCHAR(25) NOT NULL,
    PRIMARY KEY(id)
);