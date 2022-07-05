SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

-- --------------------------------------------------------

--
-- 表的结构 `admin`
--

CREATE TABLE `admin` (
  `id` int(11) UNSIGNED NOT NULL,
  `username` varchar(60) NOT NULL,
  `password` varchar(120) DEFAULT NULL,
  `name` varchar(60) DEFAULT NULL,
  `role` int(11) UNSIGNED NOT NULL,
  `create_time` int(11) UNSIGNED DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- 转存表中的数据 `admin`
--

INSERT INTO `admin` (`id`, `username`, `password`, `name`, `role`, `create_time`) VALUES
(1, 'admin', '76e7cb24c890edb3cfc8b3f66282d8a94274ef83', 'admin', 1, 1649656761)

-- --------------------------------------------------------

--
-- 表的结构 `admin_auth`
--

CREATE TABLE `admin_auth` (
  `id` int(11) UNSIGNED NOT NULL,
  `name` varchar(60) DEFAULT NULL,
  `method` varchar(20) DEFAULT NULL,
  `path` varchar(200) DEFAULT NULL,
  `enable` tinyint(1) DEFAULT '1',
  `parent_id` int(11) UNSIGNED DEFAULT NULL,
  `create_time` int(11) UNSIGNED DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- 转存表中的数据 `admin_auth`
--

INSERT INTO `admin_auth` (`id`, `name`, `method`, `path`, `enable`, `parent_id`, `create_time`) VALUES
(1, '最高权限', '*', '/admin/**', 1, NULL, 1649656863),
(2, 'api最高权限', '*', '/api/**', 1, NULL, 1649656863);

-- --------------------------------------------------------

--
-- 表的结构 `admin_auth_map`
--

CREATE TABLE `admin_auth_map` (
  `role` int(11) UNSIGNED NOT NULL,
  `auth` int(11) UNSIGNED NOT NULL,
  `create_time` int(11) UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- 转存表中的数据 `admin_auth_map`
--

INSERT INTO `admin_auth_map` (`role`, `auth`, `create_time`) VALUES
(1, 1, 1649656761),
(1, 2, 1650012144);

-- --------------------------------------------------------

--
-- 表的结构 `admin_role`
--

CREATE TABLE `admin_role` (
  `id` int(11) UNSIGNED NOT NULL,
  `name` varchar(60) DEFAULT NULL,
  `describe` varchar(100) DEFAULT NULL,
  `create_time` int(11) UNSIGNED DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- 转存表中的数据 `admin_role`
--

INSERT INTO `admin_role` (`id`, `name`, `describe`, `create_time`) VALUES
(1, '管理员', '最高权限', 1649656761);

-- --------------------------------------------------------

--
-- 表的结构 `taskqueue`
--

CREATE TABLE `taskqueue` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `type` varchar(120) DEFAULT NULL,
  `param` text,
  `fail_msg` text,
  `status` tinyint(4) UNSIGNED DEFAULT '0' COMMENT '0等待  1运行中',
  `retry_times` int(11) UNSIGNED DEFAULT '0',
  `create_time` int(11) UNSIGNED DEFAULT NULL,
  `begin_time` int(11) UNSIGNED DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- 转储表的索引
--

--
-- 表的索引 `admin`
--
ALTER TABLE `admin`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `admin_auth`
--
ALTER TABLE `admin_auth`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `admin_auth_map`
--
ALTER TABLE `admin_auth_map`
  ADD PRIMARY KEY (`role`,`auth`) USING BTREE;

--
-- 表的索引 `admin_role`
--
ALTER TABLE `admin_role`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `taskqueue`
--
ALTER TABLE `taskqueue`
  ADD PRIMARY KEY (`id`),
  ADD KEY `begin_time_union_status` (`begin_time`,`status`) USING BTREE;

--
-- 在导出的表使用AUTO_INCREMENT
--

--
-- 使用表AUTO_INCREMENT `admin`
--
ALTER TABLE `admin`
  MODIFY `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- 使用表AUTO_INCREMENT `admin_auth`
--
ALTER TABLE `admin_auth`
  MODIFY `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- 使用表AUTO_INCREMENT `admin_role`
--
ALTER TABLE `admin_role`
  MODIFY `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- 使用表AUTO_INCREMENT `taskqueue`
--
ALTER TABLE `taskqueue`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
