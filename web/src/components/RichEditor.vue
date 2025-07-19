<template>
  <div>
    <v-dialog
      v-model="envEditorDialog"
      max-width="800"
      persistent
      :transition="false"
    >
      <div style="position: relative;">
        <codemirror
          v-if="envEditorDialog"
          class="EnvironmentMaximizedEditor"
          :style="{ border: '1px solid lightgray' }"
          v-model="text"
          :options="cmOptions"
          :placeholder="$t('enterExtraVariablesJson')"
        />

        <v-btn
          dark
          fab
          small
          color="blue-grey"
          style="
            position: absolute;
            right: 0;
            top: 0;
            margin: 10px;
          "
          @click="closeDialog()"
        >
          <v-icon>mdi-arrow-collapse</v-icon>
        </v-btn>
      </div>
    </v-dialog>

    <v-btn
      dark
      fab
      small
      color="blue-grey"
      @click="envEditorDialog = true"
    >
      <v-icon>mdi-arrow-expand</v-icon>
    </v-btn>

  </div>
</template>

<script>
/* eslint-disable import/no-extraneous-dependencies,import/extensions */
import { codemirror } from 'vue-codemirror';
import 'codemirror/lib/codemirror.css';
import 'codemirror/mode/vue/vue.js';
import 'codemirror/addon/display/placeholder.js';

export default {
  props: {
    value: String,
  },

  components: {
    codemirror,
  },

  watch: {
    envEditorDialog(val) {
      this.$emit('maximize', {
        maximized: val,
      });
    },

    value() {
      this.text = this.value;
    },
  },

  created() {
    this.text = this.value;
  },

  data() {
    return {

      cmOptions: {
        tabSize: 2,
        mode: 'application/json',
        lineNumbers: true,
        line: true,
        lint: true,
        indentWithTabs: false,
      },

      text: null,
      envEditorDialog: false,
    };
  },

  methods: {
    closeDialog() {
      if (this.text !== this.value) {
        this.$emit('input', this.text);
      }
      this.envEditorDialog = false;
    },
  },
};
</script>
<style lang="scss">
.EnvironmentMaximizedEditor {
  .CodeMirror {
    font-size: 14px;
    height: 600px !important;
  }
}
</style>
