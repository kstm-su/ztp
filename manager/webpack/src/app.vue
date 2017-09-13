<template>
  <div class="container">
    <main-sidebar ref="mainSidebar"></main-sidebar>
    <main-content></main-content>
  </div>
</template>

<script>
  import MainSidebar from './components/main-sidebar';
  import MainContent from './components/main-content';

  export default {
    components: {
      MainSidebar,
      MainContent,
    },
    methods: {
      toggleMainSidebar() {
        this.$refs.mainSidebar.toggle();
      },
    },
  };
</script>

<style lang="scss">
  @import '../node_modules/vue-material/src/core/stylesheets/variables.scss';
  @import './config.scss';

  [v-cloak] {
    display: none;
  }

  html,
  body {
    height: 100%;
    overflow: hidden;
  }

  body {
    display: flex;
  }

  .container {
    min-height: 100%;
    display: flex;
    flex: 1;
    flex-flow: column nowrap;
    transition: $swift-ease-out;

    @media (min-width: 1281px) {
      padding-left: $sizebar-size;
    }
  }

  .main-header {
    .nav-trigger {
      @media (min-width: 1281px) {
        display: none;
      }
    }
  }

  .router-enter-active,
  .router-leave-active {
    transition: $swift-ease-out;
  }

  .router-enter-active {
    transition-duration: .25s;
  }

  .router-enter,
  .router-leave-to {
    opacity: 0;
  }

  .router-enter {
    transform: translate3D(0, 10%, 0);
  }
</style>
