CREATE DATABASE IF NOT EXISTS mystate;
DROP TABLE IF EXISTS trunk ;	
DROP TABLE IF EXISTS participant;

CREATE TABLE IF NOT EXISTS `trunk`(
   `trunkid` INT UNSIGNED AUTO_INCREMENT,
   `participantid` INT NOT NULL,
   `name` VARCHAR(40) NOT NULL,
   `price` DECIMAL(15,2) NOT NULL,
   `description` VARCHAR(200) NOT NULL,
   PRIMARY KEY ( `trunkid` )
	)ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE  IF NOT EXISTS  `participant` (
		`participantid` INT UNSIGNED AUTO_INCREMENT,
		`name` VARCHAR(40) NOT NULL,
		`email` VARCHAR(200) NOT NULL,
		`cash` DECIMAL(15,2) NOT NULL,   
		PRIMARY KEY (participantid)
	)ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO participant (name, email, cash) VALUES (`Tom`, `Admin@example.com`, `1100.00`);
INSERT INTO participant (name, email, cash) VALUES (`Jack`, `User@example.com`, `1150.00`);
INSERT INTO trunk (participantid, name,price, description) VALUES (1,`Linux CD`, `1.00`, `Complete OS on a CD`); 
INSERT INTO trunk (participantid, name,price, description) VALUES (2,`ComputerABC`, `12.90`, `a book about OS computer!`);
INSERT INTO trunk (participantid, name,price, description) VALUES (2,`Magazines`, `6.90`, `Stack of Computer Magezines computer!`);