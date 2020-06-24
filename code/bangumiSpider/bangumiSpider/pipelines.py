# -*- coding: utf-8 -*-

# Define your item pipelines here
#
# Don't forget to add your pipeline to the ITEM_PIPELINES setting
# See: https://docs.scrapy.org/en/latest/topics/item-pipeline.html
import pymysql
from bangumiSpider import settings


class BangumispiderPipeline(object):
    def __init__(self):
        self.connect = pymysql.connect(
            host=settings.MYSQL_HOST,
            db=settings.MYSQL_DBNAME,
            user=settings.MYSQL_USER,
            passwd=settings.MYSQL_PASSWD,
            charset='utf8',
            use_unicode=True)
        self.cursor = self.connect.cursor()

    def process_item(self, item, spider):
        bangumi_id = pymysql.escape_string(item['bangumi_id'])
        name = pymysql.escape_string(item['name'])
        cover_url = pymysql.escape_string(item['cover_url'])
        bangumi_score = pymysql.escape_string(item['bangumi_score'])
        vote_num = pymysql.escape_string(item['vote_num'])
        episode_num = pymysql.escape_string(item['episode_num'])
        tags = pymysql.escape_string(item['tags'])
        desc = pymysql.escape_string(item['desc'])
        staff_list = pymysql.escape_string(item['staff_list'])
        cv_list = pymysql.escape_string(item['cv_list'])
        sql_command = "insert ignore into bangumi values ( '{0}', '{1}', '{2}', '{3}', '{4}', '{5}', '{6}', '{7}', '{8}', '{9}');".format(bangumi_id, name, cover_url, bangumi_score, vote_num, episode_num, tags, desc, staff_list, cv_list)
        self.cursor.execute(sql_command)
        self.connect.commit()
        print("sql_command = " + sql_command)
        return item
