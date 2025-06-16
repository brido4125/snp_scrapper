from django.db import models

# Create your models here.

class Stock(models.Model):
    ticker = models.CharField(max_length=10, unique=True)
    company_name = models.CharField(max_length=255)
    sector = models.CharField(max_length=255)
    current_price = models.FloatField(null=True, blank=True)
    change_percent = models.FloatField(null=True, blank=True)
    last_updated = models.DateTimeField(auto_now=True)

    def __str__(self):
        return f"{self.company_name} ({self.ticker})"
