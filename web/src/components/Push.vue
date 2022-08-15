<script setup>

import {ref} from "vue";
import {Api} from "../api";

var images = ref([])
var message = ref("")

Api.localImages().then(function (r) {
  images.value = r.data
})

var push = function(name, tag) {
  Api.push(name, tag).then(function (r) {
    message.value = "Push has been started"
  })
}

</script>

<template>
  <h2>Local images</h2>
  <table class="table">
    <thead>
    <tr>
      <th scope="col">ID</th>
      <th scope="col">Name</th>
      <th scope="col">Tag</th>
      <th scope="col">Size</th>
      <th scope="col">Action</th>
    </tr>
    </thead>
    <tbody>
    <tr v-for="image in images">
      <th scope="row">{{ image.Id }}</th>
      <th scope="row">{{ image.Name }}</th>
      <th scope="row">{{ image.Tag }}</th>
      <th scope="row">{{ image.Size }}</th>
      <th scope="row">
        <button type="button"
                class="btn btn-sm btn-primary btn-block"
                @click="push(image.Name, image.Tag)">Push image
        </button>
      </th>
    </tr>
    </tbody>
  </table>
</template>

<style scoped>
</style>
