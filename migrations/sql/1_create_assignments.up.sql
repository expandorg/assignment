CREATE TABLE `assignments` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `job_id` int(10) unsigned NOT NULL,
  `task_id` int(10) unsigned NOT NULL,
  `response_id` int(10) unsigned DEFAULT NULL,
  `worker_id` int(10) unsigned NOT NULL,
  `active` tinyint(1) NOT NULL DEFAULT '1',
  `assigned_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `expires_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `worker_id` (`worker_id`,`job_id`),
  KEY `job_id` (`job_id`),
  KEY `task_id` (`task_id`),
  KEY `response_id` (`response_id`)
)