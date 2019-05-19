CREATE TABLE user (
  id BIGINT UNSIGNED AUTO_INCREMENT NOT NULL PRIMARY KEY,
  name VARCHAR(191) UNIQUE,
  salt VARCHAR(20),
  password VARCHAR(40),
  display_name TEXT,
  avatar_icon TEXT,
  created_at DATETIME NOT NULL,
  KEY name_idx_on_user (name)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE image (
  id BIGINT UNSIGNED AUTO_INCREMENT NOT NULL PRIMARY KEY,
  name VARCHAR(191),
  data LONGBLOB,
  KEY name_idx_on_image (name)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE channel (
  id BIGINT AUTO_INCREMENT NOT NULL PRIMARY KEY,
  name TEXT NOT NULL,
  description MEDIUMTEXT,
  updated_at DATETIME NOT NULL,
  created_at DATETIME NOT NULL
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE message (
  id BIGINT AUTO_INCREMENT NOT NULL PRIMARY KEY,
  channel_id BIGINT,
  user_id BIGINT,
  content TEXT,
  created_at DATETIME NOT NULL,
  KEY user_id_idx_on_message (user_id),
  KEY created_at_idx_on_message (created_at),
  KEY channel_id_idx_on_message (channel_id)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE haveread (
  user_id BIGINT NOT NULL,
  channel_id BIGINT NOT NULL,
  message_id BIGINT,
  updated_at DATETIME NOT NULL,
  created_at DATETIME NOT NULL,
  KEY user_id_idx_on_haveread (user_id),
  KEY channel_id_idx_on_haveread (channel_id),
  PRIMARY KEY(user_id, channel_id)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4;
