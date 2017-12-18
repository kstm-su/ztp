const Sequelize = require('sequelize');
const sequelize = new Sequelize(process.env.DATABASE_URL || 'sqlite://db.sqlite3');

sequelize.authenticate().then(() => {
  console.log('success to connect db');
}).then(() => require('./models').migration())
.catch(err => {
  console.error('unable to connect db:', err);
  process.exit();
});

module.exports = sequelize;
