from rest_framework import serializers
from .models import Stock

class StockSerializer(serializers.ModelSerializer):
    class Meta:
        model = Stock
        fields = ['id', 'ticker', 'company_name', 'sector', 'current_price', 'change_percent', 'last_updated']
