CREATE DATABASE ergo;
USE ergo;

CREATE TABLE transactions (
  id                INT(11) NOT NULL AUTO_INCREMENT,
  hash              VARCHAR(255) NOT NULL,
  blockNumber       VARCHAR(255) NOT NULL,
  `from`       VARCHAR(255) NOT NULL,
  `to`              VARCHAR(255) NOT NULL,
  amount            VARCHAR(255) NOT NULL,
  numConfirmation   VARCHAR(255) NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;
 

CREATE TABLE blocks (
  lastUpdateTime        INT(11) NOT NULL,
  lastUpdatedBlockNum   VARCHAR(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

INSERT INTO blocks (lastUpdateTime, lastUpdatedBlockNum) VALUES (1642670451, '720494');

CREATE TABLE addresses (
  id            INT(11) NOT NULL AUTO_INCREMENT,
  created       INT(11) NOT NULL,
  address       VARCHAR(255) NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY address (address)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;


# transactions indexes
CREATE INDEX idx_txto     ON transactions (`to`);
CREATE INDEX idx_txamount ON transactions (amount);
CREATE INDEX idx_txhash   ON transactions (hash);

# addresses indexes
CREATE INDEX idx_addraddr ON addresses (address);

INSERT INTO `ergo`.`addresses`
(`created`,
`address`)
VALUES
(1647970877,"3Wwmx2AoP5MpLApcpHkdpjdscFcBStSiYo2jQsyYdEc2MaoGaU59"),
(1647970877,"3Wy4Bpi1YD74w1CfFNxpoKeGNPSHdQgkz3qxeNuAVahxCjDFQ5Bb");

3Wwmx2AoP5MpLApcpHkdpjdscFcBStSiYo2jQsyYdEc2MaoGaU59
3Wy4Bpi1YD74w1CfFNxpoKeGNPSHdQgkz3qxeNuAVahxCjDFQ5Bb

3Wwmx2AoP5MpLApcpHkdpjdscFcBStSiYo2jQsyYdEc2MaoGaU59
3Wy4Bpi1YD74w1CfFNxpoKeGNPSHdQgkz3qxeNuAVahxCjDFQ5Bb