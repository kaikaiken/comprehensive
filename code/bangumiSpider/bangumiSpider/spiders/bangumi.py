# -*- coding: utf-8 -*-
import scrapy
from bangumiSpider.items import BangumiItem
import json
import re
import requests


class BangumiSpider(scrapy.Spider):
    name = 'bangumi'
    allowed_domains = ['bangumi.tv']
    start_urls = ['https://bangumi.tv/anime/browser?sort=rank&page=234']

    def parse(self, response):
        bangumi_list = response.xpath('//*[@id="browserItemList"]//li/@id').extract()
        for bangumi in bangumi_list:
            bangumi_id = bangumi.split('_')[-1]
            detail_url = 'https://bangumi.tv/subject/' + str(bangumi_id)
            yield scrapy.Request(url=detail_url, meta={"bangumi_id": bangumi_id}, callback=self.parse_detail)
        pre_url = response.xpath('//*[@id="columnSubjectBrowserA"]/div[2]/div/div/a[2]/@href').extract()
        if pre_url:
            previous_url = 'https://bangumi.tv/anime/browser' + pre_url[0]
            print("######## url = " + str(previous_url))
            # 把新的页面url加入待爬取页面
            yield scrapy.Request(previous_url, callback=self.parse)

    def parse_detail(self, response):
        bangumi_item = BangumiItem()
        info_dict = {}
        info_count = 1
        while True:
            try:
                info_key = response.xpath('//*[@id="infobox"]/li[' + str(info_count) + ']/span/text()').extract()[0][
                           :-2]
            except IndexError:
                break
            else:
                if response.xpath('//*[@id="infobox"]/li[' + str(info_count) + ']/a/text()').extract():
                    info_content = ','.join(
                        response.xpath('//*[@id="infobox"]/li[' + str(info_count) + ']/a/text()').extract())
                else:
                    info_content = response.xpath('//*[@id="infobox"]/li[' + str(info_count) + ']/text()').extract()[0]
                # print("i=" + str(info_count) + ', span=' + str(info_key) + ', content=' + str(info_content))
                info_count = info_count + 1
                info_dict[info_key] = str(info_content)
        # info_json = json.dumps(info_dict, indent=2, ensure_ascii=False)
        # print(info_json)
        bangumi_id = response.meta['bangumi_id']
        try:
            name = info_dict['中文名']
        except KeyError:
            name = response.xpath('//*[@id="headerSubject"]/h1/a/text()').extract()
        try:
            episode_num = info_dict['话数']
        except KeyError:
            episode_num = -1
        try:
            cover_url = 'https:' + str(response.xpath('//*[@id="bangumiInfo"]/div/div[1]/a/img/@src').extract()[0])
        except IndexError:
            cover_url = 'https:' + str(response.xpath('//*[@id="bangumiInfo"]/div/div[1]/a/img/@src').extract())
        try:
            bangumi_score = response.xpath('//*[@class="global_score"]/span[1]/text()').extract()[0]
        except IndexError:
            bangumi_score = response.xpath('//*[@class="global_score"]/span[1]/text()').extract()
        try:
            vote_num = response.xpath('//*[@class="chart_desc"]/small/span/text()').extract()[0]
        except IndexError:
            vote_num = response.xpath('//*[@class="chart_desc"]/small/span/text()').extract()
        desc = str(''.join(response.xpath('//*[@id="subject_summary"]/text()').extract())).replace('\r\n', '').replace(
            '\u3000', '').replace('\u3002', '')
        cv_list = []
        cv_count = 1
        while True:
            try:
                cv = \
                response.xpath('//*[@id="browserItemList"]/li[' + str(cv_count) + ']/div/div/span/a/text()').extract()[
                    0]
            except IndexError:
                try:
                    cv = response.xpath(
                        '//*[@id="browserItemList"]/li[' + str(cv_count + 1) + ']/div/div/span/a/text()').extract()[0]
                except IndexError:
                    break
                else:
                    cv_count = cv_count + 1
            else:
                cv_count = cv_count + 1
                cv_list.append(cv)
        # print(cv_list)
        tag_count = 1
        tag_list = []
        while True:
            try:
                tag = \
                response.xpath('//*[@id="subject_detail"]/div[3]/div/a[' + str(tag_count) + ']/span/text()').extract()[
                    0]
            except IndexError:
                break
            else:
                if tag_count < 10:
                    tag_count = tag_count + 1
                    tag_list.append(tag)
                else:
                    break
        # print(tag_list)
        bangumi_item['bangumi_id'] = str(bangumi_id)
        bangumi_item['name'] = str(name)
        bangumi_item['cover_url'] = str(cover_url)
        bangumi_item['bangumi_score'] = str(bangumi_score)
        bangumi_item['vote_num'] = str(vote_num)
        bangumi_item['episode_num'] = str(episode_num)
        bangumi_item['tags'] = list2str(tag_list)
        bangumi_item['desc'] = desc
        bangumi_item['staff_list'] = dict2str(info_dict)
        bangumi_item['cv_list'] = list2str(cv_list)
        print(bangumi_item)
        # yield bangumi_item


def list2str(l):
    result = ''
    for count, item in enumerate(l):
        result = result + item
        if count != len(l) - 1:
            result = result + ','
    return result


def dict2str(d):
    result = '{'
    for (key, value) in d.items():
        result = result + '\'' + str(key) + '\'' + ':' + '\'' + str(value) + '\'' + ','
    result = result[:-1] + '}'
    return result
