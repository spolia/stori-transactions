CREATE TABLE `users` (
  `alias` VARCHAR(45) NOT NULL,
  `first_name` VARCHAR(45) NOT NULL,
  `last_name` VARCHAR(45) NOT NULL,
  `email` VARCHAR(45) NOT NULL,
  `password` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`alias`),
  UNIQUE INDEX `email_UNIQUE` (`email` ASC));

  CREATE TABLE `movements` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `date` VARCHAR(45) NOT NULL,
  `type_movement` VARCHAR(45) NOT NULL,
  `amount` VARCHAR(45) NOT NULL,
  `alias` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_alias`
      FOREIGN KEY (`alias`)
          REFERENCES `users` (`alias`)
          ON DELETE CASCADE
          ON UPDATE CASCADE);


INSERT INTO `users` (`alias`, `first_name`, `last_name`, `email`, `password`)
VALUES ('jdoe', 'John', 'Doe', 'john.doe@example.com', 'password1'),
  ('jsmith', 'Jane', 'Smith', 'jane.smith@example.com', 'password2');
