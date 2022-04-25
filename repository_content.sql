SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";

CREATE TABLE `REPOSITORY_CONTENT` (
  `PAGE` int(11) UNSIGNED NOT NULL COMMENT 'Page of GitHub API response',
  `CONTENT` mediumblob NOT NULL COMMENT 'GZip encoded GitHub repository API response',
  `UPDATED_AT` tinytext COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT current_timestamp() COMMENT 'Timestamp when this content updated'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


ALTER TABLE `REPOSITORY_CONTENT`
  ADD PRIMARY KEY (`PAGE`);
COMMIT;

