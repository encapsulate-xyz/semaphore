import axios from 'axios';

export default {
  props: {
    projectId: Number,
  },

  methods: {
    async loadProjectResources(name) {
      return (await axios({
        method: 'get',
        url: `/api/project/${this.projectId}/${name}`,
        responseType: 'json',
      })).data;
    },
    async loadProjectResource(name, id) {
      return (await axios({
        method: 'get',
        url: `/api/project/${this.projectId}/${name}/${id}`,
        responseType: 'json',
      })).data;
    },
  },
};
