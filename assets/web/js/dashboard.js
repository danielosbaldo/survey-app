// Dashboard Charts Initialization
(function () {
  "use strict";

  function initDashboardCharts() {
    const container = document.getElementById("dashboard-container");
    if (!container) return;

    // Destroy existing charts
    if (window.dashboardCharts) {
      window.dashboardCharts.forEach((chart) => {
        if (chart && typeof chart.destroy === "function") {
          chart.destroy();
        }
      });
    }
    window.dashboardCharts = [];

    // Get data from data attributes
    const responsesTime = container.dataset.responsesTime;
    const responsesShop = container.dataset.responsesShop;
    const questionStats = container.dataset.questionStats;

    // Parse JSON data
    const timeData = responsesTime ? JSON.parse(responsesTime) : null;
    const shopData = responsesShop ? JSON.parse(responsesShop) : null;
    const qStats = questionStats ? JSON.parse(questionStats) : null;

    // Time Chart
    if (timeData && Object.keys(timeData).length > 0) {
      const timeCtx = document.getElementById("timeChart");
      if (timeCtx) {
        const sortedDates = Object.keys(timeData).sort();
        const timeChart = new Chart(timeCtx, {
          type: "line",
          data: {
            labels: sortedDates.map((d) => {
              const date = new Date(d);
              return date.toLocaleDateString("es-ES", {
                month: "short",
                day: "numeric",
              });
            }),
            datasets: [
              {
                label: "Respuestas",
                data: sortedDates.map((d) => timeData[d]),
                borderColor: "rgb(219, 39, 119)",
                backgroundColor: "rgba(219, 39, 119, 0.1)",
                tension: 0.4,
                fill: true,
              },
            ],
          },
          options: {
            responsive: true,
            maintainAspectRatio: false,
            animation: false,
            plugins: {
              legend: { display: false },
            },
            scales: {
              y: { beginAtZero: true, ticks: { stepSize: 1 } },
            },
          },
        });
        window.dashboardCharts.push(timeChart);
      }
    }

    // Shop Chart
    if (shopData && Object.keys(shopData).length > 0) {
      const shopCtx = document.getElementById("shopChart");
      if (shopCtx) {
        const colors = [
          "rgb(219, 39, 119)",
          "rgb(59, 130, 246)",
          "rgb(16, 185, 129)",
          "rgb(168, 85, 247)",
          "rgb(245, 158, 11)",
          "rgb(239, 68, 68)",
        ];

        const shopChart = new Chart(shopCtx, {
          type: "doughnut",
          data: {
            labels: Object.keys(shopData),
            datasets: [
              {
                data: Object.values(shopData),
                backgroundColor: colors,
                borderWidth: 2,
                borderColor: "#fff",
              },
            ],
          },
          options: {
            responsive: true,
            maintainAspectRatio: false,
            animation: false,
            plugins: {
              legend: { position: "bottom" },
            },
          },
        });
        window.dashboardCharts.push(shopChart);
      }
    }

    // Question Charts
    if (qStats && Array.isArray(qStats)) {
      // qStats is now an array of {id, prompt, stats} objects
      qStats.forEach((questionStat, index) => {
        if (
          index < 4 &&
          questionStat.stats &&
          Object.keys(questionStat.stats).length > 0
        ) {
          // Match canvas by question ID
          const qCtx = document.getElementById(
            "question-chart-" + questionStat.id,
          );

          if (qCtx) {
            const questionChart = new Chart(qCtx, {
              type: "bar",
              data: {
                labels: Object.keys(questionStat.stats),
                datasets: [
                  {
                    label: "Respuestas",
                    data: Object.values(questionStat.stats),
                    backgroundColor: "rgba(219, 39, 119, 0.7)",
                    borderColor: "rgb(219, 39, 119)",
                    borderWidth: 1,
                  },
                ],
              },
              options: {
                responsive: true,
                maintainAspectRatio: false,
                animation: false,
                plugins: {
                  legend: { display: false },
                },
                scales: {
                  y: { beginAtZero: true, ticks: { stepSize: 1 } },
                },
              },
            });
            window.dashboardCharts.push(questionChart);
          }
        }
      });
    }
  }

  // Initialize on DOM ready
  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", initDashboardCharts);
  } else {
    initDashboardCharts();
  }

  // Re-initialize when HTMX swaps content
  document.body.addEventListener("htmx:afterSettle", function (evt) {
    if (evt.detail.target && evt.detail.target.id === "admin-refresh") {
      initDashboardCharts();
    }
  });
})();
