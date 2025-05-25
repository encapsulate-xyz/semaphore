<template xmlns:v-slot="http://www.w3.org/1999/XSL/Transform">
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
</template>
<style lang="scss">
</style>

<script>

export default {
  props: {
    item: Object,
    user: Object,
  },

  data() {
    return {
    };
  },

  watch: {
  },

  computed: {
  },

  created() {
  },

  methods: {
  },
};
</script>
