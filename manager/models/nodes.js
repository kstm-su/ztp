const Sequelize = require('sequelize');
const db = require('../db');

const Nodes = db.define('nodes', {
  name: {
    type: Sequelize.STRING,
    allowNull: true,
  },
  mac_address: {
    type: Sequelize.STRING,
    unique: true,
    allowNull: false,
    validate: {
      is: /([0-9A-F]{2}:){5}[0-9A-F]{2}(:[0-9A-F]{2}:[0-9A-F]{2})?/i,
    },
  },
  ip_address: {
    type: Sequelize.STRING,
    unique: true,
    allowNull: true,
    validate: {
      isIP: true,
    },
  },
}, {
  timestamps: true,
  underscored: true,
});

module.exports = Nodes;
