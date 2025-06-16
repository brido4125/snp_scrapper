import { useEffect, useState } from 'react'

interface Stock {
  ticker: string;
  company_name: string;
  sector: string;
  current_price: number | null;
  change_percent: number | null;
  last_updated: string;
}

function App() {
  const [stocks, setStocks] = useState<Stock[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchStocks = async () => {
      try {
        const response = await fetch('/api/stocks/');
        if (!response.ok) {
          throw new Error('Failed to fetch stocks');
        }
        const data = await response.json();
        setStocks(data);
      } catch (err) {
        setError((err as Error).message);
      } finally {
        setLoading(false);
      }
    };

    fetchStocks();
  }, []);

  if (loading) return <div className="flex justify-center items-center h-screen">Loading...</div>;
  if (error) return <div className="flex justify-center items-center h-screen text-red-500">Error: {error}</div>;

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4 text-center">S&P 500 Stocks</h1>
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
        {stocks.map((stock) => (
          <div key={stock.ticker} className="bg-white p-4 rounded-lg shadow-md">
            <h2 className="text-xl font-semibold">{stock.ticker}</h2>
            <p className="text-gray-600">{stock.company_name}</p>
            <p className="text-gray-500">{stock.sector}</p>
            <p className="text-green-500">
              {stock.current_price != null ? `$${stock.current_price.toFixed(2)}` : '-'}
            </p>
            <p
              className={`${typeof stock.change_percent === 'number' ? (stock.change_percent >= 0 ? 'text-green-500' : 'text-red-500') : 'text-gray-500'}`}
            >
              {typeof stock.change_percent === 'number' ? `${stock.change_percent.toFixed(2)}%` : '-'}
            </p>
            <p className="text-gray-400 text-sm">Updated: {new Date(stock.last_updated).toLocaleString()}</p>
          </div>
        ))}
      </div>
    </div>
  );
}

export default App;
