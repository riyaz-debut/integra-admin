-- MySQL dump 10.13  Distrib 8.0.29, for Linux (x86_64)
--
-- Host: localhost    Database: integradb
-- ------------------------------------------------------
-- Server version	8.0.29

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `chaincode_lists`
--

DROP TABLE IF EXISTS `chaincode_lists`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chaincode_lists` (
  `id` int NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `label` varchar(255) DEFAULT NULL,
  `version` varchar(255) DEFAULT NULL,
  `sequence` int DEFAULT NULL,
  `status` int DEFAULT NULL,
  `url` varchar(255) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `chaincode_lists`
--

LOCK TABLES `chaincode_lists` WRITE;
/*!40000 ALTER TABLE `chaincode_lists` DISABLE KEYS */;
INSERT INTO `chaincode_lists` VALUES (1,'example_cc','example_cc','1',1,1,'https://github.com/harjot-debut/example_cc/raw/main/example_cc_2.zip','2022-06-24 11:30:44','2022-08-03 09:47:55'),(2,'example_cc_2','example_cc_2','1',1,1,'https://github.com/harjot-debut/example_cc/raw/main/example_cc.zip','2022-06-23 05:17:01','2022-08-03 09:46:17');
/*!40000 ALTER TABLE `chaincode_lists` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `chaincode_logs`
--

DROP TABLE IF EXISTS `chaincode_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chaincode_logs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `cu_id` int DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `label` varchar(255) DEFAULT NULL,
  `version` varchar(255) DEFAULT NULL,
  `sequence` int DEFAULT NULL,
  `org_id` int DEFAULT NULL,
  `org_name` varchar(255) DEFAULT NULL,
  `org_msp` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `cc_logs_fk` (`cu_id`),
  CONSTRAINT `cc_logs_fk` FOREIGN KEY (`cu_id`) REFERENCES `chaincode_updates` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `chaincode_logs`
--

