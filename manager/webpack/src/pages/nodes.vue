<template>
  <div>
    <md-table>
      <md-table-header>
        <md-table-row>
          <md-table-head>#</md-table-head>
          <md-table-head>name</md-table-head>
          <md-table-head>image</md-table-head>
          <md-table-head>MAC address</md-table-head>
          <md-table-head>IP address</md-table-head>
        </md-table-row>
      </md-table-header>
      <md-table-body>
        <md-table-row v-for="node in nodes">
          <md-table-head>
            <router-link :to="`/nodes/${node.id}`">
              {{ node.id }}
            </router-link>
          </md-table-head>
          <md-table-cell>
            <router-link :to="`/nodes/${node.id}`">{{ node.name }}</router-link>
          </md-table-cell>
          <md-table-cell>
            <span v-if="node.image == null" class="null">null</span>
            <router-link v-else :to="`/images/${node.image.id}`">{{ node.image.id }}: {{ node.image.name }}</router-link>
          </md-table-cell>
          <md-table-cell>{{ node.mac_address }}</md-table-cell>
          <md-table-cell>
            <span v-if="node.ip_address == null" class="null">null</span>
            <span v-else>{{ node.ip_address }}</span>
          </md-table-cell>
        </md-table-row>
      </md-table-body>
    </md-table>
    <form @submit.prevent="submit">
      <div>
        <label>
          name
          <input v-model="newNode.name" />
        </label>
      </div>
      <div>
        <label>
          image
          <select v-model="newNode.image_id">
            <option v-for="image in images" :value="image.id">
              {{ image.id }}: {{ image.name }}
            </option>
          </select>
        </label>
      </div>
      <div>
        <label>
          MAC address
          <input v-model="newNode.mac_address" />
        </label>
      </div>
      <div>
        <button>add</button>
      </div>
    </form>
  </div>
</template>

<script>
export default {
  data() {
    return {
      nodes: [],
      images: [],
      newNode: {
        name: '',
        image_id: 0,
        mac_address: '',
      },
    };
  },
  methods: {
    fetch() {
      return this.$http.get('/nodes').then(resp => {
        this.nodes = resp.data.sort((a, b) => {
          return new Date(b.updated_at) - new Date(a.updated_at);
        });
      });
    },
    submit() {
      return this.$http.post('/nodes', this.newNode).then(resp => {
        this.nodes.unshift(resp.data);
        this.newNode = {
          name: '',
          image_id: 0,
          mac_address: '',
        };
      });
    },
    fetchImages() {
      return this.$http.get('/images').then(resp => {
        this.images = resp.data;
      });
    },
  },
  mounted() {
    this.fetch();
    this.fetchImages();
  },
};
</script>

<style scoped>
.null {
  opacity: .5;
}
</style>
