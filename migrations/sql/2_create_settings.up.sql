CREATE TABLE `settings` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `limit` int(11) DEFAULT NULL,
  `repeat` tinyint(1) NOT NULL DEFAULT '1',
  `job_id` int(11) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `job_id` (`job_id`)
)