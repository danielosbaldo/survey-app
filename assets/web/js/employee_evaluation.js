// Employee Evaluation Charts Initialization
(function () {
  "use strict";

  function initEmployeeEvaluationCharts() {
    const container = document.getElementById("employee-evaluation-container");
    if (!container) return;

    // Destroy existing charts
    if (window.employeeEvaluationCharts) {
      window.employeeEvaluationCharts.forEach((chart) => {
        if (chart && typeof chart.destroy === "function") {
          chart.destroy();
        }
      });
    }
    window.employeeEvaluationCharts = [];

    // Get data from data attributes
    const responsesTime = container.dataset.responsesTime;
    const questionAverages = container.dataset.questionAverages;

    // Parse JSON data
    const timeData = responsesTime ? JSON.parse(responsesTime) : null;
    const avgData = questionAverages ? JSON.parse(questionAverages) : null;

    // Time Chart - Evaluations Over Time
    if (timeData && Object.keys(timeData).length > 0) {
      const timeCtx = document.getElementById("evaluationTimeChart");
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
                label: "Evaluaciones",
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
        window.employeeEvaluationCharts.push(timeChart);
      }
    }

    // Amabilidad Chart
    if (avgData && Array.isArray(avgData) && avgData.length > 0) {
      // Find the amabilidad question
      const amabilidadQuestion = avgData.find(
        (q) =>
          q.prompt.includes("amabilidad") || q.prompt.includes("Amabilidad"),
      );

      const amabilidadCtx = document.getElementById("amabilidadChart");

      if (
        amabilidadCtx &&
        amabilidadQuestion &&
        amabilidadQuestion.stats &&
        Object.keys(amabilidadQuestion.stats).length > 0
      ) {
        // Sort the stats by key (1-5) for proper ordering
        const sortedKeys = Object.keys(amabilidadQuestion.stats).sort(
          (a, b) => {
            return parseInt(a) - parseInt(b);
          },
        );

        const amabilidadChart = new Chart(amabilidadCtx, {
          type: "bar",
          data: {
            labels: sortedKeys,
            datasets: [
              {
                label: "Evaluaciones",
                data: sortedKeys.map((key) => amabilidadQuestion.stats[key]),
                backgroundColor: [
                  "rgba(239, 68, 68, 0.7)", // Red for 1
                  "rgba(251, 146, 60, 0.7)", // Orange for 2
                  "rgba(250, 204, 21, 0.7)", // Yellow for 3
                  "rgba(132, 204, 22, 0.7)", // Light green for 4
                  "rgba(34, 197, 94, 0.7)", // Green for 5
                ],
                borderColor: [
                  "rgb(239, 68, 68)",
                  "rgb(251, 146, 60)",
                  "rgb(250, 204, 21)",
                  "rgb(132, 204, 22)",
                  "rgb(34, 197, 94)",
                ],
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
              y: {
                beginAtZero: true,
                ticks: { stepSize: 1 },
                title: {
                  display: true,
                  text: "Cantidad de respuestas",
                },
              },
              x: {
                title: {
                  display: true,
                  text: "Calificación (1 = Malo, 5 = Excelente)",
                },
              },
            },
          },
        });
        window.employeeEvaluationCharts.push(amabilidadChart);
      }
    }
  }

  // Initialize on DOM ready
  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", initEmployeeEvaluationCharts);
  } else {
    initEmployeeEvaluationCharts();
  }

  // Re-initialize when HTMX swaps content
  document.body.addEventListener("htmx:afterSettle", function (evt) {
    if (evt.detail.target && evt.detail.target.id === "admin-refresh") {
      initEmployeeEvaluationCharts();
    }
  });
})();
