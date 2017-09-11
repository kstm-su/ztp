<template>
  <md-sidenav class="md-left md-fixed main-sidebar" md-swipeable ref="sidenav">
    <md-list class="md-dense links">
      <md-list-item>
        <router-link exact to="/">root</router-link>
      </md-list-item>
      <md-list-item>
        <router-link exact to="/foo">foo</router-link>
      </md-list-item>
      <md-list-item>
        <router-link exact to="/bar">bar</router-link>
      </md-list-item>
    </md-list>
  </md-sidenav>
</template>

<script>
  export default {
    methods: {
      toggle() {
        this.$refs.sidenav.toggle();
      },
    },
  };
</script>

<style lang="scss">
  @import '../../node_modules/vue-material/src/core/stylesheets/variables.scss';
  @import '../config.scss';

  .main-sidebar {
    .md-sidenav-content {
      width: $sizebar-size;
      display: flex;
      flex-flow: column;
      overflow: hidden;

      @media (min-width: 1281px) {
        top: 0;
        pointer-events: auto;
        transform: translate3d(0, 0, 0) !important;
        box-shadow: $material-shadow-2dp;
      }

      .links.md-list {
        padding: 0;
      }
    }

    .md-backdrop {
      @media (min-width: 1281px) {
        opacity: 0;
        pointer-events: none;
      }
    }
  }
</style>
