import axios from 'axios';

export default {
  props: {
    projectId: Number,
  },

  methods: {
    async loadProjectEndpoint(endpoint) {
      return (await axios({
        method: 'get',
        url: `/api/project/${this.projectId}${endpoint}`,
        responseType: 'json',
      })).data;
    },

    async loadProjectResources(name) {
      return this.loadProjectEndpoint(`/${name}`);
    },

    async loadProjectResource(name, id) {
      return this.loadProjectEndpoint(`/${name}/${id}`);
    },
  },
};
