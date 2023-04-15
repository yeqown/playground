<template>
  <div id="app">
    <video src=""/>
    <button @click="hdlPushingClk">开始推流</button>
  </div>
</template>

<script>
// import HelloWorld from './components/HelloWorld.vue'
// import FFmpge from 'fluent-ffmpeg' 

export default {
  name: 'App',
  components: {
    // HelloWorld
  },
  methods: {
    captureAndPlay() {
      // 老的浏览器可能根本没有实现 mediaDevices，所以我们可以先设置一个空的对象
    if (navigator.mediaDevices === undefined) {
        navigator.mediaDevices = {};
    }

    // 一些浏览器部分支持 mediaDevices。我们不能直接给对象设置 getUserMedia 
    // 因为这样可能会覆盖已有的属性。这里我们只会在没有getUserMedia属性的时候添加它。
    if (navigator.mediaDevices.getUserMedia === undefined) {
        navigator.mediaDevices.getUserMedia = function (constraints) {

            // 首先，如果有getUserMedia的话，就获得它
            var getUserMedia = navigator.webkitGetUserMedia || navigator.mozGetUserMedia;

            // 一些浏览器根本没实现它 - 那么就返回一个error到promise的reject来保持一个统一的接口
            if (!getUserMedia) {
                return Promise.reject(new Error('getUserMedia is not implemented in this browser'));
            }

            // 否则，为老的navigator.getUserMedia方法包裹一个Promise
            return new Promise(function (resolve, reject) {
                getUserMedia.call(navigator, constraints, resolve, reject);
            });
        }
    }

    navigator.mediaDevices.getUserMedia({ audio: true, video: {
        width: {ideal: 240, min: 240, max: 1080} ,
        height: {ideal: 240, min: 240, max: 1080},
      }}).then(function (stream) {
            var video = document.querySelector('video');
            // 旧的浏览器可能没有srcObject
            if ("srcObject" in video) {
                video.srcObject = stream;
            } else {
                // 防止在新的浏览器里使用它，应为它已经不再支持了
                video.src = window.URL.createObjectURL(stream);
            }
            video.onloadedmetadata = function (e) {
                console.log(e)
                video.play();
            };

            // TODO: push stream and merge
            // FFmpge({source: stream})
            // command.
        })
        .catch(function (err) {
            console.log(err.name + ": " + err.message);
        });
    },

    pushingStream() {
      console.log("pushing")
    },

    hdlPushingClk() {
      console.log("clicked pushing")
    }
  },
  mounted() {
    console.log("mounted")
    this.captureAndPlay()
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
