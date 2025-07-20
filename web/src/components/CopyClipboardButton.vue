<template>
  <div>
    <textarea
      style="position: absolute; left: -9999px; top: -9999px;"
      ref="copy_to_clipboard_textarea"
    ></textarea>
    <v-btn
      icon
      color="white"
      @click="copy()"
    >
      <v-icon>mdi-content-copy</v-icon>
    </v-btn>
  </div>
</template>

<script>
import EventBus from '@/event-bus';

export default {
  props: {
    text: String,
    successMessage: {
      type: String,
      default: 'Text copied to clipboard!'
    },
  },
  methods: {
    copy() {
      try {
        const el = this.$refs.copy_to_clipboard_textarea;
        el.value = this.text;
        el.focus();
        el.select();
        const successful = document.execCommand('copy');

        if (!successful) {
          throw new Error('Fallback copy failed');
        }

        EventBus.$emit('i-snackbar', {
          color: 'success',
          text: this.successMessage
        });
      } catch (e) {
        EventBus.$emit('i-snackbar', {
          color: 'error',
          text: `Can't copy to clipboard: ${e.message}`
        });
      }
    }
  }
};
</script>
