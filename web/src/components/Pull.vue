<script setup>


import {ref} from "vue";
import {Api} from "../api";

var images = ref([])
var message = ref("")
var error = ref("")

Api.remoteImages().then(function (r) {
  images.value = r.data
})

var push = function(name, tag) {
  error.value = ""
  message.value = ""
  Api.pull(name, tag).then(function () {
    message.value = "Container has been pulled successfully"
  }, function (e) {
    error.value = e.response.data.error
  })
}


</script>

<template>
  <h2>Images at Storj bucket</h2>
  <p class="alert alert-danger" role="alert" v-if="error">{{ error }}</p>
  <p class="alert alert-success" role="alert" v-if="message" v-html="message"></p>
  <table class="table">
    <thead>
    <tr>
      <th scope="col">Name</th>
      <th scope="col">Tag</th>
      <th scope="col">Action</th>
    </tr>
    </thead>
    <tbody>
    <tr v-for="image in images">
      <th scope="row">{{ image.Name }}</th>
      <th scope="row">{{ image.Tag }}</th>
      <th scope="row">
        <button type="button"
                class="btn btn-sm btn-primary btn-block"
                @click="push(image.Name, image.Tag)">Pull image
        </button>
      </th>
    </tr>
    </tbody>
  </table>
</template>

<style scoped>
</style>
