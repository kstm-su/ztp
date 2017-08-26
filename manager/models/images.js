const Sequelize = require('sequelize');
const db = require('../db');

const Images = db.define('images', {
  path: {
    type: Sequelize.STRING,
    unique: true,
    allowNull: true,
  },
  config: {
    type: Sequelize.TEXT,
    allowNull: false,
  },
  size: {
    type: Sequelize.INTEGER,
  },
  name: {
    type: Sequelize.STRING,
    allowNull: false,
  },
  description: {
    type: Sequelize.TEXT,
    allowNull: false,
  },
  error: {
    type: Sequelize.TEXT,
    allowNull: true,
  },
}, {
  timestamps: true,
  underscored: true,
});

module.exports = Images;
