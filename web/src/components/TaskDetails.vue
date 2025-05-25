<template xmlns:v-slot="http://www.w3.org/1999/XSL/Transform">
  <div>

    <h3 v-if="template">Template information</h3>
    <div v-if="template" class="mb-5">
      <div>App: <b>{{ getAppTitle(template.app) }}</b></div>
      <div>
        Template: <RouterLink  :to="`/project/${projectId}/templates/${template.id}`">
          {{ template.name }}
        </RouterLink>
      </div>
    </div>

    <h3 v-if="item.commit_hash">Commit info</h3>

    <div v-if="item.commit_hash" class="mb-5">
      <div>Commit message: <b>{{ item.commit_message }}</b></div>
      <div>Commit hash: <code>{{ item.commit_hash }}</code></div>
    </div>

    <h3>Running info</h3>

    <div class="mb-5">
      <div>Message: <b>{{ item.message || '—' }}</b></div>

      <div v-if="item.user_id != null">{{ $t('author') }}: <b>{{ user?.name || '—' }}</b></div>
      <div v-else-if="item.integration_id != null">
        {{ $t('integration') }}: {{ item.integration_id }}
      </div>

      <div>{{ $t('created') }}: {{ item.created | formatDate }}</div>
      <div>{{ $t('started') }}: {{ item.start | formatDate }}</div>
      <div>{{ $t('end') }}: {{ item.end | formatDate }}</div>
      <div>{{ $t('duration') }}: {{ [item.start, item.end] | formatMilliseconds }}</div>
    </div>

    <h3 v-if="item?.params">Task parameters</h3>
    <div class="mb-5" v-if="item?.params">
      <div>Limit: {{ item.params.limit }}</div>
      <div>Debug: {{ item.params.debug }}</div>
      <div>Debug level: {{ item.params.debug_level }}</div>
      <div>Diff <code>--diff</code>: {{ item.params.diff }}</div>
      <div>Dry run <code>--check</code>: {{ item.params.dry_run }}</div>
      <div>Environment: {{ item.environment }}</div>
    </div>

  </div>
</template>
<style lang="scss">
</style>

<script>

import ProjectMixin from '@/components/ProjectMixin';
import AppsMixin from '@/components/AppsMixin';

export default {
  props: {
    item: Object,
    user: Object,
    projectId: Number,
  },

  mixins: [ProjectMixin, AppsMixin],

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
