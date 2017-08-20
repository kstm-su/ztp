<template>
  <div>
    <h2>{{ image.id }}: {{ image.name }}</h2>
    <div>
      <h3>description</h3>
      {{ image.description }}
    </div>
    <div>
      <h3>updated at</h3>
      {{ image.updated_at }}
    </div>
    <div>
      <h3>created at</h3>
      {{ image.created_at }}
    </div>
    <div>
      <h3>config</h3>
      <code><pre>{{ image.config }}</pre></code>
    </div>
    <div>
      <h3>nodes</h3>
      <table>
        <thead>
          <tr>
            <th>#</th>
            <th>name</th>
            <th>MAC address</th>
            <th>IP address</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="node in image.nodes">
            <td>
              <router-link :to="`/nodes/${node.id}`">
                {{ node.id }}
              </router-link>
            </td>
            <td>
              <router-link :to="`/nodes/${node.id}`">{{ node.name }}</router-link>
            </td>
            <td>{{ node.mac_address }}</td>
            <td>
              <span v-if="node.ip_address == null" class="null">null</span>
              <span v-else>{{ node.ip_address }}</span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      image: {},
    };
  },
  methods: {
    fetch() {
      let id = +this.$route.params.id;
      return this.$http.get(`/images/${id}`).then(resp => {
        this.image = resp.data;
      });
    },
  },
  mounted() {
    this.fetch();
  },
};
</script>
