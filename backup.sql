-- MySQL dump 10.13  Distrib 8.0.31, for Win64 (x86_64)
--
-- Host: localhost    Database: ecommerce
-- ------------------------------------------------------
-- Server version	8.0.31

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
-- Table structure for table `carts`
--

DROP TABLE IF EXISTS `carts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `carts` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned DEFAULT NULL,
  `product_id` bigint unsigned DEFAULT NULL,
  `quantity` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_products_cart` (`product_id`),
  CONSTRAINT `fk_products_cart` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `carts`
--

LOCK TABLES `carts` WRITE;
/*!40000 ALTER TABLE `carts` DISABLE KEYS */;
INSERT INTO `carts` VALUES (26,1,1,4);
/*!40000 ALTER TABLE `carts` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `chats`
--

DROP TABLE IF EXISTS `chats`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chats` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL,
  `product_id` bigint unsigned NOT NULL,
  `admin_id` bigint unsigned DEFAULT '0',
  `message` text NOT NULL,
  `is_admin` tinyint(1) DEFAULT '0',
  `read` tinyint(1) DEFAULT '0',
  `attachment` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_products_chats` (`product_id`),
  KEY `fk_users_chats` (`user_id`),
  CONSTRAINT `fk_products_chats` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`),
  CONSTRAINT `fk_users_chats` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=60 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `chats`
--

LOCK TABLES `chats` WRITE;
/*!40000 ALTER TABLE `chats` DISABLE KEYS */;
INSERT INTO `chats` VALUES (3,'2025-01-26 10:06:05.694','2025-01-26 10:06:05.694',1,1,0,'hello',0,0,''),(4,'2025-01-26 11:33:44.683','2025-01-26 11:33:44.683',1,1,0,'h;iuh;iu',0,0,''),(7,'2025-01-26 11:15:45.745','2025-01-26 11:15:45.745',1,1,0,'hy',0,0,''),(8,'2025-01-26 11:16:10.697','2025-01-26 11:16:10.697',1,1,0,'poo',0,0,''),(9,'2025-01-26 11:31:16.363','2025-01-26 11:31:16.363',1,1,3,'pop',1,0,''),(10,'2025-01-26 11:32:56.190','2025-01-26 11:32:56.190',1,1,3,'pop2',1,0,''),(11,'2025-01-26 11:34:48.025','2025-01-26 11:34:48.025',1,1,3,'pop5',1,0,''),(12,'2025-01-26 11:35:52.764','2025-01-26 11:35:52.764',1,1,0,'yghyh',0,0,''),(13,'2025-01-26 11:39:04.879','2025-01-26 11:39:04.879',1,1,0,'kguy',0,0,''),(14,'2025-01-26 11:40:27.351','2025-01-26 11:40:27.351',1,1,0,'kgyugyugoyu',0,0,''),(15,'2025-01-26 11:40:40.652','2025-01-26 11:40:40.652',1,1,0,'foyuiytfiyt',0,0,''),(16,'2025-01-26 11:40:59.063','2025-01-26 11:40:59.063',1,1,0,'bhblyugyiyt',0,0,''),(17,'2025-01-26 11:41:26.298','2025-01-26 11:41:26.298',1,1,3,'huluigy',1,0,''),(18,'2025-01-26 11:42:30.691','2025-01-26 11:42:30.691',1,1,0,'ugygulygly',0,0,''),(19,'2025-01-28 04:20:02.175','2025-01-28 04:20:02.175',1,1,3,'from admin',1,0,''),(20,'2025-01-28 04:22:44.946','2025-01-28 04:22:44.946',1,1,3,'admin again',1,0,''),(21,'2025-01-28 04:33:46.708','2025-01-28 04:33:46.708',1,1,3,'hi',1,0,''),(22,'2025-01-28 04:37:44.618','2025-01-28 04:37:44.618',1,1,3,'pop',1,0,''),(23,'2025-01-28 04:45:07.740','2025-01-28 04:45:07.740',1,1,3,'ol',1,0,''),(24,'2025-01-28 05:24:46.166','2025-01-28 05:24:46.166',1,1,3,'suo',1,0,''),(25,'2025-01-28 05:25:08.248','2025-01-28 05:25:08.248',1,1,0,'nics',0,0,''),(26,'2025-01-28 05:50:49.554','2025-01-28 05:50:49.554',1,1,3,'pop',1,0,''),(27,'2025-01-28 05:50:59.897','2025-01-28 05:50:59.897',1,1,0,'nice',0,0,''),(28,'2025-01-28 05:58:03.320','2025-01-28 05:58:03.320',1,1,0,'helli',0,0,''),(29,'2025-01-28 05:59:44.532','2025-01-28 05:59:44.532',1,1,0,'pop',0,0,''),(30,'2025-01-28 06:02:31.323','2025-01-28 06:02:31.323',1,1,0,'lll',0,0,''),(31,'2025-01-28 06:14:18.352','2025-01-28 06:14:18.352',1,1,0,'pop',0,0,''),(32,'2025-01-28 06:18:47.304','2025-01-28 06:18:47.304',1,1,0,'pop',0,0,''),(33,'2025-01-28 06:19:40.848','2025-01-28 06:19:40.848',1,1,0,'you did',0,0,''),(34,'2025-01-28 06:20:50.468','2025-01-28 06:20:50.468',1,1,0,'what',0,0,''),(35,'2025-01-28 06:22:06.837','2025-01-28 06:22:06.837',1,1,0,'what',0,0,''),(36,'2025-01-28 06:24:39.719','2025-01-28 06:24:39.719',1,1,0,'what',0,0,''),(37,'2025-01-28 06:24:53.342','2025-01-28 06:24:53.342',1,1,0,'fuck',0,0,''),(38,'2025-01-28 06:25:02.064','2025-01-28 06:25:02.064',1,1,3,'hey',1,0,''),(39,'2025-01-28 06:25:43.698','2025-01-28 06:25:43.698',1,1,0,'sup',0,0,''),(40,'2025-01-28 06:28:01.447','2025-01-28 06:28:01.447',1,1,0,'what',0,0,''),(41,'2025-01-28 06:29:32.375','2025-01-28 06:29:32.375',1,1,0,'pop',0,0,''),(42,'2025-01-28 06:33:19.701','2025-01-28 06:33:19.701',1,1,0,'nice',0,0,''),(43,'2025-01-28 07:10:46.516','2025-01-28 07:10:46.516',1,1,0,'hey',0,0,''),(44,'2025-01-28 07:11:31.400','2025-01-28 07:11:31.400',1,1,0,'pop',0,0,''),(45,'2025-01-28 07:13:47.068','2025-01-28 07:13:47.068',1,1,0,'hey',0,0,''),(46,'2025-01-28 07:15:12.919','2025-01-28 07:15:12.919',1,1,0,'what',0,0,''),(47,'2025-01-28 06:48:19.416','2025-01-28 06:48:19.416',1,1,0,'hello',0,0,''),(48,'2025-01-28 06:48:39.336','2025-01-28 06:48:39.336',1,1,3,'how can i help you',1,0,''),(49,'2025-01-28 06:49:24.308','2025-01-28 06:49:24.308',1,1,0,'i would like to have this product but the price is too high',0,0,''),(50,'2025-01-28 06:57:29.379','2025-01-28 06:57:29.379',1,1,0,'hey',0,0,''),(51,'2025-01-28 06:57:57.564','2025-01-28 06:57:57.564',1,1,3,'i dont even under stand this UI',1,0,''),(52,'2025-01-28 06:58:18.507','2025-01-28 06:58:18.507',1,1,0,'why are you saying this',0,0,''),(53,'2025-01-28 06:58:39.186','2025-01-28 06:58:39.186',1,1,3,'look at the arrange ment ',1,0,''),(54,'2025-01-28 07:12:04.112','2025-01-28 07:12:04.112',1,2,0,'how much',0,0,''),(55,'2025-01-28 07:12:30.711','2025-01-28 07:12:30.711',1,2,3,'you saw the price',1,0,''),(56,'2025-02-11 07:39:35.728','2025-02-11 07:39:35.728',1,2,0,'pop',0,0,''),(57,'2025-02-11 07:40:14.083','2025-02-11 07:40:14.083',1,2,0,'how much is this',0,0,''),(58,'2025-02-11 07:41:24.399','2025-02-11 07:41:24.399',1,2,3,'20000',1,0,''),(59,'2025-02-23 04:16:31.434','2025-02-23 04:16:31.434',1,1,0,'nice',0,0,'');
/*!40000 ALTER TABLE `chats` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `orders`
--

DROP TABLE IF EXISTS `orders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `orders` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `product_id` bigint unsigned NOT NULL,
  `quantity` bigint NOT NULL,
  `total_price` double NOT NULL,
  `status` varchar(20) DEFAULT 'pending',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `seller` tinyint(1) DEFAULT '0',
  `buyer` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `fk_orders_user` (`user_id`),
  KEY `fk_orders_product` (`product_id`),
  CONSTRAINT `fk_orders_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`),
  CONSTRAINT `fk_orders_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `orders`
--

LOCK TABLES `orders` WRITE;
/*!40000 ALTER TABLE `orders` DISABLE KEYS */;
/*!40000 ALTER TABLE `orders` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pictures`
--

DROP TABLE IF EXISTS `pictures`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `pictures` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `url` varchar(255) NOT NULL,
  `product_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_products_pictures` (`product_id`),
  CONSTRAINT `fk_products_pictures` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pictures`
--

LOCK TABLES `pictures` WRITE;
/*!40000 ALTER TABLE `pictures` DISABLE KEYS */;
INSERT INTO `pictures` VALUES (1,'2025-01-18 18:35:27.502','2025-01-18 18:35:27.502','/staticproductimage/ce098174-a58d-4c4a-b585-109317345364-1838960677.png',1),(2,'2025-01-18 18:35:27.502','2025-01-18 18:35:27.502','/staticproductimage/2eaedd85-2bb6-4ab4-9287-516955c133e8-344991014.png',1),(3,'2025-01-28 07:09:51.048','2025-01-28 07:09:51.048','/staticproductimage/e433be47-fa71-4c44-acc2-c44129be15e7-3396943133.png',2),(4,'2025-01-28 07:09:51.048','2025-01-28 07:09:51.048','/staticproductimage/a9041450-b72f-4769-be67-38d2c7df8a56-3823468020.png',2),(5,'2025-02-16 00:40:26.853','2025-02-16 00:40:26.853','/staticproductimage/6304645d-aa0b-4d21-8129-f36df6a0a435-3427221210.png',3),(6,'2025-02-16 00:40:26.853','2025-02-16 00:40:26.853','/staticproductimage/7eea9014-91c6-4edb-84bf-f520f53245c2-332551421.png',3),(7,'2025-02-16 00:40:26.853','2025-02-16 00:40:26.853','/staticproductimage/2a16a149-1c03-4ee9-bdfd-69ae30cb1d93-92958450.png',3),(8,'2025-02-16 00:43:22.948','2025-02-16 00:43:22.948','/staticproductimage/6d0d01b9-712f-45cf-931f-bc50c33809b2-4141222742.png',4);
/*!40000 ALTER TABLE `pictures` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `products`
--

DROP TABLE IF EXISTS `products`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `products` (
  `user_id` bigint unsigned DEFAULT NULL,
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `name` varchar(100) NOT NULL,
  `description` text NOT NULL,
  `category` longtext,
  `price` decimal(10,2) NOT NULL,
  `stock` bigint NOT NULL,
  `sku` varchar(100) DEFAULT NULL,
  `size` varchar(50) NOT NULL,
  `color` varchar(50) NOT NULL,
  `brand` varchar(100) NOT NULL,
  `material` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_products_sku` (`sku`),
  KEY `fk_users_product` (`user_id`),
  FULLTEXT KEY `name` (`name`,`description`,`brand`),
  CONSTRAINT `fk_users_product` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `products`
--

LOCK TABLES `products` WRITE;
/*!40000 ALTER TABLE `products` DISABLE KEYS */;
INSERT INTO `products` VALUES (1,1,'2025-01-18 18:35:27.310','2025-01-18 18:35:27.310','Shirt','Nicely fitted suites of you desire','Men\'s Wear',5500.00,4,'SHIRT-1737225327308496000','M','Blue','Nike','Cotton'),(3,2,'2025-01-28 07:09:50.755','2025-01-28 07:09:50.755','Trousers','your joggers, fitted and nice','Men\'s Wear',7000.00,30,'TROUSERS-1738048190754311000','M','Gray','Chinos','Cotton'),(1,3,'2025-02-16 00:40:26.515','2025-02-16 00:40:26.515','Bag','Nice and durable female bag with high quality','Women\'s Wear',15000.00,15,'BAG-1739666426513519700','M','Black','lether brand','Polyester'),(1,4,'2025-02-16 00:43:22.849','2025-02-16 00:43:22.849','Shoe','Simple footwear for any occation','Women\'s Wear',7000.00,22,'SHOE-1739666602848033200','M','Red','Quark','Polyester');
/*!40000 ALTER TABLE `products` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `profiles`
--

DROP TABLE IF EXISTS `profiles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `profiles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `first_name` longtext,
  `last_name` longtext,
  `image` longtext,
  `user_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_profiles_user_id` (`user_id`),
  CONSTRAINT `fk_users_profile` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `profiles`
--

LOCK TABLES `profiles` WRITE;
/*!40000 ALTER TABLE `profiles` DISABLE KEYS */;
INSERT INTO `profiles` VALUES (1,'Ahibi','emmanuel','',1),(3,'Ahibi','emmanuel','',3);
/*!40000 ALTER TABLE `profiles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `reviews`
--

DROP TABLE IF EXISTS `reviews`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `reviews` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `rating` bigint NOT NULL,
  `comment` text,
  `product_id` bigint unsigned NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_users_reviews` (`user_id`),
  KEY `fk_products_reviews` (`product_id`),
  CONSTRAINT `fk_products_reviews` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`),
  CONSTRAINT `fk_users_reviews` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `reviews`
--

LOCK TABLES `reviews` WRITE;
/*!40000 ALTER TABLE `reviews` DISABLE KEYS */;
INSERT INTO `reviews` VALUES (1,'2025-01-18 18:36:08.774','2025-01-18 18:36:08.774',0,'i love this shirt',1,1),(2,'2025-01-18 18:36:29.588','2025-01-18 18:36:29.588',0,'pls how much is it pls',1,1),(3,'2025-01-18 18:36:53.758','2025-01-18 18:36:53.758',0,'any body ',1,1),(4,'2025-01-18 18:37:08.541','2025-01-18 18:37:08.541',0,'hello',1,1),(5,'2025-01-18 18:37:19.938','2025-01-18 18:37:19.938',0,'what\'s sup',1,1),(6,'2025-01-18 18:37:57.131','2025-01-18 18:37:57.131',0,'well dnd',1,1),(7,'2025-01-18 18:38:00.822','2025-01-18 18:38:00.822',0,'dkn;iuere',1,1),(8,'2025-01-18 18:38:07.787','2025-01-18 18:38:07.787',0,'whpuhper\\',1,1),(9,'2025-01-18 18:38:11.124','2025-01-18 18:38:11.124',0,'ihwe[98u[8ew',1,1),(10,'2025-01-18 18:38:14.439','2025-01-18 18:38:14.439',0,'lkjwe\'oij\'ioe',1,1),(11,'2025-01-18 18:38:18.023','2025-01-18 18:38:18.023',0,';lkwe]09]0re',1,1),(12,'2025-01-18 18:38:49.895','2025-01-18 18:38:49.895',0,'popup',1,1),(13,'2025-01-18 18:38:59.427','2025-01-18 18:38:59.427',0,'😂😂😂',1,1),(14,'2025-01-28 07:10:43.142','2025-01-28 07:10:43.142',0,'i love this',2,1),(15,'2025-01-28 07:02:29.141','2025-01-28 07:02:29.141',0,'hello',2,3),(16,'2025-01-28 07:06:31.416','2025-01-28 07:06:31.416',0,'what up',2,3),(17,'2025-01-28 06:50:19.858','2025-01-28 06:50:19.858',0,'😂😂😂',2,3),(18,'2025-01-28 06:55:02.187','2025-01-28 06:55:02.187',0,'testing',2,1),(19,'2025-01-28 06:55:15.916','2025-01-28 06:55:15.916',0,'ok',2,3),(20,'2025-02-23 04:01:23.032','2025-02-23 04:01:23.032',0,'here again',1,1);
/*!40000 ALTER TABLE `reviews` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `username` varchar(191) DEFAULT NULL,
  `password` longtext,
  `email` varchar(191) DEFAULT NULL,
  `verification_code` longtext,
  `email_verified` tinyint(1) DEFAULT '0',
  `role` longtext,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_users_username` (`username`),
  UNIQUE KEY `uni_users_email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'2025-01-18 18:33:30.364','2025-01-18 18:34:20.160','emmanuelahibi3','$2a$10$RILbXByB5Xdd18bxbiO8CeafteXCAg2yfSbCnxzbiej7itFhAwSQe','emmanuelahibi3@gmail.com','used',1,'user'),(3,'2025-01-22 08:00:04.735','2025-01-22 08:00:27.753','phychemma4','$2a$10$n9rIWWMw3caXiSJNla.zpuVfdT2eFKIUr2xVTV2QL0zjpP3hZ4.GW','phychemma4@gmail.com','used',1,'admin');
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

-- Dump completed on 2025-02-27 23:30:34