LOCK TABLES `chaincode_logs` WRITE;
/*!40000 ALTER TABLE `chaincode_logs` DISABLE KEYS */;
INSERT INTO `chaincode_logs` VALUES (1,10,'example_cc','example_cc','1',1,2,'Org1','Org1MSP','2022-06-30 09:15:04','2022-06-30 09:15:04'),(9,10,'example_cc','example_cc','1',1,2,'Org2','Org2MSP','2022-06-24 17:00:44','2022-07-29 16:43:52');
/*!40000 ALTER TABLE `chaincode_logs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `chaincode_updates`
--

DROP TABLE IF EXISTS `chaincode_updates`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chaincode_updates` (
  `id` int NOT NULL AUTO_INCREMENT,
  `cc_id` int DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `label` varchar(255) DEFAULT NULL,
  `version` varchar(255) DEFAULT NULL,
  `sequence` int DEFAULT NULL,
  `status` int DEFAULT NULL,
  `url` varchar(255) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `cu_fk` (`cc_id`),
  CONSTRAINT `cu_fk` FOREIGN KEY (`cc_id`) REFERENCES `chaincode_lists` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `chaincode_updates`
--

LOCK TABLES `chaincode_updates` WRITE;
/*!40000 ALTER TABLE `chaincode_updates` DISABLE KEYS */;
INSERT INTO `chaincode_updates` VALUES (10,1,'example_cc','example_cc','1',2,1,'https://github.com/harjot-debut/example_cc/raw/main/example_cc_3.zip','2022-08-03 07:01:30','2022-08-03 07:01:30'),(11,2,'example_cc_2','example_cc_2','1',2,1,'https://github.com/harjot-debut/example_cc/raw/main/example_cc_4.zip','2022-08-03 23:35:03','2022-08-03 23:35:03'),(14,1,'example_cc','example_cc','1',3,1,'https://github.com/harjot-debut/example_cc/raw/main/example_cc_3.zip','2022-06-24 11:30:44','2022-08-16 10:25:23');
/*!40000 ALTER TABLE `chaincode_updates` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `org_signatures`
--

DROP TABLE IF EXISTS `org_signatures`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `org_signatures` (
  `id` int NOT NULL AUTO_INCREMENT,
  `org_id` int DEFAULT NULL,
  `org_name` varchar(255) DEFAULT NULL,
  `org_msp` varchar(255) DEFAULT NULL,
  `signby_id` int DEFAULT NULL,
  `signature` text,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `org_signatures`
--

LOCK TABLES `org_signatures` WRITE;
/*!40000 ALTER TABLE `org_signatures` DISABLE KEYS */;
/*!40000 ALTER TABLE `org_signatures` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `organizations`
--

DROP TABLE IF EXISTS `organizations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `organizations` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `msp_id` varchar(255) DEFAULT NULL,
  `peers_count` int DEFAULT NULL,
  `config` text,
  `modified_config` text,
  `status` int DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `organizations`
--

LOCK TABLES `organizations` WRITE;
/*!40000 ALTER TABLE `organizations` DISABLE KEYS */;
INSERT INTO `organizations` VALUES (1,'AdminOrg','adminOrgMSP',2,NULL,NULL,NULL,'2022-06-30 10:28:35','2022-06-30 10:28:35'),(2,'Org1','Org1MSP',2,NULL,NULL,NULL,'2022-06-30 10:28:35','2022-06-30 10:28:35'),(3,'Org2','Org2MSP',2,NULL,NULL,NULL,'2022-06-30 10:28:35','2022-06-30 10:28:35');
/*!40000 ALTER TABLE `organizations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `peers`
--

DROP TABLE IF EXISTS `peers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `peers` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `org_id` int DEFAULT NULL,
  `org_name` varchar(255) DEFAULT NULL,
  `url` varchar(255) DEFAULT NULL,
  `ip` varchar(255) DEFAULT NULL,
  `cert` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `peers`
--

LOCK TABLES `peers` WRITE;
/*!40000 ALTER TABLE `peers` DISABLE KEYS */;
INSERT INTO `peers` VALUES (1,'peer0.org3.example.com',4,'Org3','grpcs://localhost:10051','160.788.78.00','peeroOrg3cert','2022-06-30 11:20:12','2022-06-30 11:20:12'),(2,'peer1.org3.example.com',4,'Org3','grpcs://localhost:11051','168.188.58.60','peeroOrg3cert','2022-06-30 11:20:12','2022-06-30 11:20:12');
/*!40000 ALTER TABLE `peers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `user_name` varchar(255) DEFAULT NULL,
  `password` varchar(255) DEFAULT NULL,
  `org_id` int DEFAULT NULL,
  `org_name` varchar(255) DEFAULT NULL,
  `org_msp` varchar(255) DEFAULT NULL,
  `role` varchar(255) DEFAULT NULL,
  `status` int DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`),
  UNIQUE KEY `user_name` (`user_name`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'adminuser','$2a$10$Jcu7tCu1XXjGbt7Sm3rBgOTZoIbwdSahIk3GWCAtOM7Ab4yX/DB2m',1,'adminOrg','adminMSP','admin',1,'2022-07-01 04:46:53','2022-07-01 04:46:53'),(2,'user1@org1','$2a$10$5JNybMffRvuLMHk/jypcuu6M6vIrO2FR8Nn5Kropd3Toocv9.8cD2',2,'Org1','Org1MSP','user',1,'2022-07-01 04:47:56','2022-07-01 04:47:56'),(3,'user2@org2','$2a$10$XHUEN5.YaHpGFIm7JLmYqu33hHxUoNvx3JAi4EDBX1yS65VW8dsWC',3,'Org2','Org2MSP','user',1,'2022-07-01 04:49:46','2022-07-01 04:49:46'),(6,'testuser','$2a$10$zKVqsqNhN9oJRVzLgEg9sOLglo2HLom4LxGCtTOSxd.9VxaHzMOTq',3,'Org2','Org2MSP','user',0,'2022-08-16 14:05:38','2022-08-16 14:05:38'),(7,'teset','$2a$10$2hXW/e0gPgR15McqcMm9WOnCE1k8pPZauE43Dydl6FRbqPvuF5.XW',3,'3','Org2MSP','user',0,'2022-08-16 14:13:12','2022-08-16 14:13:12'),(8,'testser','$2a$10$LVq0vWT7D87WJB.HdCVouuIxarXJoUx7iXatHpY.RC6Xj15pA63Ae',3,'3','Org2MSP','user',0,'2022-08-16 14:18:18','2022-08-16 14:18:18'),(9,'tester','$2a$10$fGOxlA9A9Mt967PYbcwC..Jesut8Md5wurH1XsAMm1OmofeZ.gy9G',3,'Org2','Org2MSP','user',0,'2022-08-16 14:24:03','2022-08-16 14:24:03'),(12,'testuserer','$2a$10$UDgKHDEfutAbKOd7pGeqVORgfPwAXc6ujxszCQ.LOFp8CUdhsSpoq',3,'Org2','Org2MSP','user',0,'2022-08-16 14:25:30','2022-08-16 14:25:30');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-08-17 16:31:39
