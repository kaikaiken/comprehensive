"""bangumi_project URL Configuration

The `urlpatterns` list routes URLs to views. For more information please see:
    https://docs.djangoproject.com/en/3.0/topics/http/urls/
Examples:
Function views
    1. Add an import:  from my_app import views
    2. Add a URL to urlpatterns:  path('', views.home, name='home')
Class-based views
    1. Add an import:  from other_app.views import Home
    2. Add a URL to urlpatterns:  path('', Home.as_view(), name='home')
Including another URLconf
    1. Import the include() function: from django.urls import include, path
    2. Add a URL to urlpatterns:  path('blog/', include('blog.urls'))
"""
from django.contrib import admin
from django.urls import re_path
from system import views

urlpatterns = [
    re_path(r'^bangumi/(?P<bangumi_id>[^/]+)$', views.bangumi_detail),
    re_path(r'^bangumiAll/(?P<page>[^/]+)$', views.all_bangumis),
    re_path(r'^CB/(?P<bangumi_id>[^/]+)$', views.bangumi_rec_CB),
    re_path(r'^CF/(?P<id_list>[^/]+)$', views.bangumi_rec_CF),
    re_path(r'^HotBangumiList$', views.Hot_Bangumi),
    re_path(r'^bangumiSearch/(?P<bangumi_name>[^/]+)$', views.bangumi_search),
    re_path(r'^admin/', admin.site.urls),
]
