# This is an auto-generated Django model module.
# You'll have to do the following manually to clean this up:
#   * Rearrange models' order
#   * Make sure each model has one field with primary_key=True
#   * Make sure each ForeignKey and OneToOneField has `on_delete` set to the desired behavior
#   * Remove `managed = False` lines if you wish to allow Django to create, modify, and delete the table
# Feel free to rename the models, but don't rename db_table values or field names.
from django.db import models


class Bangumi(models.Model):
    bangumi_id = models.DecimalField(primary_key=True, max_digits=10, decimal_places=0)
    name = models.CharField(max_length=25, blank=True, null=True)
    cover_url = models.CharField(max_length=100, blank=True, null=True)
    bangumi_score = models.DecimalField(max_digits=2, decimal_places=1, blank=True, null=True)
    vote_num = models.DecimalField(max_digits=10, decimal_places=0, blank=True, null=True)
    episode_num = models.DecimalField(max_digits=5, decimal_places=0, blank=True, null=True)
    tags = models.CharField(max_length=600, blank=True, null=True)
    desc = models.CharField(max_length=2000, blank=True, null=True)
    staff_list = models.CharField(max_length=3000, blank=True, null=True)
    cv_list = models.CharField(max_length=400, blank=True, null=True)

    class Meta:
        #managed = False
        db_table = 'bangumi'
        ordering = ('name', )
        verbose_name = "番剧"
        verbose_name_plural = verbose_name

    def __str__(self):
        return self.name