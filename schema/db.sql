CREATE DATABASE ergo;
USE ergo;

CREATE TABLE transactions (
  id                INT(11) NOT NULL AUTO_INCREMENT,
  tx_id             VARCHAR(255) NOT NULL,
  blockNumber       VARCHAR(255) NOT NULL,
  `to`              VARCHAR(255) NOT NULL,
  amount            VARCHAR(255) NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY tx_id (tx_id)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;
 
--tx_id = hash
-- `from`            VARCHAR(255) NOT NULL,

CREATE TABLE blocks (
  lastUpdateTime        INT(11) NOT NULL,
  lastUpdatedBlockNum   VARCHAR(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- INSERT INTO blocks (lastUpdateTime, lastUpdatedBlockNum) VALUES (1642670451, '11215055');

CREATE TABLE addresses (
  id            INT(11) NOT NULL AUTO_INCREMENT,
  created       INT(11) NOT NULL,
  address       VARCHAR(255) NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY address (address)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

-- privateKey    VARCHAR(255) NOT NULL,

# transactions indexes
CREATE INDEX idx_txto     ON transactions (`to`);
CREATE INDEX idx_txamount ON transactions (amount);
CREATE INDEX idx_txhash   ON transactions (tx_id);

# addresses indexes
CREATE INDEX idx_addraddr ON addresses (address);

