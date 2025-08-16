<template>
  <v-form
    ref="form"
    lazy-validation
    v-model="formValid"
    v-if="item != null"
  >
    <v-alert
      :value="formError"
      :type="(formError || '').includes('already activated') ? 'warning' : 'error'"
    >{{ formError }}
    </v-alert>

    <v-textarea
      class="mt-4"
      rows="4"
      auto-grow
      v-model="item.key"
      label="Subscription Key"
      :rules="[v => !!v || $t('key_required')]"
      required
      :disabled="formSaving"
      outlined
      dense
    ></v-textarea>

    <div style="text-align: right; margin-bottom: 30px; margin-top: -5px;">
      <v-btn @click="save" style="width: 100%;" color="primary" :disabled="formSaving">
        <v-progress-circular
          v-if="formSaving"
          indeterminate
          color="white"
          :size="24"
        ></v-progress-circular>
        <span v-else>Activate new key</span>
      </v-btn>
    </div>

    <v-card
      v-if="item.plan"
      class="mb-3"
      style="background: var(--highlighted-card-bg-color)"
    >
      <v-card-title>Plan &amp; status</v-card-title>
      <v-card-text class="pb-2">
        <v-row>
          <v-col class="py-0">
            <v-list class="py-0" style="background: unset;">
              <v-list-item class="pa-0">
                <v-list-item-content>
                  <v-list-item-title>Plan</v-list-item-title>
                  <v-list-item-subtitle>{{ item.plan }}</v-list-item-subtitle>
                </v-list-item-content>
              </v-list-item>
              <v-list-item class="pa-0">
                <v-list-item-content>
                  <v-list-item-title>Expires at</v-list-item-title>
                  <v-list-item-subtitle>{{ item.expiresAt }}</v-list-item-subtitle>
                </v-list-item-content>
              </v-list-item>
            </v-list>
          </v-col>
          <v-col class="py-0">
            <v-list class="py-0" style="background: unset;">
              <v-list-item class="pa-0">
                <v-list-item-content>
                  <v-list-item-title>Pro users</v-list-item-title>
                  <v-list-item-subtitle>{{ item.users }}</v-list-item-subtitle>
                </v-list-item-content>
              </v-list-item>
              <v-list-item class="pa-0">
                <v-list-item-content>
                  <v-list-item-title>Status</v-list-item-title>
                  <v-list-item-subtitle style="display: flex; align-items: center;">
                    <div
                      style="
                        border-radius: 100px;
                        width: 8px;
                        height: 8px;
                        background: greenyellow;
                        margin-right: 5px;
                        margin-top: 1px;
                      "
                    ></div>
                    <div>{{ item.state }}</div>
                  </v-list-item-subtitle>
                </v-list-item-content>
              </v-list-item>
            </v-list>
          </v-col>
        </v-row>

        <div style="
          margin-top: 20px;
          font-weight: bold;
          color: #00bc00;
        ">Renews in {{ (new Date() - new Date(item.expiresAt)) | formatMilliseconds }}</div>
      </v-card-text>
    </v-card>

    <div v-else class="mb-4 mt-2">

      <div>
        Don't have subscription key? <a
        target="_blank"
        href="https://portal.semaphoreui.com/auth/login?new_project=premium"
      >Get one</a>.
      </div>
    </div>

  </v-form>
</template>
<script>
import ItemFormBase from '@/components/ItemFormBase';

export default {
  mixins: [ItemFormBase],

  data() {
    return {
      tab: 0,
    };
  },

  computed: {
    isNew() {
      return false;
    },

    statusColor() {
      switch (this.item.state) {
        case 'expired':
          return 'error';
        case 'active':
          return 'success';
        default:
          return '';
      }
    },
  },

  methods: {
    async afterSave() {
      await this.loadData();
    },

    getItemsUrl() {
      return '/api/subscription';
    },

    getSingleItemUrl() {
      return '/api/subscription';
    },

    getRequestOptions() {
      return {
        method: 'post',
        url: '/api/subscription',
      };
    },
  },
};
</script>
