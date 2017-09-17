<template>
  <div>
    <md-table-card>
      <md-toolbar>
        <h1 class="md-title">Images</h1>
        <md-table-alternate-header md-selected-label="selected">
          <md-button class="md-accent" @click="rebuild">
            <md-icon>sync</md-icon>
            rebuild
          </md-button>
          <md-button class="md-accent" @click="confirmDelete">
            <md-icon>delete</md-icon>
            delete
          </md-button>
        </md-table-alternate-header>
      </md-toolbar>
      <md-table ref="table">
        <md-table-header>
          <md-table-row>
            <md-table-head>#</md-table-head>
            <md-table-head>name</md-table-head>
            <md-table-head>description</md-table-head>
            <md-table-head>size</md-table-head>
            <md-table-head>status</md-table-head>
          </md-table-row>
        </md-table-header>
        <transition-group name="list" tag="md-table-body">
          <md-table-row
            v-for="image in sortedImages"
            :key="image.id"
            :id="`image-${image.id}`"
            class="list-item"
            :md-item="image"
            md-selection
          >
            <md-table-cell>{{ image.id }}</md-table-cell>
            <md-table-cell>{{ image.name }}</md-table-cell>
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
              <md-button class="md-icon-button edit-button" @click.stop="edit(image)">
                <md-icon>edit</md-icon>
              </md-button>
            </md-table-cell>
          </md-table-row>
        </transition-group>
      </md-table>
    </md-table-card>
    <md-button class="md-fab md-fab-bottom-right" @click="edit(newImage)">
      <md-icon>add</md-icon>
    </md-button>
    <md-dialog :open-from="editingSelector" :close-to="editingSelector" ref="editDialog">
      <md-dialog-title>
        <span v-if="editing.id == null">new image</span>
        <span v-else>edit image #{{ editing.id }}</span>
      </md-dialog-title>
      <md-dialog-content>
        <form>
          <md-input-container>
            <label>image name</label>
            <md-input v-model="editing.name"></md-input>
          </md-input-container>
          <md-input-container>
            <label>LinuxKit config</label>
            <md-textarea v-model="editing.config"></md-textarea>
          </md-input-container>
          <md-input-container>
            <label>image description</label>
            <md-textarea v-model="editing.description"></md-textarea>
          </md-input-container>
          <md-input-container>
            <label>image size</label>
            <md-input type="number" v-model.number="editing.size"></md-input>
          </md-input-container>
        </form>
      </md-dialog-content>
      <md-dialog-actions>
        <md-button class="md-primary" @click="closeEditDialog">Cancel</md-button>
        <md-button class="md-primary" @click="submit">OK</md-button>
      </md-dialog-actions>
    </md-dialog>
    <md-dialog-confirm
      md-title="delete images"
      md-content="Are you sure to delete?"
      ref="confirmDeleteDialog"
      @close="deleteSelectedImages"
    >
    </md-dialog-confirm>
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
        editing: {},
      };
    },
    computed: {
      sortedImages() {
        return this.images.sort((a, b) => {
          return new Date(b.updated_at) - new Date(a.updated_at);
        });
      },
      editingSelector() {
        if (this.editing.id == null) {
          return null;
        }
        return `#image-${this.editing.id}`;
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
        this.editing.build = true;
        if (this.editing.id == null) {
          return this.$http.post('/images', this.editing).then(resp => {
            this.images.unshift(resp.data);
            this.newImage = {
              name: '',
              config: '',
              description: '',
              size: this.editing.size,
            };
            this.closeEditDialog();
          });
        }
        return this.$http.put(`/images/${this.editing.id}`, this.editing).then(resp => {
          this.closeEditDialog();
        });
      },
      confirmDelete() {
        this.$refs.confirmDeleteDialog.open();
      },
      deleteSelectedImages(type) {
        if (type !== 'ok') {
          return;
        }
        this.$refs.table.selectedRows.forEach(image => {
          this.$http.delete(`/images/${image.id}`).then(resp => {
            this.images = this.images.filter(v => i.id !== image.id);
          });
        });
      },
      rebuild() {
        this.$refs.table.selectedRows.forEach(image => {
          this.$http.put(`/images/${image.id}`, {
            build: true,
          });
        });
      },
      edit(image) {
        this.editing = image;
        this.$refs.editDialog.open();
      },
      closeEditDialog() {
        this.$refs.editDialog.close();
        this.editing = {};
      },
    },
    mounted() {
      this.fetch();
    },
  };
</script>

<style scoped>
  >>> .md-dialog {
    min-width: 80%;
  }

  >>> .md-spinner {
    vertical-align: middle;
  }

  tr:not(:hover) .edit-button {
    opacity: 0;
  }

  .list-item {
    transition: all ease .5s;
  }

  .list-enter,
  .list-leave-to {
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
