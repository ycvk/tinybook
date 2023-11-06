-- 以下存储过程 用于生成users表的测试数据
DELIMITER //
CREATE PROCEDURE InsertUsers()
BEGIN
    DECLARE i INT DEFAULT 0;
    DECLARE v_email VARCHAR(191);
    DECLARE v_password LONGTEXT;
    DECLARE v_ctime BIGINT;
    DECLARE v_utime BIGINT;
    DECLARE v_nickname LONGTEXT;
    DECLARE v_birthday LONGTEXT;
    DECLARE v_about_me LONGTEXT;
    DECLARE v_phone VARCHAR(20);

    -- 暂时禁用索引和约束 以提高插入速度 但是要注意的是可能会损害数据完整性 如果在插入过程中出现错误 会导致索引和约束不生效
    -- 所以在插入完成后 记得要重新启用索引和约束
    ALTER TABLE `users`
        DISABLE KEYS;

    -- 开启事务
    START TRANSACTION;

    -- 插入100万条数据
    WHILE i < 1000000
        DO
            SET v_email = CONCAT('user', i, '@test.insert.com');
            SET v_password = '$2a$10$iHQWgti3td42eFX5akkP5eWAVNOtcbxnpVdzYwNOBggqE1ZPDOLXG';-- 密码为hello#world123
            SET v_ctime = UNIX_TIMESTAMP();
            SET v_utime = UNIX_TIMESTAMP();
            SET v_nickname = CONCAT('Nick', i);
            SET v_birthday = '1990-01-01';
            SET v_about_me = 'About Me';
            SET v_phone = CONCAT('1234', i);

            INSERT INTO `users` (`email`, `password`, `ctime`, `utime`, `nickname`, `birthday`, `about_me`, `phone`)
            VALUES (v_email, v_password, v_ctime, v_utime, v_nickname, v_birthday, v_about_me, v_phone);

            SET i = i + 1;

            -- 每插入10000条数据 提交一次事务 以防止事务过大 可以根据自己的机器情况调整
            IF i MOD 10000 = 0 THEN
                COMMIT;
                START TRANSACTION;
            END IF;
        END WHILE;

    -- 提交剩余的数据
    COMMIT;

-- 重新启用索引和约束
    ALTER TABLE `users`
        ENABLE KEYS;
END;
//
DELIMITER ;

-- 执行存储过程
CALL InsertUsers();
