<template xmlns:v-slot="http://www.w3.org/1999/XSL/Transform">
  <div style="overflow: hidden;" class="pb-5">

    <div class="pl-5 pt-5 d-flex" style="column-gap: 10px;">
      <div class="AnsibleServerStatus AnsibleServerStatus--ok">
        <div class="AnsibleServerStatus__count">{{ okServers }}</div>
        <div>OK SERVERS</div>
      </div>

      <div class="AnsibleServerStatus AnsibleServerStatus--bad">
        <div class="AnsibleServerStatus__count">{{ notOkServers }}</div>
        <div>NOT OK SERVERS</div>
      </div>
    </div>

    <v-btn-toggle class="pl-5 mt-8 mb-3" dense v-model="tab" mandatory>
      <v-btn value="notOkServers">
        Not ok servers
      </v-btn>
      <v-btn value="allServers">
        All servers
      </v-btn>
    </v-btn-toggle>

    <v-simple-table v-if="tab === 'notOkServers'">
      <template v-slot:default>
        <thead>
        <tr>
          <th style="width: 150px;">Server</th>
          <th style="width: 200px;">Task</th>
          <th style="calc(100% - 350px);">Error</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="(task, index) in failedTasks" :key="index">
          <td style="width: 150px;">{{ task.host }}</td>
          <td style="width: 200px;">{{ task.task }}</td>
          <td>
            <div style="overflow: hidden; color: red; max-width: 400px; text-overflow: ellipsis">
              {{ task.answer }}
            </div>
          </td>
        </tr>
        </tbody>
      </template>
    </v-simple-table>

    <v-simple-table v-else-if="tab === 'allServers'">
      <template v-slot:default>
        <thead>
        <tr>
          <th>Host</th>
          <th>Changed</th>
          <th>Failed</th>
          <th>Ignored</th>
          <th>Ok</th>
          <th>Rescued</th>
          <th>Skipped</th>
          <th>Unreachable</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="(host, index) in hosts" :key="index">
          <td>{{ host.host }}</td>
          <td>{{ host.changed }}</td>
          <td>{{ host.failed }}</td>
          <td>{{ host.ignored }}</td>
          <td>{{ host.ok }}</td>
          <td>{{ host.rescued }}</td>
          <td>{{ host.skipped }}</td>
          <td>{{ host.unreachable }}</td>
        </tr>
        </tbody>
      </template>
    </v-simple-table>
  </div>
</template>
<style lang="scss">
  .AnsibleServerStatus {
    text-align: center;
    width: 250px;
    font-weight: bold;
    color: white;
    font-size: 24px;
    line-height: 1.2;
    border-radius: 8px;
  }

  .AnsibleServerStatus__count {
    font-size: 100px;
  }

  .AnsibleServerStatus--ok {
    background-color: green;
  }

  .AnsibleServerStatus--bad {
    background-color: red;
  }
</style>

<script>

export default {
  props: {
    stages: Array,
  },

  data() {
    return {
      okServers: 0,
      notOkServers: 0,
      tab: 'notOkServers',
    };
  },

  watch: {
    stages() {
      this.calcStats();
    },
  },

  computed: {
    failedTasks() {
      const running = (this.stages || [])
        .filter((stage) => stage.type === 'running')[0];
      return running?.result.failed || {};
    },
    hosts() {
      const running = (this.stages || [])
        .filter((stage) => stage.type === 'print_result')[0];
      return running?.result.hosts || [];
    },
  },

  created() {
    this.calcStats();
  },

  methods: {
    calcStats() {
      this.hosts.forEach((host) => {
        if (host.failed > 0 || host.unreachable > 0) {
          this.notOkServers += 1;
        } else {
          this.okServers += 1;
        }
      });
    },
  },
};
</script>
