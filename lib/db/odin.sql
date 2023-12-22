-- phpMyAdmin SQL Dump
-- version 5.2.0
-- https://www.phpmyadmin.net/
--
-- Host: localhost
-- Generation Time: Nov 19, 2023 at 12:12 PM
-- Server version: 10.4.24-MariaDB
-- PHP Version: 8.1.6

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `odin`
--

-- --------------------------------------------------------

--
-- Table structure for table `apikey`
--

CREATE TABLE `apikey` (
  `apikey` varchar(255) NOT NULL,
  `ownerid` varchar(255) NOT NULL,
  `active` tinyint(1) NOT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `appointments`
--

CREATE TABLE `appointments` (
  `userid` varchar(255) NOT NULL,
  `appointmentid` varchar(25) NOT NULL,
  `title` varchar(25) NOT NULL,
  `description` varchar(255) NOT NULL,
  `done` tinyint(1) NOT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `apt`
--

CREATE TABLE `apt` (
  `aptname` varchar(255) NOT NULL,
  `code` int(255) NOT NULL,
  `aptid` int(11) NOT NULL,
  `description` text DEFAULT NULL,
  `active` tinyint(1) NOT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `assets`
--

CREATE TABLE `assets` (
  `asset_id` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` text NOT NULL,
  `describers` text NOT NULL,
  `active` tinyint(1) NOT NULL,
  `hardware` tinyint(1) NOT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `events`
--

CREATE TABLE `events` (
  `id` int(11) NOT NULL,
  `event_id` varchar(255) NOT NULL,
  `operating_system` int(11) NOT NULL,
  `handled` tinyint(1) NOT NULL,
  `level` int(11) NOT NULL,
  `data` text NOT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `hashes`
--

CREATE TABLE `hashes` (
  `userid` varchar(255) NOT NULL,
  `hash` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `ioc`
--

CREATE TABLE `ioc` (
  `ioc_id` int(11) NOT NULL,
  `type` varchar(255) NOT NULL,
  `value` text NOT NULL,
  `source` text NOT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL,
  `virusid` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `minion`
--

CREATE TABLE `minion` (
  `minionid` varchar(255) NOT NULL,
  `name` char(100) DEFAULT NULL,
  `uname` varchar(255) DEFAULT NULL,
  `userid` varchar(255) DEFAULT NULL,
  `groupid` varchar(255) DEFAULT NULL,
  `homedir` varchar(255) DEFAULT NULL,
  `ostype` char(20) NOT NULL,
  `description` text NOT NULL,
  `installed` tinyint(1) NOT NULL,
  `mothershipid` varchar(255) NOT NULL,
  `address` text NOT NULL,
  `motherships` text DEFAULT NULL,
  `tunnel_address` text NOT NULL,
  `tls` tinyint(1) NOT NULL,
  `ownerid` varchar(255) NOT NULL,
  `lastseen` varchar(255) NOT NULL,
  `is_dropper` tinyint(1) NOT NULL,
  `generate_command` text NOT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `motherships`
--

CREATE TABLE `motherships` (
  `ownerid` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `msid` varchar(255) NOT NULL,
  `address` text NOT NULL,
  `implant_tunnel` text NOT NULL,
  `admin_tunnel` text NOT NULL,
  `other_motherships` text NOT NULL,
  `description` text NOT NULL,
  `tls` tinyint(1) NOT NULL,
  `certpem` text NOT NULL,
  `keypem` text NOT NULL,
  `active` tinyint(1) NOT NULL,
  `generate_command` text NOT NULL,
  `machine_data` text NOT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `outputs`
--

CREATE TABLE `outputs` (
  `output_id` int(11) NOT NULL,
  `service_id` int(11) DEFAULT NULL,
  `command` text DEFAULT NULL,
  `output` text DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `plugins`
--

CREATE TABLE `plugins` (
  `owner` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `hash` varchar(255) NOT NULL,
  `plugin_type` int(11) NOT NULL,
  `description` text NOT NULL,
  `active` tinyint(1) NOT NULL,
  `signed` tinyint(1) NOT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `scans`
--

CREATE TABLE `scans` (
  `scan_id` varchar(100) NOT NULL,
  `name` varchar(255) NOT NULL,
  `scan_type` enum('Bug Bounty','Pentest','Black Ops') DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `services`
--

CREATE TABLE `services` (
  `service_id` int(11) NOT NULL,
  `target_id` varchar(100) NOT NULL,
  `service_name` varchar(255) NOT NULL,
  `port` int(11) DEFAULT NULL,
  `protocol` varchar(20) NOT NULL,
  `state` tinyint(1) NOT NULL,
  `version` varchar(255) DEFAULT NULL,
  `AT` int(11) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `targets`
--

CREATE TABLE `targets` (
  `target_id` varchar(100) NOT NULL,
  `scan_id` varchar(100) NOT NULL,
  `host` varchar(255) NOT NULL,
  `host_ip` int(10) UNSIGNED NOT NULL,
  `target_ip` int(10) UNSIGNED NOT NULL,
  `firewall_name` varchar(25) NOT NULL DEFAULT 'NONE',
  `decoys` text NOT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `user`
--

CREATE TABLE `user` (
  `userid` varchar(255) NOT NULL,
  `ownerid` varchar(255) NOT NULL,
  `username` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `active` tinyint(1) NOT NULL,
  `anonymous` tinyint(1) NOT NULL,
  `verified` tinyint(1) NOT NULL,
  `admin` tinyint(1) NOT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `virus`
--

CREATE TABLE `virus` (
  `aptid` int(11) NOT NULL,
  `virusid` varchar(255) NOT NULL,
  `hash` text NOT NULL,
  `virustype` varchar(255) NOT NULL,
  `filetype` varchar(255) NOT NULL,
  `communicationmode` varchar(255) NOT NULL,
  `ostype` varchar(255) NOT NULL,
  `description` text NOT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `webdata`
--

CREATE TABLE `webdata` (
  `target_id` varchar(100) NOT NULL,
  `webdata` int(11) NOT NULL,
  `directory_path` text DEFAULT NULL,
  `parameter_path` text DEFAULT NULL,
  `file_path` text DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `yara_rule`
--

CREATE TABLE `yara_rule` (
  `yr_id` int(11) NOT NULL,
  `ioc_id` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `meta` text NOT NULL,
  `condition` text NOT NULL,
  `actions` text NOT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Indexes for dumped tables
--

--
-- Indexes for table `apikey`
--
ALTER TABLE `apikey`
  ADD PRIMARY KEY (`apikey`),
  ADD KEY `ownerid` (`ownerid`);

--
-- Indexes for table `appointments`
--
ALTER TABLE `appointments`
  ADD PRIMARY KEY (`appointmentid`),
  ADD KEY `userid` (`userid`);

--
-- Indexes for table `apt`
--
ALTER TABLE `apt`
  ADD PRIMARY KEY (`aptid`);

--
-- Indexes for table `assets`
--
ALTER TABLE `assets`
  ADD PRIMARY KEY (`asset_id`);

--
-- Indexes for table `events`
--
ALTER TABLE `events`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `hashes`
--
ALTER TABLE `hashes`
  ADD PRIMARY KEY (`userid`),
  ADD UNIQUE KEY `hash` (`hash`);

--
-- Indexes for table `ioc`
--
ALTER TABLE `ioc`
  ADD PRIMARY KEY (`ioc_id`);

--
-- Indexes for table `minion`
--
ALTER TABLE `minion`
  ADD PRIMARY KEY (`minionid`),
  ADD UNIQUE KEY `minionid` (`minionid`),
  ADD KEY `userid` (`userid`),
  ADD KEY `mothershipid` (`mothershipid`),
  ADD KEY `ownerid` (`ownerid`);

--
-- Indexes for table `motherships`
--
ALTER TABLE `motherships`
  ADD PRIMARY KEY (`msid`),
  ADD UNIQUE KEY `ownerid` (`ownerid`),
  ADD UNIQUE KEY `msid` (`msid`);

--
-- Indexes for table `outputs`
--
ALTER TABLE `outputs`
  ADD PRIMARY KEY (`output_id`),
  ADD KEY `service_id` (`service_id`);

--
-- Indexes for table `plugins`
--
ALTER TABLE `plugins`
  ADD PRIMARY KEY (`hash`),
  ADD KEY `owner` (`owner`);

--
-- Indexes for table `services`
--
ALTER TABLE `services`
  ADD PRIMARY KEY (`service_id`);

--
-- Indexes for table `user`
--
ALTER TABLE `user`
  ADD PRIMARY KEY (`userid`),
  ADD UNIQUE KEY `userid` (`userid`);

--
-- Indexes for table `virus`
--
ALTER TABLE `virus`
  ADD PRIMARY KEY (`aptid`);

--
-- Indexes for table `webdata`
--
ALTER TABLE `webdata`
  ADD PRIMARY KEY (`webdata`);

--
-- Indexes for table `yara_rule`
--
ALTER TABLE `yara_rule`
  ADD PRIMARY KEY (`yr_id`),
  ADD KEY `ioc_id` (`ioc_id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `apt`
--
ALTER TABLE `apt`
  MODIFY `aptid` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `events`
--
ALTER TABLE `events`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ioc`
--
ALTER TABLE `ioc`
  MODIFY `ioc_id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `outputs`
--
ALTER TABLE `outputs`
  MODIFY `output_id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `services`
--
ALTER TABLE `services`
  MODIFY `service_id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `virus`
--
ALTER TABLE `virus`
  MODIFY `aptid` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `webdata`
--
ALTER TABLE `webdata`
  MODIFY `webdata` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `yara_rule`
--
ALTER TABLE `yara_rule`
  MODIFY `yr_id` int(11) NOT NULL AUTO_INCREMENT;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `apikey`
--
ALTER TABLE `apikey`
  ADD CONSTRAINT `apikey_ibfk_1` FOREIGN KEY (`ownerid`) REFERENCES `user` (`userid`);

--
-- Constraints for table `appointments`
--
ALTER TABLE `appointments`
  ADD CONSTRAINT `appointments_ibfk_1` FOREIGN KEY (`userid`) REFERENCES `user` (`userid`);

--
-- Constraints for table `hashes`
--
ALTER TABLE `hashes`
  ADD CONSTRAINT `hashes_ibfk_1` FOREIGN KEY (`userid`) REFERENCES `user` (`userid`);

--
-- Constraints for table `minion`
--
ALTER TABLE `minion`
  ADD CONSTRAINT `minion_ibfk_1` FOREIGN KEY (`userid`) REFERENCES `user` (`userid`),
  ADD CONSTRAINT `minion_ibfk_2` FOREIGN KEY (`mothershipid`) REFERENCES `motherships` (`msid`),
  ADD CONSTRAINT `minion_ibfk_3` FOREIGN KEY (`ownerid`) REFERENCES `user` (`userid`);

--
-- Constraints for table `motherships`
--
ALTER TABLE `motherships`
  ADD CONSTRAINT `motherships_ibfk_1` FOREIGN KEY (`ownerid`) REFERENCES `user` (`userid`);

--
-- Constraints for table `outputs`
--
ALTER TABLE `outputs`
  ADD CONSTRAINT `outputs_ibfk_1` FOREIGN KEY (`service_id`) REFERENCES `services` (`service_id`);

--
-- Constraints for table `plugins`
--
ALTER TABLE `plugins`
  ADD CONSTRAINT `plugins_ibfk_1` FOREIGN KEY (`owner`) REFERENCES `user` (`userid`);

--
-- Constraints for table `virus`
--
ALTER TABLE `virus`
  ADD CONSTRAINT `virus_ibfk_1` FOREIGN KEY (`aptid`) REFERENCES `apt` (`aptid`);

--
-- Constraints for table `yara_rule`
--
ALTER TABLE `yara_rule`
  ADD CONSTRAINT `yara_rule_ibfk_1` FOREIGN KEY (`ioc_id`) REFERENCES `ioc` (`ioc_id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
