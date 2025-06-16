from django.urls import path
from .views import StockListView

urlpatterns = [
    path('stocks/', StockListView.as_view(), name='stock-list'),
]
