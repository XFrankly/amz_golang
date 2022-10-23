package main

/*
比较时间 sql
SELECT * FROM product
WHERE date(add_time) BETWEEN '2019-11-05' AND '2019-11-30'



時間直接對比 sql
		select * from test where create_time between ‘2015-03-03 17:39:05’ and ‘2016-03-03 17:39:52’;
		方法一：直接比较
		select * from test where create_time between ‘2015-03-03 17:39:05’ and ‘2016-03-03 17:39:52’;

		方法二：用unix_timestamp函数，将字符型的时间，转成unix时间戳
		select * from test where unix_timestamp(create_time) > unix_timestamp(‘2011-03-03 17:39:05’) and unix_timestamp(create_time) < unix_timestamp(‘2011-03-03 17:39:52’);
		个人觉得这样比较更踏实点儿。

		Oracle（Date，TimeStamp等）：
		方法一：将字符串转换为日期类型
		select * from test where create_time between to_date(‘2015-03-03 17:39:05’) and to_date(‘2016-03-03 17:39:52’);

		二.存储日期类型的字段为数值类型
		MySql（bigint）：
		方法一：将日期字符串转换为时间戳
		select * from test where create_time > unix_timestamp(‘2011-03-03 17:39:05’) and create_time< unix_timestamp(‘2011-03-03 17:39:52’);

		方法二：将时间戳转换为日期类型
		select * from test where from_unixtime(create_time/1000) between ‘2014-03-03 17:39:05’ and ‘2015-03-03 17:39:52’);
*/
