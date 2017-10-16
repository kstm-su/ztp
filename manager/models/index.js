const Images = require('./images');
const Nodes = require('./nodes');

Nodes.belongsTo(Images);
Images.hasMany(Nodes);

module.exports = {
  Images,
  Nodes,
  migration(succ, err) {
    var pImages = new Promise((resolve, reject) => {
      Images.sync().then(resolve).catch(reject);
    });
    var pNodes = new Promise((resolve, reject) => {
      Nodes.sync().then(resolve).catch(reject);
    });
    Promise.all([pImages, pNodes]).then(succ).catch(err);
  },
};
