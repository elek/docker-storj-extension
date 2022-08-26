<script setup>

import {Api} from '../api'
import {ref} from "vue";

const bucket = ref("")
const grant = ref("")
var message = ref("")
var error = ref("")
var save = function () {
  message.value = ""
  error.value = ""
  if (bucket.value === "") {
    error.value = "Bucket  is required"
    return
  }
  if (grant.value === "") {
    error.value = "Grant  is required"
    return
  }
  Api.configure(bucket.value, grant.value).then(function () {
    message.value = "Configuration has been saved successfully"
  }, function (e) {
    error.value = e.response.data.error

  })
}
</script>

<template>

  <form>
    <p>You need a bucket name and and access grant to use Storj Decentralized cloud</p>
    <p class="alert alert-danger" role="alert" v-if="error">{{ error }}</p>
    <p class="alert alert-success" role="alert" v-if="message" v-html="message"></p>

    <input v-model="bucket" type="text" id="bucket" class="mb-3 form-control"
           placeholder="Bucket name"
           required autofocus>

    <input v-model="grant" type="text" id="grant" class="mb-3 form-control"
           placeholder="Access grant"
           required>

    <button type="button"
            class="btn btn-lg btn-warning btn-block"
            @click="save()">Save configuration
    </button>
  </form>
</template>

<style scoped>
</style>
