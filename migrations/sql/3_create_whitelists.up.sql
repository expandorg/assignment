CREATE TABLE `whitelists` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `worker_id` int(11) unsigned NOT NULL,
  `job_id` int(11) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `job_worker` (`job_id`,`worker_id`)
)