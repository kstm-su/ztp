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
      <transition-group name="list" tag="md-table-body">
        <md-table-row v-for="image in sortedImages" :key="image.id" class="list-item">
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
            <transition name="status" mode="out-in">
              <span v-if="image.error" key="error">
                <md-icon class="color-red">error</md-icon>
                error
              </span>
              <span v-else-if="image.path" key="ready">
                <md-icon class="color-green">check_circle</md-icon>
                ready
              </span>
              <span v-else key="building">
                <md-spinner md-indeterminate :md-size="20"></md-spinner>
                building
              </span>
            </transition>
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
      </transition-group>
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
    computed: {
      sortedImages() {
        return this.images.sort((a, b) => {
          new Date(b.updated_at) - new Date(a.updated_at);
        });
      },
    },
    sockets: {
      image(image) {
        this.images = this.images.map(i => {
          if (i.id === image.id) {
            return image;
          }
          return i;
        });
      },
    },
    methods: {
      fetch() {
        return this.$http.get('/images').then(resp => {
          this.images = resp.data;
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
        }).then(resp => this.images.map(image => {
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
  .list-item {
    transition: all ease .5s;
  }

  .list-enter,
  .list-leave-active {
    opacity: 0;
  }

  .status-enter-active,
  .status-leave-active {
    transition: opacity ease .25s;
  }

  .status-enter,
  .status-leave-to {
    opacity: 0;
  }
</style>
