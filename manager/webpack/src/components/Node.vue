<template>
  <div>
    <h2>{{ node.id }}: {{ node.name }}</h2>
    <div>
      <h3>MAC address</h3>
      {{ node.mac_address }}
    </div>
    <div>
      <h3>IP address</h3>
      {{ node.ip_address }}
    </div>
    <div>
      <h3>updated at</h3>
      {{ node.updated_at }}
    </div>
    <div>
      <h3>created at</h3>
      {{ node.created_at }}
    </div>
    <div>
      <h3>image</h3>
      <router-link :to="`/images/${node.image.id}`">
        {{ node.image.id }}: {{ node.image.name }}
      </router-link>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      node: {},
    };
  },
  methods: {
    fetch() {
      let id = +this.$route.params.id;
      return this.$http.get(`/nodes/${id}`).then(resp => {
        this.node = resp.data;
      });
    },
  },
  mounted() {
    this.fetch();
  },
};
</script>
