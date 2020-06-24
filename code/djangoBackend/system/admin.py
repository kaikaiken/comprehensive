from django.contrib import admin
from . import models
# Register your models here.


class BangumiAdmin(admin.ModelAdmin):
    list_display = ('bangumi_id', 'name')
    list_filter = ('name',)

admin.register(models.Bangumi)
