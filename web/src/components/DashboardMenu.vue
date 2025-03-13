<template>
  <div>
    <v-tabs class="pl-4">
      <v-tab
        v-if="projectType === ''"
        key="history"
        :to="`/project/${projectId}/history`"
      >{{ $t('history') }}
      </v-tab>

      <v-tab key="stats" :to="`/project/${projectId}/stats`">{{ $t('stats') }}</v-tab>

      <v-tab key="activity" :to="`/project/${projectId}/activity`">{{ $t('activity') }}</v-tab>

      <v-tab
        v-if="canUpdateProject"
        key="settings"
        :to="`/project/${projectId}/settings`"
      >{{ $t('settings') }}
      </v-tab>

      <v-tab
        v-if="projectType === ''"
        key="runners"
        :to="`/project/${projectId}/runners`"
      >
        {{ $t('runners') }}
        <v-icon class="ml-1" large color="hsl(348deg, 86%, 61%)">mdi-professional-hexagon</v-icon>
      </v-tab>
    </v-tabs>

    <v-divider style="margin-top: -1px;" />
  </div>
</template>
<script>
import PermissionsCheck from '@/components/PermissionsCheck';
import {
  TEMPLATE_TYPE_ACTION_TITLES,
  TEMPLATE_TYPE_ICONS,
  TEMPLATE_TYPE_TITLES,
} from '@/lib/constants';

export default {

  mixins: [PermissionsCheck],

  props: {
    projectId: Number,
    projectType: String,
    canUpdateProject: Boolean,
  },

  data() {
    return {
      dateRanges: [{
        text: 'Past week',
        value: 'last_week',
      }, {
        text: 'Past month',
        value: 'last_month',
      }, {
        text: 'Past year',
        value: 'last_year',
      }],
      users: [{
        text: 'All users',
        value: null,
      }],
      user: null,
      TEMPLATE_TYPE_ICONS,
      TEMPLATE_TYPE_TITLES,
      TEMPLATE_TYPE_ACTION_TITLES,
      stats: null,
      dateRange: 'last_week',
    };
  },
};
</script>
