<template>
  <div>
    <table>
      <thead>
        <tr>
          <th>#</th>
          <th>name</th>
          <th>image</th>
          <th>MAC address</th>
          <th>IP address</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="node in nodes">
          <td>
            <router-link :to="`/nodes/${node.id}`">
              {{ node.id }}
            </router-link>
          </td>
          <td>
            <router-link :to="`/nodes/${node.id}`">{{ node.name }}</router-link>
          </td>
          <td>
            <span v-if="node.image == null" class="null">null</span>
            <router-link v-else :to="`/images/${node.image.id}`">{{ node.image.id }}: {{ node.image.name }}</router-link>
          </td>
          <td>{{ node.mac_address }}</td>
          <td>
            <span v-if="node.ip_address == null" class="null">null</span>
            <span v-else>{{ node.ip_address }}</span>
          </td>
        </tr>
      </tbody>
    </table>
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
