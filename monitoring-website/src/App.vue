<script setup>
import { ref } from 'vue'
import EnergyChart from './components/EnergyChart'

</script>

<template>
  <div>
    <div class="mb-5 top"></div>
    <div class="grid grid-cols-5 gap-10 mr-10 ml-10">
      <div class="col-span-2 h-100 mb-2 card">

      </div>
      <div class="col-span-3 h-100 mb-2 card">
        <EnergyChart />
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use "./styles/colors.scss" as colors;

.top {
  background-color: colors.$primary-color;
  height: 4vh;
}

.card{
  background-color: colors.$card-color
}
</style>