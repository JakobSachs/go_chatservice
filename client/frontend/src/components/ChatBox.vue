<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "ChatBox",
  props: {
    text: String,
  },
  model: {
    prop: "text",
    event: "input",
  },
  data() {
    return {
      text: "",
    };
  },
});
</script>
<template>
  <div class="chat-box">
    <input type="text" v-model="text" />
    <button :disabled="!text">Send</button>
  </div>
</template>

<style scoped></style>
