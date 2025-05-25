<template xmlns:v-slot="http://www.w3.org/1999/XSL/Transform">
  <div>

    <h3>Template information</h3>
    <div class="mb-4">
      <div>Template: {{ template?.name }}</div>
      <div>App: {{ template?.app }}</div>
    </div>

    <h3>Start information</h3>
    <v-row>
      <v-col class="pr-4">
        <v-list two-line subheader class="pa-0">
          <v-list-item class="pa-0">
            <v-list-item-content v-if="item.user_id != null">
              <v-list-item-title>{{ $t('author') }}</v-list-item-title>
              <v-list-item-subtitle>{{ user?.name || '-' }}</v-list-item-subtitle>
            </v-list-item-content>
            <v-list-item-content v-else-if="item.integration_id != null">
              <v-list-item-title>{{ $t('integration') }}</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </v-list>
      </v-col>
      <v-col class="pr-4">
        <v-list two-line subheader class="pa-0">
          <v-list-item class="pa-0">
            <v-list-item-content>
              <v-list-item-title>{{ $t('started') || '-' }}</v-list-item-title>
              <v-list-item-subtitle>
                {{ item.start | formatDate }}
              </v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>
        </v-list>
      </v-col>
      <v-col>
        <v-list-item class="pa-0">
          <v-list-item-content>
            <v-list-item-title>{{ $t('duration') || '-' }}</v-list-item-title>
            <v-list-item-subtitle>
              {{ [item.start, item.end] | formatMilliseconds }}
            </v-list-item-subtitle>
          </v-list-item-content>
        </v-list-item>
      </v-col>
    </v-row>

    <h3>Parameters</h3>

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
