<script setup>
import axios from 'axios'
import { ref, watch } from 'vue'

const APIPATH = "http://localhost:3000/api"

const currentRF = ref(0)
const maxRF = ref(0)
const indicators = ref([
  "hideIndicator",
  "hideIndicator",
  "hideIndicator",
  "",
  "",
  ""
])
const percentFull = ref(0)

function getRFInfo() {
  axios.get(APIPATH + "/energy-computer/").then(response => {
    for (let i = 0; i < response.data.length; i++) {
      const computer = response.data[i];
      currentRF.value += computer.currentRF
      maxRF.value += computer.maxRF
    }

    currentRF.value /= 1000000
    currentRF.value = Math.round(currentRF.value * 100) / 100

    maxRF.value /= 1000000
    maxRF.value = Math.round(maxRF.value * 100) / 100
  })
}
getRFInfo()

// Calculate percent and show on battery indicator
function calculatePercent() {
  percentFull.value = currentRF.value / (maxRF.value / 100)
  percentFull.value = Math.round(percentFull.value * 100) / 100

  if (percentFull.value > 0) {
    indicators.value[5] = ""
  }
  else {
    indicators.value[5] = "hideIndicator"
  }
  if (percentFull.value > 16.66) {
    indicators.value[4] = ""
  }
  else {
    indicators.value[4] = "hideIndicator"
  }
  if (percentFull.value > 33.33) {
    indicators.value[3] = ""
  }
  else {
    indicators.value[3] = "hideIndicator"
  }
  if (percentFull.value > 49.99) {
    indicators.value[2] = ""
  }
  else {
    indicators.value[2] = "hideIndicator"
  }
  if (percentFull.value > 66.66) {
    indicators.value[1] = ""
  }
  else {
    indicators.value[1] = "hideIndicator"
  }
  if (percentFull.value > 83.33) {
    indicators.value[0] = ""
  }
  else {
    indicators.value[0] = "hideIndicator"
  }
}

watch(currentRF, (newVal) => {
  calculatePercent()
})

setInterval(getRFInfo, 1000 * 3)
</script>

<template>
  <div class="flex flex-row h-full">
    <div class="flex-none w-3/5">
      <div class="flex flex-col h-full">
        <div class="h-1/2 flex flex-col">
          <p class="text-center title">Network Info</p>
          <p class="text-center mt-3.5 text-2xl">Current: {{ currentRF }} RF</p>
          <p class="text-center mt-2 text-2xl">Max: {{ maxRF }} RF</p>
        </div>
        <div class="horizontalBorder h-1/2 flex flex-col">
          <p class="text-center sumRF mt-5">{{ "+100" }} {{ "k" }}RF/t</p>
          <div class="flex flex-row justify-center mt-1.5">
            <p class="w-3/6 text-xl out">Out: {{ 100 }} {{ "k" }}RF/t</p>
            <p class="in text-xl">In: {{ 200 }} {{ "k" }}RF/t</p>
          </div>
        </div>
      </div>
    </div>
    <div class="flex-1 verticalBorder" id="battery">
      <div class="flex flex-row justify-center mt-2">
        <svg class="text-center inline" width="240" height="240" viewBox="0 0 240 267" fill="none"
          xmlns="http://www.w3.org/2000/svg"><!--267-->
          <rect x="2.5" y="2.5" width="235" height="262" rx="13.5" stroke="black" stroke-width="5" />
          <rect x="15" y="217.004" width="210" height="26.5936" fill="#13ACDE" :class="indicators[5]" />
          <rect x="15" y="178.709" width="210" height="25.5299" fill="#13ACDE" :class="indicators[4]" />
          <rect x="15" y="139.351" width="210" height="26.5936" fill="#13ACDE" :class="indicators[3]" />
          <rect x="15" y="101.056" width="210" height="26.5936" fill="#13ACDE" :class="indicators[2]" />
          <rect x="15" y="62.761" width="210" height="25.5299" fill="#13ACDE" :class="indicators[1]" />
          <rect x="15" y="23.4024" width="210" height="26.5936" fill="#13ACDE" :class="indicators[0]" />
        </svg>
      </div>
      <p class="text-center text-5xl mt-2">{{ percentFull }}%</p>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use "../styles/colors.scss" as colors;

.hideIndicator {
  display: none;
}

.title {
  font-size: 2rem;
  color: colors.$primary-color;
}

.horizontalBorder {
  border-top: 8px solid black;
}

.verticalBorder {
  border-left: 8px solid black;
}

.sumRF {
  font-size: 2.6rem;
}

.out {
  color: #FF2C2C;
  //font-size: 1.2rem;
}

.in {
  color: #49FF5B;
  //font-size: 1.2rem;
}
</style>