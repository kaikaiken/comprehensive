# -*- coding:utf-8 -*-
from django.shortcuts import render, HttpResponse
from django.http import JsonResponse
from .models import Bangumi
import json
import pandas as pd
import numpy as np

from sklearn.metrics.pairwise import cosine_similarity
import sklearn
import json
from sklearn.feature_extraction.text import CountVectorizer
import pymysql as sql
import numpy as np
import requests
from django.contrib.staticfiles import finders

from fuzzywuzzy import fuzz
import math
# Create your views here.


def all_bangumis(request, page):
    page = int(page)
    bangumis = Bangumi.objects.all().order_by('-vote_num')
    data = []
    for i, bangumi in enumerate(bangumis[page*20: page*20 + 20]):
        result = {}
        result["bangumi_id"] = int(bangumi.bangumi_id)
        result["cover_url"] = bangumi.cover_url
        result["bangumi_score"] = float(bangumi.bangumi_score)
        if bangumi.name[0:2] == '[\'':
            bangumi.name = str(bangumi.name[2:-2])
        result["name"] = bangumi.name
        data.append(result)
    data = json.dumps(data, ensure_ascii=False)
    return HttpResponse(data)


def bangumi_detail(request, bangumi_id):
    bangumi = Bangumi.objects.get(pk=bangumi_id)
    data = []
    result = {}
    result["bangumi_id"] = int(bangumi_id)
    result["name"] = bangumi.name
    result["cover_url"] = bangumi.cover_url
    result["bangumi_score"] = float(bangumi.bangumi_score)
    result["vote_num"] = int(bangumi.vote_num)
    result["episode_num"] = int(bangumi.episode_num)
    result["tags"] = bangumi.tags
    result["desc"] = bangumi.desc
    result["staff_list"] = bangumi.staff_list
    result["cv_list"] = bangumi.cv_list
    print(result)
    data.append(result)
    data = json.dumps(data, ensure_ascii=False)
    return HttpResponse(data)


def bangumi_search(request, bangumi_name):
    bangumiAll = Bangumi.objects.all()
    #json_result = {}
    result_list = []
    for bangumi in bangumiAll:
        result = {}
        result["bangumi_id"] = int(bangumi.bangumi_id)
        result["cover_url"] = bangumi.cover_url
        result["bangumi_score"] = float(bangumi.bangumi_score)
        result["desc"] = bangumi.desc
        result["staff_list"] = bangumi.staff_list
        if bangumi.name[0:2] == '[\'':
            bangumi.name = str(bangumi.name[2:-2])
        result["name"] = bangumi.name
        value = fuzz.token_sort_ratio(bangumi_name, bangumi.name)
        #result = json.dumps(result, ensure_ascii=False)
        if len(result) < 20 and value >= 40:
            result_list.append(result)
    #son_result["data"] = result_list
    result_list = json.dumps(result_list, ensure_ascii=False)
    return HttpResponse(result_list)


def bangumi_rec_CB(request, bangumi_id):
    conn = sql.connect(host="localhost", user="root", password="test", db = "bangumi_project", charset = "utf8")
    sql_query = 'SELECT * FROM bangumi'
    df = pd.read_sql(sql_query, con=conn, index_col='bangumi_id')
    df['feature'] = ''
    count = CountVectorizer()
    create_feature(df)
    #print(df['feature'])
    count_matrix = count.fit_transform(df['feature'])

    cosine_sim = cosine_similarity(count_matrix, count_matrix)
    #print(cosine_sim)
    indices = pd.Series(df.index)
    recommended_bangumis = []
    idx = indices[indices == int(bangumi_id)].index[0]
    score_series = pd.Series(cosine_sim[idx]).sort_values(ascending = False)
    top_10_indexes = list(score_series.iloc[1:11].index)
    top_10_values = list(score_series.iloc[1:11].values)
    for i in top_10_indexes:
        recommended_bangumis.append(list(df.index)[i])
    recommended_result = {}
    for (i, bangumi) in enumerate(recommended_bangumis):
        recommended_result[str(int(bangumi))] = top_10_values[i]
    data = []
    for id in recommended_result:
        try:
            bangumi = Bangumi.objects.get(pk=id)
            result = {}
            result["bangumi_id"] = id
            result["cover_url"] = bangumi.cover_url
            result["bangumi_score"] = float(bangumi.bangumi_score)
            if bangumi.name[0:2] == '[\'':
                bangumi.name = str(bangumi.name[2:-2])
            result["name"] = bangumi.name
            result["simi_value"] = recommended_result[id]
            data.append(result)
        except:
            print(id)
    data = json.dumps(data, ensure_ascii=False)
    return HttpResponse(data)
    conn.close()





