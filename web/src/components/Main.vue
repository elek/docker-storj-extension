<script setup>

import {Api} from '../api'
import {ref} from "vue";

const status = ref("unknown")
var message = ref("")
var error = ref("")

var refresh = function (){
  Api.status().then(function (s) {
    status.value = s.status
  })
}
refresh()

var stop = function () {
  Api.stop().then(function () {
    message.value = "Container is beeing stopped"
    refresh()
  }, function (e) {
    error = e
  })
}

var start = function () {
  Api.start().then(function () {
    message.value = "Container is starting"
    refresh()
  }, function (e) {
    error = e
  })
}

</script>

<template>
  <div class="container">
    <p>This extension helps to start a local docker registry to push/pull images from/to <a href="https://storj.io">Storj
      Decentralized cloud</a></p>

    <div v-if="status === 'unkown'">
      <p>Checking the state of the registry container. </p>
    </div>

    <div v-if="status === 'missing'">
      <p>Registry is not yet configured. Please configured it under the
        <router-link to="/config">configuration menu</router-link>
      </p>
    </div>

    <div v-if="status === 'stopped'">
      <p>Registry is configured, but not running.</p>
      <button type="button"
              class="btn btn-lg btn-warning btn-block"
              @click="start()">Start registry
      </button>
    </div>

    <div v-if="status === 'running'">
      <p>Registry is up and running. You can use with with the prefix <code>localhost:9999</code>.</p>
      <button type="button"
              class="btn btn-lg btn-warning btn-block"
              @click="stop()">Stop registry
      </button>
    </div>


  </div>

</template>

<style scoped>
</style>
