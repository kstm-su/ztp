const Images = require('./images');
const Nodes = require('./nodes');

Nodes.belongsTo(Images);
Images.hasMany(Nodes);

module.exports = {
  Images,
  Nodes,
  migration() {
    Images.sync();
    Nodes.sync();
  },
};
