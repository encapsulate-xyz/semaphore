<template xmlns:v-slot="http://www.w3.org/1999/XSL/Transform">
  <div>

    <h3>Template information</h3>
    <div class="mb-4">
      <div>Template: {{ template?.name }}</div>
      <div>App: {{ template?.app }}</div>
    </div>

    <h3>Commit info</h3>

    <div class="mb-4">
      <div>Commit message: {{ item.commit_message }}</div>
      <div>Commit hash: {{ item.commit_hash }}</div>
    </div>

    <h3>Running info</h3>

    <div class="mb-4">
      <div>Message: {{ item.message || '-' }}</div>

      <div v-if="item.user_id != null">{{ $t('author') }}: {{ user?.name || '-' }}</div>
      <div v-else-if="item.integration_id != null">
        {{ $t('integration') }}: {{ item.integration_id }}
      </div>

      <div>{{ $t('started') }}: {{ item.start | formatDate }}</div>
      <div>{{ $t('duration') }}: {{ [item.start, item.end] | formatMilliseconds }}</div>
    </div>

    <h3>Task parameters</h3>
    <div class="mb-4">
      <div>Limit: {{ item.params.limit }}</div>
      <div>Debug: {{ item.params.debug }}</div>
      <div>Debug level: {{ item.params.debug_level }}</div>
      <div>Diff: {{ item.params.diff }}</div>
      <div>Environment: {{ item.enviroment }}</div>
    </div>

  </div>
</template>
<style lang="scss">
</style>

<script>

import ProjectMixin from '@/components/ProjectMixin';

export default {
  props: {
    item: Object,
    user: Object,
    projectId: Number,
  },

  mixins: [ProjectMixin],

  data() {
    return {
      template: null,
    };
  },

  watch: {
    async item() {
      if (this.item?.template_id !== this.template?.id) {
        await this.loadData();
      }
    },
  },

  computed: {
  },

  async created() {
    await this.loadData();
  },

  methods: {
    async loadData() {
      this.template = await this.loadProjectResource('templates', this.item.template_id);
    },
  },
};
</script>
