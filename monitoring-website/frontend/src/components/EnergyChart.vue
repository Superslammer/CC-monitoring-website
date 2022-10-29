<script setup>
import { ref, onMounted } from 'vue'
import * as d3 from 'd3'

const data = [
  { date: 1666750500, amount: 10 },
  { date: 1666750620, amount: 11 },
  { date: 1666750800, amount: 10 },
  { date: 1666751100, amount: 20 },
  { date: 1666751400, amount: 21 },
  { date: 1666751700, amount: 27 },
  { date: 1666751800, amount: 29.5 },
  { date: 1666752000, amount: 30 },
  { date: 1666752300, amount: 0 },
  { date: 1666752600, amount: 30 },
  { date: 1666752900, amount: 41 },
]
const xAxisTicks = ref([])

const width = ref(500)
const height = ref(500)


onMounted(() => {
  var svgEl = document.getElementById("energyGraph")
  var svgRect = svgEl.getBoundingClientRect()

  width.value = svgRect.width
  height.value = svgRect.height

  generateGraph()
})

function generateGraph() {
  let margin = ({ top: 15, right: 15, bottom: 10, left: 15 })
  let svg = d3.select("svg")

  // Creating line
  let scaleX = d3.scaleTime()
    .domain([new Date(data[0].date * 1000), new Date(data[data.length - 1].date * 1000)])
    .range([margin.left, width.value - margin.right])

  let scaleY = d3.scaleLinear()
    .domain([d3.min(data, d => d.amount), d3.max(data, d => d.amount)])
    .range([height.value - margin.bottom, margin.top])

  let line = d3.line()
    .x(d => scaleX(d.date * 1000))
    .y(d => scaleY(d.amount))

  // Creating axises
  let xAxis = svg => svg
    .attr("transform", `translate(0,${height.value - margin.bottom})`)
    .call(d3.axisBottom(scaleX)
      .ticks(d3.timeMinute.every(5))
      .tickFormat(getXaxisTicks))
    .call(g => g.select(".domain")
      .remove())

  let yAxis = svg => svg
    .attr("transform", `translate(${margin.left},0)`)
    .call(d3.axisRight(scaleY)
      .tickSize(width.value - margin.left - margin.right)
      .tickFormat(formatYaxisTicks))
    .call(g => g.select(".domain")
      .remove())
    .call(g => g.selectAll(".tick line")
      .attr("stroke-opacity", 0.5))
    .call(g => g.selectAll(".tick text")
      .attr("x", 4)
      .attr("dy", -4))

  function getXaxisTicks(d) {
    xAxisTicks.value.push(d)
    return d
  }

  function formatYaxisTicks(d) {
    return d + " MRF"
  }

  // Applying axises and line

  svg.append("g")
    .attr("id", "remove")
    .call(xAxis)

  svg.select("#remove").remove()

  svg.append("path").attr("d", line(data))
    .classed("graphStroke", true)

  svg.append("g")
    .call(yAxis)
}
</script>

<template>
  <div>
    <svg id="energyGraph" class="line"></svg>
  </div>
</template>

<style lang="scss" scoped>
.line {
  fill: none;
}

svg {
  width: 100%;
  height: 21rem;
}
</style>