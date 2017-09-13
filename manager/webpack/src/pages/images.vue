<template>
  <div>
    <table>
      <thead>
        <tr>
          <th>#</th>
          <th>name</th>
          <th>description</th>
          <th>size</th>
          <th>status</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="image in images">
          <td>
            <router-link :to="`/images/${image.id}`">
              {{ image.id }}
            </router-link>
          </td>
          <td>
            <router-link :to="`/images/${image.id}`">
              {{ image.name }}
            </router-link>
          </td>
          <td>{{ image.description }}</td>
          <td>{{ image.size }}MB</td>
          <td>
            <span v-if="image.error">error</span>
            <span v-else-if="image.path">ready</span>
            <span v-else>building</span>
          </td>
        </tr>
      </tbody>
    </table>
    <form @submit.prevent="submit">
      <div>
        <label>
          name
          <input v-model="newImage.name" />
        </label>
      </div>
      <div>
        <label>
          config
          <textarea v-model="newImage.config"></textarea>
        </label>
      </div>
      <div>
        <label>
          description
          <textarea v-model="newImage.description"></textarea>
        </label>
      </div>
      <div>
        <label>
          size
          <input type="number" v-model.number="newImage.size" />
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
      images: [],
      newImage: {
        name: '',
        config: '',
        description: '',
        size: 1024,
      },
    };
  },
  methods: {
    fetch() {
      return this.$http.get('/images').then(resp => {
        this.images = resp.data.sort((a, b) => {
          return new Date(b.updated_at) - new Date(a.updated_at);
        });
      });
    },
    submit() {
      return this.$http.post('/images', this.newImage).then(resp => {
        this.images.unshift(resp.data);
        this.newImage = {
          name: '',
          config: '',
          description: '',
          size: 1024,
        };
      });
    },
  },
  mounted() {
    this.fetch();
  },
};
</script>

<style scoped>
.null {
  opacity: .5;
}
</style>
