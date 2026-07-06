-- phpMyAdmin SQL Dump
-- version 5.2.0
-- https://www.phpmyadmin.net/
--
-- Host: localhost
-- Generation Time: Feb 22, 2026 at 09:24 PM
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
  `owner_id` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `blogs`
--

CREATE TABLE `blogs` (
  `id` int(11) NOT NULL,
  `uuid` varchar(255) NOT NULL,
  `ownerid` varchar(255) NOT NULL,
  `title` varchar(255) NOT NULL,
  `author` varchar(255) NOT NULL,
  `content` longtext NOT NULL,
  `maintag` varchar(255) DEFAULT NULL,
  `types` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (json_valid(`types`)),
  `categories` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (json_valid(`categories`)),
  `subcats` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (json_valid(`subcats`)),
  `tags` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (json_valid(`tags`)),
  `public` tinyint(1) DEFAULT 0,
  `archived` tinyint(1) DEFAULT 0,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `categories`
--

CREATE TABLE `categories` (
  `id` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `parent_id` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `comments`
--

CREATE TABLE `comments` (
  `id` int(11) NOT NULL,
  `comment_uuid` varchar(255) NOT NULL,
  `blog_uuid` varchar(255) NOT NULL,
  `comment` text NOT NULL,
  `commentor` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `enum_types`
--

CREATE TABLE `enum_types` (
  `id` int(11) NOT NULL,
  `type` varchar(100) NOT NULL,
  `description` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

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
-- Table structure for table `exploits`
--

CREATE TABLE `exploits` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `exploit_id` varchar(255) NOT NULL,
  `lhost` varchar(255) NOT NULL,
  `lport` int(11) NOT NULL,
  `address` varchar(255) NOT NULL,
  `target` varchar(255) NOT NULL,
  `average_severity` int(11) NOT NULL,
  `grouped` tinyint(1) DEFAULT 0,
  `works` tinyint(1) DEFAULT 0,
  `grouped_vulns` text DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
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
  `owner_id` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `scans`
--

INSERT INTO `scans` (`scan_id`, `name`, `scan_type`, `owner_id`, `created_at`, `updated_at`) VALUES
('044c62d9abbca9bca9a17a865dc01589', 'example', 'Bug Bounty', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '2025-03-22 02:38:30', '2025-03-22 02:38:30'),
('12cbf5f4922792f5eb72d2972abe8483', 'example', 'Bug Bounty', '123456', '2025-02-11 10:01:33', '2025-02-11 10:01:33'),
('58819d63e178f1e8cb46e12f23a15795', 'example', 'Bug Bounty', '123456', '2025-02-11 10:00:26', '2025-02-11 10:00:26'),
('67a4a7f7f807673e0a114494ab60e309', 'example', 'Bug Bounty', '123456', '2025-03-20 08:52:27', '2025-03-20 08:52:27'),
('b9c7690613a10b0f5e49bc3b13673c17', 'example', 'Bug Bounty', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '2025-03-22 02:44:31', '2025-03-22 02:44:31'),
('f9328b7c81720ca278a697a7b8d81f75', 'example', 'Bug Bounty', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '2025-03-22 02:37:44', '2025-03-22 02:37:44');

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
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL,
  `data` text DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `services`
--

INSERT INTO `services` (`service_id`, `target_id`, `service_name`, `port`, `protocol`, `state`, `version`, `created_at`, `updated_at`, `data`) VALUES
(2, '60077e5848038468386282765ccd41b1', 'SSH', 22, 'TCP', 1, '2.2.2.2', '2025-02-23 10:38:55', '2025-02-23 10:38:55', '');

-- --------------------------------------------------------

--
-- Table structure for table `targets`
--

CREATE TABLE `targets` (
  `target_id` varchar(100) NOT NULL,
  `scan_id` varchar(100) NOT NULL,
  `host` varchar(255) NOT NULL,
  `host_ip` varchar(45) NOT NULL,
  `target_ip` varchar(45) NOT NULL,
  `firewall_name` varchar(25) NOT NULL DEFAULT 'NONE',
  `decoys` text NOT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `targets`
--

INSERT INTO `targets` (`target_id`, `scan_id`, `host`, `host_ip`, `target_ip`, `firewall_name`, `decoys`, `created_at`, `updated_at`) VALUES
('08d87209ae32dcfb96a84c8ac84056fc', '12cbf5f4922792f5eb72d2972abe8483', 'target2.com', '16843009', '16843009', '', '8.8.8.8,2.2.4.4', '2025-02-21 10:35:47', '2025-02-21 10:35:47'),
('2ff226e95a842f1d5f1c7c329b5dd00c', '12cbf5f4922792f5eb72d2972abe8483', 'testtarget.com', '16843009', '16843009', 'CloudFlare', '8.8.8.8,2.2.4.4', '2025-02-21 10:43:04', '2025-02-21 10:43:04'),
('4676429abe558d0669a46ba7cf470bb4', '12cbf5f4922792f5eb72d2972abe8483', 'domain.com', '16843009', '16843009', '', '8.8.8.8,2.2.4.4', '2025-02-13 09:40:03', '2025-02-13 09:40:03'),
('5e76b9cef4e34c0045fd9c0d98cdaf09', '12cbf5f4922792f5eb72d2972abe8483', 'testtarget.com', '16843009', '16843009', '', '8.8.8.8,2.2.4.4', '2025-02-21 10:41:45', '2025-02-21 10:41:45'),
('60077e5848038468386282765ccd41b1', '12cbf5f4922792f5eb72d2972abe8483', 'workingtarget.com', '16843009', '16843009', 'CloudFlare', '8.8.8.8,2.2.4.4', '2025-02-21 10:47:54', '2025-02-21 10:47:54'),
('dd0b32a96a4ef1321afaebcd084a9596', '12cbf5f4922792f5eb72d2972abe8483', 'testtarget.com', '16843009', '16843009', 'CloudFlare', '8.8.8.8,2.2.4.4', '2025-02-21 10:44:32', '2025-02-21 10:44:32');

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

--
-- Dumping data for table `user`
--

INSERT INTO `user` (`userid`, `ownerid`, `username`, `email`, `password`, `active`, `anonymous`, `verified`, `admin`, `created_at`, `updated_at`) VALUES
('05f65378-685c-45cb-aa5d-f69bc395ba9d', '05f65378-685c-45cb-aa5d-f69bc395ba9d', 'test3', 'test3@mail.com', '$2a$10$mNkYXA5JTpuBxsb9Lfp4m.u4b/LslIn9TppecSwwIlg3lRtcTqKmO', 1, 0, 1, 1, '2025-03-10 07:21:26', '2025-03-10 07:21:26'),
('123456', '12345', 'test user', 'user@mail.com', '$2a$10$GA3fSdt0KF6jMSeTuZdkruNaUhkBEmqTEYZNJu8s.bJ8QhPNAWc6O', 1, 0, 1, 1, '2025-01-20 08:32:53', '2025-01-20 08:32:53'),
('1234567', '12345', 'test user', 'user2@mail.com', '$2a$10$DT0OAtCuMIuXVuNAXhHfsupSz4uZ.2Oman/JZXNv6tqOza3nUfvKK', 1, 0, 1, 1, '2025-01-20 09:17:23', '2025-01-20 09:17:23'),
('27c7dfcf-168c-4744-a572-848489f16634', '27c7dfcf-168c-4744-a572-848489f16634', 'test user3', 'user4@mail.com', '$2a$10$hQRRrNvrQi1q71AB/VVw.eC6PLeRnNY8PauQ0d7ifuQQq1jrxVpKG', 1, 0, 1, 1, '2025-03-14 07:57:30', '2025-03-14 07:57:30'),
('a09003a2-60ef-4123-88f1-58db25136fea', 'a09003a2-60ef-4123-88f1-58db25136fea', 'test', 'test@mail.com', '$2a$10$n.80KxrpBgPs1GoLpa3wh.SvhlG823kKaea5ye202ng8Trdpw81Pm', 1, 0, 1, 1, '2025-03-10 06:27:41', '2025-03-10 06:27:41'),
('cd89975f-525f-4bb8-b3b8-12f4a145535a', 'cd89975f-525f-4bb8-b3b8-12f4a145535a', 'test', 'test2@mail.com', '$2a$10$XjyfXwZfHRhdZXhMyUAtsulH49kNZ2bgCqtobN.s5Um/6obhBPPUi', 1, 0, 1, 1, '2025-03-10 06:40:57', '2025-03-10 06:40:57'),
('ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', 'sam', 'sam@mail.com', '$2a$10$3mKi67y4MbFAog9E1W2Gpe/q2gKuf8MTiiZAhJmWuiATzb6biyECS', 1, 0, 1, 1, '2025-03-14 08:03:58', '2025-03-14 08:03:58');

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
-- Table structure for table `vulnerabilities`
--

CREATE TABLE `vulnerabilities` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `vulnerability_id` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `severity` int(11) NOT NULL,
  `target` varchar(255) NOT NULL,
  `payload` longtext DEFAULT NULL,
  `attack_type` varchar(255) DEFAULT NULL,
  `grouped` tinyint(1) DEFAULT 0,
  `authenticated` tinyint(1) DEFAULT 0,
  `works` tinyint(1) DEFAULT 0,
  `details` longtext DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Table structure for table `webdata`
--

CREATE TABLE `webdata` (
  `target_id` varchar(100) NOT NULL,
  `directory_path` text DEFAULT NULL,
  `parameter_path` text DEFAULT NULL,
  `file_path` text DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `webdata`
--

INSERT INTO `webdata` (`target_id`, `directory_path`, `parameter_path`, `file_path`, `created_at`, `updated_at`) VALUES
('60077e5848038468386282765ccd41b1', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjpbImRvbWFpbi5jby9kb21haW4iLCJkb21haW4uY29tL2Fub3RoZXIvZW5kcG9pbnQiXX0.fMc753BPrpfYAxH4Mki6thImMwX4EgNBEZii1jla9ZI', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjpbImRvbWFpbi5jby9kb21haW4_ZG9tPWhhY2tlZCIsImRvbWFpbi5jb20vYW5vdGhlci9lbmRwb2ludD9uYW1lPWVuZCJdfQ.4fyx8U5iRc8qT-C81imRGgaaNIpdYcF7eA2lR8EPJIY', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjpbImRvbWFpbi5jby9kb21haW4vZG9tYWluLnBkZiIsImRvbWFpbi5jb20vYW5vdGhlci9lbmRwb2ludC9lbmRwb2ludC5wbmciXX0.V0hGoM4K_ExgSEaz6WL3NM9jgprKaII_nnWFb2hxEWA', '2025-02-23 07:10:09', '2025-02-23 07:10:09');

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
  ADD KEY `owner_id` (`owner_id`);

--
-- Indexes for table `blogs`
--
ALTER TABLE `blogs`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `uuid` (`uuid`);

--
-- Indexes for table `categories`
--
ALTER TABLE `categories`
  ADD PRIMARY KEY (`id`),
  ADD KEY `parent_id` (`parent_id`);

--
-- Indexes for table `comments`
--
ALTER TABLE `comments`
  ADD PRIMARY KEY (`id`),
  ADD KEY `blog_uuid` (`blog_uuid`);

--
-- Indexes for table `enum_types`
--
ALTER TABLE `enum_types`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `type` (`type`);

--
-- Indexes for table `events`
--
ALTER TABLE `events`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `exploits`
--
ALTER TABLE `exploits`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `uniq_exploit_id` (`exploit_id`);

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
-- Indexes for table `plugins`
--
ALTER TABLE `plugins`
  ADD PRIMARY KEY (`hash`),
  ADD KEY `owner` (`owner`);

--
-- Indexes for table `scans`
--
ALTER TABLE `scans`
  ADD PRIMARY KEY (`scan_id`),
  ADD UNIQUE KEY `scan_id` (`scan_id`),
  ADD KEY `owner_id` (`owner_id`);

--
-- Indexes for table `services`
--
ALTER TABLE `services`
  ADD PRIMARY KEY (`service_id`),
  ADD KEY `target_id` (`target_id`);

--
-- Indexes for table `targets`
--
ALTER TABLE `targets`
  ADD PRIMARY KEY (`target_id`),
  ADD UNIQUE KEY `target_id` (`target_id`);

--
-- Indexes for table `user`
--
ALTER TABLE `user`
  ADD PRIMARY KEY (`userid`),
  ADD UNIQUE KEY `userid` (`userid`),
  ADD UNIQUE KEY `email` (`email`);

--
-- Indexes for table `virus`
--
ALTER TABLE `virus`
  ADD PRIMARY KEY (`aptid`);

--
-- Indexes for table `vulnerabilities`
--
ALTER TABLE `vulnerabilities`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `uniq_vuln_id` (`vulnerability_id`);

--
-- Indexes for table `webdata`
--
ALTER TABLE `webdata`
  ADD KEY `target_id` (`target_id`);

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
-- AUTO_INCREMENT for table `blogs`
--
ALTER TABLE `blogs`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `categories`
--
ALTER TABLE `categories`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `comments`
--
ALTER TABLE `comments`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `enum_types`
--
ALTER TABLE `enum_types`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `events`
--
ALTER TABLE `events`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `exploits`
--
ALTER TABLE `exploits`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ioc`
--
ALTER TABLE `ioc`
  MODIFY `ioc_id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `services`
--
ALTER TABLE `services`
  MODIFY `service_id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT for table `virus`
--
ALTER TABLE `virus`
  MODIFY `aptid` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `vulnerabilities`
--
ALTER TABLE `vulnerabilities`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

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
-- Constraints for table `assets`
--
ALTER TABLE `assets`
  ADD CONSTRAINT `assets_ibfk_1` FOREIGN KEY (`owner_id`) REFERENCES `user` (`userid`);

--
-- Constraints for table `categories`
--
ALTER TABLE `categories`
  ADD CONSTRAINT `categories_ibfk_1` FOREIGN KEY (`parent_id`) REFERENCES `categories` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `comments`
--
ALTER TABLE `comments`
  ADD CONSTRAINT `comments_ibfk_1` FOREIGN KEY (`blog_uuid`) REFERENCES `blogs` (`uuid`) ON DELETE CASCADE ON UPDATE CASCADE;

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
-- Constraints for table `plugins`
--
ALTER TABLE `plugins`
  ADD CONSTRAINT `plugins_ibfk_1` FOREIGN KEY (`owner`) REFERENCES `user` (`userid`);

--
-- Constraints for table `scans`
--
ALTER TABLE `scans`
  ADD CONSTRAINT `scans_ibfk_1` FOREIGN KEY (`owner_id`) REFERENCES `user` (`userid`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `services`
--
ALTER TABLE `services`
  ADD CONSTRAINT `services_ibfk_1` FOREIGN KEY (`target_id`) REFERENCES `targets` (`target_id`);

--
-- Constraints for table `virus`
--
ALTER TABLE `virus`
  ADD CONSTRAINT `virus_ibfk_1` FOREIGN KEY (`aptid`) REFERENCES `apt` (`aptid`);

--
-- Constraints for table `webdata`
--
ALTER TABLE `webdata`
  ADD CONSTRAINT `webdata_ibfk_1` FOREIGN KEY (`target_id`) REFERENCES `targets` (`target_id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `yara_rule`
--
ALTER TABLE `yara_rule`
  ADD CONSTRAINT `yara_rule_ibfk_1` FOREIGN KEY (`ioc_id`) REFERENCES `ioc` (`ioc_id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
