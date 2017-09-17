<template>
  <div>
    <md-table-card>
      <md-toolbar>
        <h1 class="md-title">Nodes</h1>
        <md-table-alternate-header md-selected-label="selected">
          <md-button class="md-accent">
            <md-icon>replay</md-icon>
            restart
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
            <md-table-head>host name</md-table-head>
            <md-table-head>image</md-table-head>
            <md-table-head>MAC address</md-table-head>
            <md-table-head>IP address</md-table-head>
          </md-table-row>
        </md-table-header>
        <transition-group name="list" tag="md-table-body">
          <md-table-row
            v-for="node in sortedNodes"
            :key="node.id"
            :id="`node-${node.id}`"
            class="list-item"
            :md-item="node"
            md-selection
          >
            <md-table-cell>{{ node.id }}</md-table-cell>
            <md-table-cell>{{ node.name }}</md-table-cell>
            <md-table-cell>
              <md-select v-model="node.image_id" @change="change(node)">
                <md-option class="null" :value="null">&lt;null&gt;</md-option>
                <md-option v-for="image in images" v-model="image.id">#{{ image.id }}: {{ image.name }}</md-option>
              </md-select>
            </md-table-cell>
            <md-table-cell>{{ node.mac_address.toUpperCase() }}</md-table-cell>
            <md-table-cell>
              <span v-if="node.ip_address == null" class="null">null</span>
              <span v-else>{{ node.ip_address }}</span>
            </md-table-cell>
            <md-table-cell>
              <md-button class="md-icon-button edit-button" @click.stop="edit(node)">
                <md-icon>edit</md-icon>
              </md-button>
            </md-table-cell>
          </md-table-row>
        </transition-group>
      </md-table>
    </md-table-card>
    <md-button class="md-fab md-fab-bottom-right" @click="edit(newNode)">
      <md-icon>add</md-icon>
    </md-button>
    <md-dialog :open-from="editingSelector" :close-to="editingSelector" ref="editDialog">
      <md-dialog-title>
        <span v-if="editing.id == null">new node</span>
        <span v-else>edit node #{{ editing.id }}</span>
      </md-dialog-title>
      <md-dialog-content>
        <form>
          <md-input-container>
            <label>host name</label>
            <md-input v-model="editing.name"></md-input>
          </md-input-container>
          <md-input-container>
            <label>image</label>
            <md-select v-model="editing.image_id">
              <md-option class="null" :value="null">&lt;null&gt;</md-option>
              <md-option v-for="image in sortedImages" v-model="image.id">
                #{{ image.id }}: {{ image.name }}
              </md-option>
            </md-select>
          </md-input-container>
          <md-input-container>
            <label>MAC address</label>
            <md-input v-model="editing.mac_address"></md-input>
          </md-input-container>
        </form>
      </md-dialog-content>
      <md-dialog-actions>
        <md-button class="md-primary" @click="closeEditDialog">Cancel</md-button>
        <md-button class="md-primary" @click="submit">OK</md-button>
      </md-dialog-actions>
    </md-dialog>
    <md-dialog-confirm
      md-title="delete nodes"
      md-content="Are you sure to delete?"
      ref="confirmDeleteDialog"
      @close="deleteSelectedNodes"
    >
    </md-dialog-confirm>
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
          image_id: null,
          mac_address: '',
        },
        editing: {},
      };
    },
    computed: {
      sortedNodes() {
        return this.nodes.sort((a, b) => {
          return new Date(b.updated_at) - new Date(a.updated_at);
        });
      },
      sortedImages() {
        return this.images.sort((a, b) => {
          return new Date(b.updated_at) - new Date(a.updated_at);
        });
      },
      editingSelector() {
        if (this.editing.id == null) {
          return null;
        }
        return `#node-${this.editing.id}`;
      },
    },
    sockets: {
      node(node) {
        node.image = this.images.find(image => image.id === node.image_id);
        this.nodes = this.nodes.map(i => {
          if (i.id === node.id) {
            return node;
          }
          return i;
        });
      },
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
        return this.$http.get('/nodes').then(resp => {
          this.nodes = resp.data;
        });
      },
      submit() {
        delete this.editing.image;
        if (this.editing.id == null) {
          return this.$http.post('/nodes', this.editing).then(resp => {
            this.nodes.unshift(resp.data);
            this.newNode = {
              name: '',
              image_id: null,
              mac_address: '',
            };
            this.closeEditDialog();
          });
        }
        return this.$http.put(`/nodes/${this.editing.id}`, this.editing).then(resp => {
          this.closeEditDialog();
        });
      },
      fetchImages() {
        return this.$http.get('/images').then(resp => {
          this.images = resp.data;
        });
      },
      confirmDelete() {
        this.$refs.confirmDeleteDialog.open();
      },
      deleteSelectedNodes(type) {
        if (type !== 'ok') {
          return;
        }
        this.$refs.table.selectedRows.forEach(node => {
          this.$http.delete(`/nodes/${node.id}`).then(resp => {
            this.nodes = this.nodes.filter(v => v.id !== node.id);
          });
        });
      },
      edit(node) {
        this.editing = node;
        this.$refs.editDialog.open();
      },
      change(node) {
        this.editing = node;
        this.submit();
      },
      closeEditDialog() {
        this.$refs.editDialog.close();
        this.editing = {};
      },
    },
    mounted() {
      this.fetch();
      this.fetchImages();
    },
  };
</script>

<style scoped>
  >>> .md-dialog {
    min-width: 80%;
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
