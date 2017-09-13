<template>
  <div>
    <md-table>
      <md-table-header>
        <md-table-row>
          <md-table-head >#</md-table-head>
          <md-table-head>name</md-table-head>
          <md-table-head>description</md-table-head>
          <md-table-head>size</md-table-head>
          <md-table-head>status</md-table-head>
          <md-table-head>path</md-table-head>
        </md-table-row>
      </md-table-header>
      <md-table-body>
        <md-table-row v-for="image in images" :key="image.id">
          <md-table-cell>
            <router-link :to="`/images/${image.id}`">
              {{ image.id }}
            </router-link>
          </md-table-cell>
          <md-table-cell>
            <router-link :to="`/images/${image.id}`">
              {{ image.name }}
            </router-link>
          </md-table-cell>
          <md-table-cell>{{ image.description }}</md-table-cell>
          <md-table-cell>{{ image.size }}MB</md-table-cell>
          <md-table-cell>
            <span v-if="image.error">
              <md-icon class="color-red">error</md-icon>
              error
            </span>
            <span v-else-if="image.path">
              <md-icon class="color-green">check circle</md-icon>
              ready
            </span>
            <span v-else>
              <md-spinner md-indeterminate :md-size="20"></md-spinner>
              building
            </span>
          </md-table-cell>
          <md-table-cell>
            <span v-if="image.path">
              {{ image.path }}
            </span>
            <span v-else>
              <a href="#" @click.prevent="rebuild(image.id)">rebuild</a>
            </span>
          </md-table-cell>
        </md-table-row>
      </md-table-body>
    </md-table>
    <form @submit.prevent="submit">
      <md-input-container>
        <label>image name</label>
        <md-input v-model="newImage.name"></md-input>
      </md-input-container>
      <md-input-container>
        <label>LinuxKit config</label>
        <md-textarea v-model="newImage.config"></md-textarea>
      </md-input-container>
      <md-input-container>
        <label>image description</label>
        <md-textarea v-model="newImage.description"></md-textarea>
      </md-input-container>
      <md-input-container>
        <label>image size</label>
        <md-input type="number" v-model.number="newImage.size"></md-input>
      </md-input-container>
      <md-button class="md-raised md-primary" @click="submit">add</md-button>
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
      rebuild(id) {
        return this.$http.put(`/images/${id}`, {
          build: true,
        }).then(resp => this.images = this.images.map(image => {
          if (image.id === id) {
            return resp.data;
          }
          return image;
        }));
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

  .color-red {
    color: red;
  }

  .color-green {
    color: green;
  }
</style>
