import pandas as pd
import requests
import yfinance as yf
import numpy as np
from django.core.management.base import BaseCommand
from stocks.models import Stock
from django.db import transaction

class Command(BaseCommand):
    help = 'Updates the S&P 500 stock list from Wikipedia and fetches price data.'

    def handle(self, *args, **options):
        # 1. Fetch S&P 500 list from Wikipedia
        self.stdout.write("Fetching S&P 500 list from Wikipedia...")
        try:
            url = 'https://en.wikipedia.org/wiki/List_of_S%26P_500_companies'
            headers = {
                'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3'
            }
            response = requests.get(url, headers=headers)
            response.raise_for_status()
            tables = pd.read_html(response.text)
            sp500_table = tables[0]
            self.stdout.write(self.style.SUCCESS("Successfully fetched S&P 500 list."))
        except Exception as e:
            self.stdout.write(self.style.ERROR(f"Failed to fetch S&P 500 list: {e}"))
            return

        # 2. Get all tickers and fetch data from yfinance
        tickers = sp500_table['Symbol'].unique().tolist()
        if not tickers:
            self.stdout.write(self.style.ERROR("No tickers found."))
            return
            
        self.stdout.write(f"Fetching price data for {len(tickers)} tickers...")
        data = yf.download(tickers=tickers, period='2d', group_by='ticker')
        self.stdout.write(self.style.SUCCESS("Successfully fetched price data."))

        # 3. Update database
        self.stdout.write("Updating database...")
        with transaction.atomic():
            for index, row in sp500_table.iterrows():
                ticker_symbol = row['Symbol']
                
                stock, created = Stock.objects.update_or_create(
                    ticker=ticker_symbol,
                    defaults={
                        'company_name': row['Security'],
                        'sector': row['GICS Sector']
                    }
                )
                
                try:
                    ticker_data = data[ticker_symbol]
                    if not ticker_data.empty and len(ticker_data) > 1:
                        # Use last two days to calculate change
                        previous_close = ticker_data['Close'].iloc[-2]
                        current_price = ticker_data['Close'].iloc[-1]
                        
                        if previous_close > 0:
                            change_percent = ((current_price - previous_close) / previous_close) * 100
                        else:
                            change_percent = 0.0

                        # Replace NaN with None for JSON compatibility
                        stock.current_price = None if np.isnan(current_price) else current_price
                        stock.change_percent = None if np.isnan(change_percent) else change_percent
                        stock.save()

                except KeyError:
                    self.stdout.write(self.style.WARNING(f"No data for {ticker_symbol}, skipping price update."))
                except Exception as e:
                    self.stdout.write(self.style.ERROR(f"Error processing {ticker_symbol}: {e}"))

        self.stdout.write(self.style.SUCCESS('Successfully updated S&P 500 stock data.'))
