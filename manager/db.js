const Sequelize = require('sequelize');
const sequelize = new Sequelize('sqlite://db.sqlite3');

const Images = sequelize.define('images', {
  path: {
    type: Sequelize.STRING,
    unique: true,
    allowNull: true,
  },
  config: {
    type: Sequelize.TEXT,
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

const Nodes = sequelize.define('nodes', {
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

Nodes.belongsTo(Images);
Images.hasMany(Nodes);

sequelize.authenticate().then(() => {
  console.log('success to connect db');
  Images.sync();
  Nodes.sync();
}).catch(err => {
  console.error('unable to connect db:', err);
});

module.exports = (req, res, next) => {
  req.images = Images;
  req.nodes = Nodes;
  next();
};
