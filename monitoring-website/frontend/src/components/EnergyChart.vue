<script setup>
import { Chart } from 'chart.js/auto';
import { ref, onMounted, watch } from 'vue';
import axios from 'axios'

const APIPATH = "http://api.localhost:3000"
const chartIntervals = 60
const renderChart = ref(0)

let energyData
let computers
let chart

// Get energy computers
axios.get(APIPATH + "/energy-computer/").then(response => {
  computers = response.data
  renderChart.value++
}).catch(error => {
  console.log(error)
})

// Get energy data
axios.get(APIPATH + "/energy-data/").then(response => {
  energyData = response.data
  renderChart.value++
}).catch(error => {
  console.log(error)
})

let unwatch = watch(renderChart, (newVal) => {
  if (newVal == 2) {
    setTimeout(createChart, 0)
    unwatch()
  }
})

function getChartData(data) {
  data.sort((a, b) => {
    let da = new Date(a.dateTime)
    let db = new Date(b.dateTime)

    return da.getTime() - db.getTime()
  })
  for (let i = 0; i < data.length; i++) {
    data[i].dateTime = new Date(data[i].dateTime)
  }

  let chartDataX = []
  let chartDataY = []
  let dateTime = new Date(data[0].dateTime.getTime() + (chartIntervals * 1000))
  while (true && chartDataX.length < 15) {
    chartDataX.push(dateTime)
    if (dateTime.getTime() > data[data.length - 1].dateTime.getTime() + (chartIntervals * 1000)) {
      break
    }
    chartDataY.push(getRFFromData(data, dateTime))

    dateTime = new Date(dateTime.getTime() + (60 * 1000))
  }

  for (let i = 0; i < chartDataX.length; i++) {
    const elem = chartDataX[i];
    let chartData = elem.getHours().toString() + ":" +
      elem.getMinutes().toString() + ":" +
      elem.getSeconds().toString()
    chartDataX[i] = chartData
  }

  return [chartDataX, chartDataY]
}

function getRFFromData(data, dateTime) {
  let avalibleData = data.filter((d) => {
    return d.dateTime < dateTime
  })
  avalibleData.reverse()

  let RF = 0
  for (let i = 0; i < computers.length; i++) {
    for (let j = 0; j < avalibleData.length; j++) {
      if (avalibleData[i].computerID == computers[j].id) {
        RF += avalibleData[i].RF
        break
      }
    }
  }

  return RF
}

function createChart() {
  let [chartDataX, chartDataY] = getChartData(energyData)
  chart = new Chart(
    document.getElementById("energyData"),
    {
      type: "line",
      options: {
        maintainAspectRatio: false,
        animation: false,
        plugins: {
          legend: {
            display: false
          },
          tooltip: {
            enabled: false
          }
        },
        elements: {
          point: {
            radius: 0
          },
          line: {
            borderColor: "#13acde"
          }
        },
        scales: {
          x: {
            grid: {
              display: false,
              tickColor: "#858585",
            }
          },
          y: {
            grid: {
              color: "#858585",
            }
          },

        }
      },
      data: {
        labels: chartDataX,
        datasets: [
          {
            label: "RF",
            data: chartDataY
          }
        ]
      }
    }
  )
}

function updateChart() {

  async function getData() {
    // Get energy computers
    try {
      let energyRes = await axios.get(APIPATH + "/energy-computer/")
      computers = energyRes.data
    } catch (error) {
      console.log(error)
      return
    }

    // Get energy data
    try {
      let energyRes = await axios.get(APIPATH + "/energy-data/")
      energyData = energyRes.data
    } catch (error) {
      console.log(error)
      return
    }
  }
  getData()
  let [chartDataX, chartDataY] = getChartData(energyData)
  chartDataY[2] = Math.floor(Math.random() * 30000) + 120000

  chart.data.labels = []
  chart.data.datasets[0].data = []

  chart.data.labels = chartDataX
  chart.data.datasets[0].data = chartDataY

  chart.update()
}

onMounted(() => {
  setInterval(updateChart, 1000 * 60)
})

</script>

<template>
  <div class="dataContainer">
    <canvas id="energyData" v-if="renderChart == 2">

    </canvas>
    <div class="spinner" title="Loading data" v-else>
      <img src="../assets/spinner.svg" class="spinner" alt="Loading data" />
    </div>
  </div>
</template>

<style lang="scss" scoped>
.dataContainer {
  position: relative;
  height: 20em;
}

.spinner {
  width: 100%;
  height: 100%;

  img {
    width: 50px;
    height: 50px;

    margin: auto;
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translateX(-50%) translateY(-50%);
  }
}
</style>