<template>
  <div>
    <span
      key="placeholder"
      v-if="!editable"
      class="SingleLineEditable__content"
      :class="{
        'SingleLineEditable__content--placeholder': !value,
      }"
    >
      {{ value || placeholder || $t('empty') }}
    </span>

    <span
      key="content"
      v-else
      ref="editableDiv"
      :contenteditable="true"
      class="SingleLineEditable__editable"
      @paste="handlePaste"
      @keypress="handleKeypress"
    >
      {{ value }}
    </span>

    <v-btn
      icon @click="edit()"
      v-if="canEdit && !editable"
      class="SingleLineEditable__button"
    >
      <v-icon>mdi-pencil</v-icon>
    </v-btn>

    <v-btn icon @click="save()" v-if="editable" class="ml-2 SingleLineEditable__button">
      <v-icon color="green">mdi-check</v-icon>
    </v-btn>

    <v-btn icon @click="cancel()" v-if="editable" class="SingleLineEditable__button">
      <v-icon color="red">mdi-close</v-icon>
    </v-btn>
  </div>
</template>
<style scoped>
  .SingleLineEditable__editable {
    min-width: 50px;
    outline: none;
    display: inline-block;
    background-color: rgba(128, 128, 128, 0.4);
    border-radius: 4px;
    padding-left: 5px;
    padding-right: 5px;
    margin-left: -5px;
    margin-right: -5px;
  }

  .SingleLineEditable__content {
    opacity: 0.7;
    padding-left: 5px;
    padding-right: 5px;
    margin-left: -5px;
    margin-right: -5px;
  }

  .SingleLineEditable__content--placeholder {
    color: grey;
  }

  .SingleLineEditable__button {
    margin-bottom: -7px;
    margin-top: -10px;
  }
</style>
<script>
export default {
  props: {
    value: String,
    placeholder: String,
    canEdit: Boolean,
  },
  data() {
    return {
      editable: false,
    };
  },
  watch: {
    value(newVal) {
      if (this.$refs.editableDiv.innerText !== newVal) {
        this.$refs.editableDiv.innerText = newVal || '';
      }
    },
  },
  methods: {

    handlePaste(e) {
      if (!this.editable) {
        return;
      }
      e.preventDefault();
      const text = (e.clipboardData || window.clipboardData).getData('text');
      document.execCommand('insertText', false, text.replace(/\n/g, ' '));
    },

    edit() {
      this.editable = true;
      setTimeout(() => {
        this.$refs.editableDiv.focus({ preventScroll: true });
      }, 100);
    },

    save() {
      this.editable = false;
      this.$emit('input', this.$refs.editableDiv.innerText);
      this.$emit('save');
    },

    cancel() {
      this.$refs.editableDiv.innerText = this.value;
      this.editable = false;
    },

    handleKeypress(e) {
      if (!this.editable) {
        return;
      }
      if (e.key === 'Enter') {
        e.preventDefault();
        this.save();
      }
    },
  },
};
</script>
