--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user`
(
    `id`       int          NOT NULL AUTO_INCREMENT,
    `name`     varchar(255) NOT NULL DEFAULT '',
    `password` varchar(255) NOT NULL DEFAULT '',
    PRIMARY KEY (`id`),
    UNIQUE (name)
) ENGINE = InnoDB
  AUTO_INCREMENT = 4
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `vote_item`
--

DROP TABLE IF EXISTS `vote_item`;
/*!40101 SET @saved_cs_client = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vote_item`
(
    `id`          int NOT NULL AUTO_INCREMENT,
    `name`        varchar(200) DEFAULT NULL,
    `description` text,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 4
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci;

--
-- Table structure for table `user_vote`
--

DROP TABLE IF EXISTS `user_vote`;
/*!40101 SET @saved_cs_client = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_vote`
(
    `user_id`      int NOT NULL ,
    `vote_item_id` int NOT NULL,
    KEY `use_fk` (`user_id`),
    KEY `vote_item_fk` (`vote_item_id`),
    CONSTRAINT `use_fk` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT `vote_item_fk` FOREIGN KEY (`vote_item_id`) REFERENCES `vote_item` (`id`) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT user_vote_pk
        PRIMARY KEY (user_id, vote_item_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Dumping data for table `user`
--

INSERT INTO `user`
VALUES (1, 'test', '$2a$10$on2mGRDvXJqn6wk.1ukauusNtpWuAWsQ.1i3zYnFXaAr.AGo2zpte'),
       (2, 'John', '$2a$10$on2mGRDvXJqn6wk.1ukauusNtpWuAWsQ.1i3zYnFXaAr.AGo2zpte'),
       (3, 'test_user', '$2a$10$on2mGRDvXJqn6wk.1ukauusNtpWuAWsQ.1i3zYnFXaAr.AGo2zpte');
