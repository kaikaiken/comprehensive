# -*- coding: utf-8 -*-

# Define here the models for your scraped items
#
# See documentation in:
# https://docs.scrapy.org/en/latest/topics/items.html

import scrapy


class BangumispiderItem(scrapy.Item):
    # define the fields for your item here like:
    # name = scrapy.Field()
    pass


class BangumiItem(scrapy.Item):
    bangumi_id = scrapy.Field()  # 番剧id
    name = scrapy.Field()  # 番剧名
    cover_url = scrapy.Field()  # 海报url
    bangumi_score = scrapy.Field()  # Bangumi网站的分数
    vote_num = scrapy.Field()  # 评分人数
    episode_num = scrapy.Field()
    tags = scrapy.Field()  # 前几个标签
    desc = scrapy.Field()  # 番剧描述
    staff_list = scrapy.Field()  # 制作者（以json表示）
    cv_list = scrapy.Field()  # 声优列表
