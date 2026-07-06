-- Updated apikey table structure
DROP TABLE IF EXISTS `apikey`;
CREATE TABLE `apikey` (
  `key_id` varchar(255) NOT NULL,          -- Unique ID for the credential set
  `tool_name` varchar(255) NOT NULL,      -- e.g., 'Wazuh', 'VirusTotal', 'Internal'
  `ownerid` varchar(255) NOT NULL,        -- User ID
  `credential_type` varchar(50) NOT NULL, -- 'apikey', 'oauth2', 'token'
  `encrypted_data` text NOT NULL,         -- JSON blob of credentials, encrypted
  `active` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`key_id`),
  KEY `ownerid` (`ownerid`),
  CONSTRAINT `apikey_owner_fk` FOREIGN KEY (`ownerid`) REFERENCES `user` (`userid`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

ALTER TABLE `motherships` DROP FOREIGN KEY `motherships_ibfk_1`;
ALTER TABLE `motherships` DROP INDEX `ownerid`;
ALTER TABLE `motherships` ADD CONSTRAINT `motherships_owner_fk` FOREIGN KEY (`ownerid`) REFERENCES `user` (`userid`)
    ON DELETE CASCADE ON UPDATE CASCADE;

-- Add the 'active' flag and 'describers' column to the minion table
ALTER TABLE `minion`
    ADD COLUMN `active` TINYINT(1) NOT NULL DEFAULT 1 AFTER `tls`,
ADD COLUMN `describers` TEXT DEFAULT NULL AFTER `active`;

-- Fix the existing Scan errors by ensuring 'online' logic exists if needed,
-- or stick to the 'active' flag for consistency with your mothership code.

-- Represents the command intent created by the admin
CREATE TABLE IF NOT EXISTS `odin`.`rmm_tasks` (
    `task_id` VARCHAR(255) NOT NULL,
    `ownerid` VARCHAR(255) NOT NULL,
    `command` TEXT NOT NULL,
    `target_type` VARCHAR(50) NOT NULL, -- custom, linux, windows, android
    `minion_ids` TEXT DEFAULT NULL,     -- Base64 encoded list of target IDs
    `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NULL DEFAULT NULL,
    PRIMARY KEY (`task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- Represents individual verified results created by Minions
CREATE TABLE IF NOT EXISTS `odin`.`rmm_executions` (
    `execution_id` VARCHAR(255) NOT NULL,
    `task_id` VARCHAR(255) NOT NULL,
    `minionid` VARCHAR(255) NOT NULL,
    `status` VARCHAR(50) NOT NULL, -- success, failed
    `output` LONGTEXT DEFAULT NULL,
    `verification_hash` VARCHAR(255) NOT NULL,
    `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NULL DEFAULT NULL,
    PRIMARY KEY (`execution_id`),
    KEY `task_id_idx` (`task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;


-- Represents individual security flaws tied to a target
CREATE TABLE IF NOT EXISTS `odin`.`vulnerabilities` (
    `vulnerability_id` varchar(255) NOT NULL,
    `target_id` varchar(100) NOT NULL,
    `name` int(11) NOT NULL,            -- Maps to Vulnerability enum (LFI, SQLI, etc.)
    `severity` int(11) NOT NULL,
    `payload` longtext DEFAULT NULL,
    `attack_type` int(11) NOT NULL,     -- Maps to AttackType enum (PHISHING, WEBATTACK, etc.)
    `grouped` tinyint(1) DEFAULT 0,
    `authenticated` tinyint(1) DEFAULT 0,
    `works` tinyint(1) DEFAULT 0,
    `details` longtext DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT current_timestamp(),
    `updated_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`vulnerability_id`),
    KEY `target_id` (`target_id`),
    CONSTRAINT `vuln_target_fk` FOREIGN KEY (`target_id`) REFERENCES `targets` (`target_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- Represents successful exploit configurations tied to a target
CREATE TABLE IF NOT EXISTS `odin`.`exploits` (
    `exploit_id` varchar(255) NOT NULL,
    `target_id` varchar(100) NOT NULL,
    `lhost` varchar(255) DEFAULT NULL,
    `lport` int(11) DEFAULT NULL,
    `address` varchar(255) DEFAULT NULL,
    `average_severity` int(11) NOT NULL,
    `grouped` tinyint(1) DEFAULT 0,
    `grouped_vulns` text DEFAULT NULL, -- Stored as comma-separated vulnerability_ids
    `works` tinyint(1) DEFAULT 0,
    `created_at` timestamp NULL DEFAULT current_timestamp(),
    `updated_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`exploit_id`),
    KEY `target_id` (`target_id`),
    CONSTRAINT `exploit_target_fk` FOREIGN KEY (`target_id`) REFERENCES `targets` (`target_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;


