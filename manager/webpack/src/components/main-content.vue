<template>
  <div class="main-content">
    <md-whiteframe md-tag="md-toolbar" md-elevation="1" class="main-header">
      <md-button class="md-icon-button nav-trigger" @click="$parent.toggleMainSidebar">
        <md-icon>menu</md-icon>
      </md-button>
      <h2 class="md-title">kstm ZTP</h2>
    </md-whiteframe>
    <div class="page">
      <transition name="router" mode="out-in" appear>
        <router-view></router-view>
      </transition>
    </div>
  </div>
</template>

<style lang="scss" scoped>
  @import '../../node_modules/vue-material/src/core/stylesheets/variables.scss';

  .main-content {
    min-height: 100%;
    max-height: 100%;
    display: flex;
    flex-direction: column;
  }

  .main-header {
    .nav-trigger {
      @media (min-width: 1281px) {
        display: none;
      }
    }

    .md-title {
      flex: 1;
    }
  }

  .page {
    overflow: auto;
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