def create_feature(df):
    for i in range(len(df['name'].array)):
        tags = df['tags'].array[i].split(',')
        cv_list = df['cv_list'].array[i].split(',')
        if len(cv_list) > 6:
            cv_list = cv_list[:6]
        try:
            staff_dict = df['staff_list'].array[i].replace('\t', '').replace("\"", "\'").replace("\':\'", "\":\"").replace("\',\'", "\",\"")
            staff_dict = staff_dict.replace("{\'", "{\"").replace("\'}", "\"}")
            staff_dict = json.loads(staff_dict)
        except:
            print(df['staff_list'].array[i])
        feature = tags + cv_list
        staff_feature = ["原作", "导演", "音乐", "原画"]
        for s in staff_feature:
            staff = staff_dict.get(s, '').split(',')[0]
            if len(staff) > 0:
                feature.append(staff)

        feature = ','.join(set(feature))
        #print(feature)
        df.iloc[i, 9] = feature


def bangumi_rec_CF(request, id_list):
    id_list = id_list.split(',')
    user_fav_dict = {}
    for (i, bangumi_id) in enumerate(id_list):
        print(bangumi_id)
        user_fav_dict[bangumi_id] = 1
    result = CF_recommend(user_fav_dict)
    print("####result:" + str(result))
    data = []
    for id in result:
        try:
            bangumi = Bangumi.objects.get(pk=id)
            result = {}
            result["bangumi_id"] = id
            result["cover_url"] = bangumi.cover_url
            result["bangumi_score"] = float(bangumi.bangumi_score)
            if bangumi.name[0:2] == '[\'':
                bangumi.name = str(bangumi.name[2:-2])
            result["name"] = bangumi.name
            #result["simi_value"] = result['\'' + id + '\'']
            data.append(result)
        except:
            print("***" + id)

    data = json.dumps(data, ensure_ascii=False)
    return HttpResponse(data)


def Hot_Bangumi(request):
    train = dict()  # 用户-物品的矩阵
    csv = finders.find('data_new.csv')
    for line in open(csv):
        user, score, item = line.strip().split(",")
        train.setdefault(user, {})
        train[user][item] = int(float(score))
    '''
    for k, v in train.items():
        print("Key: " + k)
        print("Value: " + str(v) + '\n')
    '''
    # 通过番剧被不同用户收藏的次数建立物品-物品矩阵
    C = dict()  # 物品-物品的共现矩阵
    N = dict()  # 番剧被多少个不同用户喜爱
    for user, items in train.items():
        for i in items.keys():
            N.setdefault(i, 0)
            N[i] += 1
            C.setdefault(i, {})
            for j in items.keys():
                if i == j: continue
                C[i].setdefault(j, 0)
                C[i][j] += 1
                # 计算相似度矩阵
    # 计算余弦相似度
    W = dict()
    for i, related_items in C.items():
        W.setdefault(i, {})
        for j, cij in related_items.items():
            W[i][j] = cij / (math.sqrt(N[i] * N[j]))
    # for k, v in W.items():
    # print("Key: " + k)
    # print("Value: " + str(v) + '\n')
    result = ''
    for k, v in W.items():
        result += str(int(k))
        result += ','
    result = result[:-1]
    return HttpResponse(result)


def CF_recommend(user_fav_dict):
    print("dict:" + str(user_fav_dict))
    train = dict()  # 用户-物品的矩阵
    csv = finders.find('data_new.csv')
    for line in open(csv):
        user, score, item = line.strip().split(",")
        train.setdefault(user, {})
        train[user][item] = int(float(score))
    '''
    for k, v in train.items():
        print("Key: " + k)
        print("Value: " + str(v) + '\n')
    '''
    # 通过番剧被不同用户收藏的次数建立物品-物品矩阵
    C = dict()  # 物品-物品的共现矩阵
    N = dict()  # 番剧被多少个不同用户喜爱
    for user, items in train.items():
        for i in items.keys():
            N.setdefault(i, 0)
            N[i] += 1
            C.setdefault(i, {})
            for j in items.keys():
                if i == j: continue
                C[i].setdefault(j, 0)
                C[i][j] += 1
                # 计算相似度矩阵
    '''
    for k, v in N.items():
        print("Key: " + k)
        print("Value: " + str(v) + '\n')
    for k, v in C.items():
        print("Key: " + k)
        print("Value: " + str(v) + '\n')
    '''
    # 计算余弦相似度
    W = dict()
    for i, related_items in C.items():
        W.setdefault(i, {})
        for j, cij in related_items.items():
            W[i][j] = cij / (math.sqrt(N[i] * N[j]))
    #for k, v in W.items():
        #print("Key: " + k)
        #print("Value: " + str(v) + '\n')
    hot_bangumis = []
    for k, v in W.items():
        hot_bangumis.append(k)
    print(len(hot_bangumis))
    rank = dict()
    #action_item = train['20728']  # 用户user产生过行为的item和评分
    #print(str(train['20728']))
    print(str(user_fav_dict))
    to_delete_list = []
    for k, v in user_fav_dict.items():
        if k not in hot_bangumis:
            to_delete_list.append(k)
    for k in to_delete_list:
        user_fav_dict.pop(k)
    print(str(user_fav_dict))
    for item, score in user_fav_dict.items():
        for j, wj in sorted(W[item].items(), key=lambda x: x[1], reverse=True)[0:3]:
            if j in user_fav_dict.keys():
                continue
            rank.setdefault(j, 0)
            rank[j] += score * wj
    final = dict(sorted(rank.items(), key=lambda x: x[1], reverse=True)[0:10])
    return final

