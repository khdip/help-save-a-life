fetch('/json/dailyReport')
  .then(res => res.json())
  .then(data => {
    const currencyMap = {};

    // Organize data by currency
    data.forEach(item => {
      if (!currencyMap[item.Currency]) {
        currencyMap[item.Currency] = {};
      }
      currencyMap[item.Currency][item.Date] = item.Amount;
    });

    // Get all unique dates
    const allDates = [...new Set(data.map(item => item.Date))].sort();

    // Prepare datasets
    const datasets = Object.entries(currencyMap).map(([currency, dateMap]) => {
      const amounts = allDates.map(date => dateMap[date] || 0);
      return {
        label: currency,
        data: amounts,
        borderColor: getRandomColor(),
        fill: false,
        tension: 0.4
      };
    });

    // Render chart
    new Chart(document.getElementById('dailyReportChart'), {
      type: 'line',
      data: {
        labels: allDates,
        datasets: datasets
      },
      options: {
        responsive: true,
        plugins: {
          legend: { position: 'top' },
          title: { display: false }
        }
      }
    });
  });

function getRandomColor() {
  return `hsl(${Math.floor(Math.random() * 360)}, 70%, 50%)`;
}
